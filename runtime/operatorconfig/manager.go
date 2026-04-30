package operatorconfig

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ReloadReason identifies why the configuration was reloaded.
type ReloadReason string

const (
	// ReloadReasonInitial indicates the initial synchronous load at startup.
	ReloadReasonInitial ReloadReason = "initial"
	// ReloadReasonReconcile indicates a reload triggered by the controller watch.
	ReloadReasonReconcile ReloadReason = "reconcile"
)

// ConfigMapReader is the minimal read contract needed by Manager.
type ConfigMapReader interface {
	Get(ctx context.Context, key types.NamespacedName, into *corev1.ConfigMap) error
}

// ConfigMapReaderFunc adapts a function to ConfigMapReader.
type ConfigMapReaderFunc func(context.Context, types.NamespacedName, *corev1.ConfigMap) error

// Get implements ConfigMapReader.
func (f ConfigMapReaderFunc) Get(ctx context.Context, key types.NamespacedName, into *corev1.ConfigMap) error {
	return f(ctx, key, into)
}

// Request identifies a reconcile target.
type Request struct {
	NamespacedName types.NamespacedName
}

// Result is reserved for future reconcile scheduling metadata.
type Result struct{}

// Options configures a Manager instance for a concrete operator config type.
type Options[T any] struct {
	Client          ConfigMapReader
	Logger          logr.Logger
	ConfigMapKey    types.NamespacedName
	ControllerName  string
	DefaultConfig   func() *T
	ParseConfigMap  func(*corev1.ConfigMap) (*T, error)
	CloneConfig     func(*T) *T
	ApplyConfig     func(*T)
	OnConfigApplied func(ReloadReason, *T)
}

// Manager encapsulates shared ConfigMap loading, caching, and reconcile logic
// for operator configuration structs.
type Manager[T any] struct {
	opts          Options[T]
	apiReader     ConfigMapReader
	currentConfig *T
	defaultConfig *T
	lastSync      time.Time
	mu            sync.RWMutex
}

// NewManager constructs a Manager with the provided options.
func NewManager[T any](opts Options[T]) (*Manager[T], error) {
	if opts.Client == nil {
		return nil, fmt.Errorf("operatorconfig: Client must be provided")
	}
	if opts.DefaultConfig == nil {
		return nil, fmt.Errorf("operatorconfig: DefaultConfig must be provided")
	}
	if opts.ParseConfigMap == nil {
		return nil, fmt.Errorf("operatorconfig: ParseConfigMap must be provided")
	}
	if opts.CloneConfig == nil {
		return nil, fmt.Errorf("operatorconfig: CloneConfig must be provided")
	}
	if opts.ControllerName == "" {
		opts.ControllerName = "operator-config-manager"
	}
	defaultSeed := opts.DefaultConfig()
	if defaultSeed == nil {
		return nil, fmt.Errorf("operatorconfig: DefaultConfig must not return nil")
	}
	defaultCfg := opts.CloneConfig(defaultSeed)
	if defaultCfg == nil {
		return nil, fmt.Errorf("operatorconfig: CloneConfig must not return nil")
	}
	currentCfg := opts.CloneConfig(defaultCfg)
	if currentCfg == nil {
		return nil, fmt.Errorf("operatorconfig: CloneConfig must not return nil")
	}
	return &Manager[T]{
		opts:          opts,
		currentConfig: currentCfg,
		defaultConfig: opts.CloneConfig(defaultCfg),
	}, nil
}

// SetAPIReader injects a non-cached reader for startup scenarios.
func (m *Manager[T]) SetAPIReader(reader ConfigMapReader) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.apiReader = reader
}

// CurrentConfig returns a cloned snapshot of the cached configuration.
func (m *Manager[T]) CurrentConfig() *T {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.currentConfig == nil {
		return nil
	}
	return m.opts.CloneConfig(m.currentConfig)
}

// ResetToDefault replaces the cached configuration with the default snapshot.
func (m *Manager[T]) ResetToDefault() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentConfig = m.opts.CloneConfig(m.defaultConfig)
	m.lastSync = time.Now()
}

// LoadInitial synchronously loads the ConfigMap before the manager starts.
func (m *Manager[T]) LoadInitial(ctx context.Context) error {
	cfg, err := m.loadAndParse(ctx)
	if err != nil {
		return err
	}
	snapshot := m.storeConfig(cfg)
	m.applyAndNotify(ReloadReasonInitial, snapshot)
	return nil
}

// Reconcile reacts to ConfigMap updates/deletes and refreshes the cache.
func (m *Manager[T]) Reconcile(ctx context.Context, req Request) (Result, error) {
	if req.NamespacedName != m.opts.ConfigMapKey {
		return Result{}, nil
	}

	cfg, err := m.loadAndParse(ctx)
	if err != nil {
		if apierrors.IsNotFound(err) {
			fallback := m.opts.CloneConfig(m.defaultConfig)
			snapshot, changed := m.storeConfigIfChanged(fallback)
			if changed {
				m.applyAndNotify(ReloadReasonReconcile, snapshot)
			}
			return Result{}, nil
		}
		m.opts.Logger.Error(err, "failed to refresh operator configuration")
		return Result{}, err
	}

	snapshot, changed := m.storeConfigIfChanged(cfg)
	if changed {
		m.applyAndNotify(ReloadReasonReconcile, snapshot)
	}
	return Result{}, nil
}

func (m *Manager[T]) loadAndParse(ctx context.Context) (*T, error) {
	reader := m.reader()
	if reader == nil {
		return nil, fmt.Errorf("operatorconfig: no client or apiReader configured")
	}

	configMap := &corev1.ConfigMap{}
	if err := reader.Get(ctx, m.opts.ConfigMapKey, configMap); err != nil {
		return nil, err
	}

	cfg, err := m.opts.ParseConfigMap(configMap)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("operatorconfig: ParseConfigMap returned nil")
	}
	return cfg, nil
}

func (m *Manager[T]) reader() ConfigMapReader {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.apiReader != nil {
		return m.apiReader
	}
	return m.opts.Client
}

func (m *Manager[T]) storeConfig(cfg *T) *T {
	snapshot := m.opts.CloneConfig(cfg)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentConfig = snapshot
	m.lastSync = time.Now()
	return m.opts.CloneConfig(snapshot)
}

func (m *Manager[T]) storeConfigIfChanged(cfg *T) (*T, bool) {
	snapshot := m.opts.CloneConfig(cfg)
	m.mu.Lock()
	defer m.mu.Unlock()
	if reflect.DeepEqual(m.currentConfig, snapshot) {
		return m.opts.CloneConfig(m.currentConfig), false
	}
	m.currentConfig = snapshot
	m.lastSync = time.Now()
	return m.opts.CloneConfig(snapshot), true
}

func (m *Manager[T]) applyAndNotify(reason ReloadReason, cfg *T) {
	if m.opts.ApplyConfig != nil {
		m.opts.ApplyConfig(m.opts.CloneConfig(cfg))
	}
	if m.opts.OnConfigApplied != nil {
		m.opts.OnConfigApplied(reason, m.opts.CloneConfig(cfg))
	}
}

// MatchesConfigMap reports whether the object is the configured ConfigMap.
func (m *Manager[T]) MatchesConfigMap(obj metav1.Object) bool {
	if obj == nil {
		return false
	}
	return types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.GetName()} == m.opts.ConfigMapKey
}
