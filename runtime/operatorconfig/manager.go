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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReloadReason identifies why the configuration was reloaded.
type ReloadReason string

const (
	// ReloadReasonInitial indicates the initial synchronous load at startup.
	ReloadReasonInitial ReloadReason = "initial"
	// ReloadReasonReconcile indicates a reload triggered by the controller watch.
	ReloadReasonReconcile ReloadReason = "reconcile"
)

// Options configures a Manager instance for a concrete operator config type.
type Options[T any] struct {
	Client          client.Client
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
	apiReader     client.Reader
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
	if opts.Logger.GetSink() == nil {
		opts.Logger = log.Log.WithName("operator-config-manager")
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
func (m *Manager[T]) SetAPIReader(reader client.Reader) {
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

// SetupWithManager wires the manager into controller-runtime so ConfigMap
// changes trigger reconciles.
func (m *Manager[T]) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named(m.opts.ControllerName).
		For(&corev1.ConfigMap{}).
		WithEventFilter(m.configMapPredicate()).
		Complete(m)
}

// Reconcile reacts to ConfigMap updates/deletes and refreshes the cache.
func (m *Manager[T]) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	if req.NamespacedName != m.opts.ConfigMapKey {
		return reconcile.Result{}, nil
	}

	cfg, err := m.loadAndParse(ctx)
	if err != nil {
		if apierrors.IsNotFound(err) {
			fallback := m.opts.CloneConfig(m.defaultConfig)
			snapshot, changed := m.storeConfigIfChanged(fallback)
			if changed {
				m.applyAndNotify(ReloadReasonReconcile, snapshot)
			}
			return reconcile.Result{}, nil
		}
		m.opts.Logger.Error(err, "failed to refresh operator configuration")
		return reconcile.Result{}, err
	}

	snapshot, changed := m.storeConfigIfChanged(cfg)
	if changed {
		m.applyAndNotify(ReloadReasonReconcile, snapshot)
	}
	return reconcile.Result{}, nil
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

func (m *Manager[T]) reader() client.Reader {
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

func (m *Manager[T]) configMapPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			if e.Object == nil {
				return false
			}
			return namespacedName(e.Object.GetNamespace(), e.Object.GetName()) == m.opts.ConfigMapKey
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectNew == nil {
				return false
			}
			return namespacedName(e.ObjectNew.GetNamespace(), e.ObjectNew.GetName()) == m.opts.ConfigMapKey
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			if e.Object == nil {
				return false
			}
			return namespacedName(e.Object.GetNamespace(), e.Object.GetName()) == m.opts.ConfigMapKey
		},
		GenericFunc: func(event.GenericEvent) bool { return false },
	}
}

func namespacedName(ns, name string) types.NamespacedName {
	return types.NamespacedName{Namespace: ns, Name: name}
}
