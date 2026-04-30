package operatorconfig

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type sampleConfig struct {
	Value string
}

const (
	defaultConfigValue = "default"
	defaultNamespace   = "default"
	liveConfigValue    = "live"
	operatorConfigName = "operator-config"
)

func defaultSampleConfig() *sampleConfig {
	return &sampleConfig{Value: defaultConfigValue}
}

func cloneSampleConfig(in *sampleConfig) *sampleConfig {
	if in == nil {
		return nil
	}
	copy := *in
	return &copy
}

func parseSampleConfig(cm *corev1.ConfigMap) (*sampleConfig, error) {
	val := cm.Data["value"]
	if val == "" {
		val = defaultConfigValue
	}
	return &sampleConfig{Value: val}, nil
}

type memoryConfigMapReader struct {
	items map[types.NamespacedName]*corev1.ConfigMap
}

func newMemoryConfigMapReader(items ...*corev1.ConfigMap) *memoryConfigMapReader {
	reader := &memoryConfigMapReader{items: make(map[types.NamespacedName]*corev1.ConfigMap, len(items))}
	for _, item := range items {
		reader.put(item)
	}
	return reader
}

func (r *memoryConfigMapReader) Get(_ context.Context, key types.NamespacedName, into *corev1.ConfigMap) error {
	item := r.items[key]
	if item == nil {
		return apierrors.NewNotFound(schema.GroupResource{Resource: "configmaps"}, key.Name)
	}
	item.DeepCopyInto(into)
	return nil
}

func (r *memoryConfigMapReader) delete(item *corev1.ConfigMap) {
	delete(r.items, types.NamespacedName{Name: item.Name, Namespace: item.Namespace})
}

func (r *memoryConfigMapReader) put(item *corev1.ConfigMap) {
	if item == nil {
		return
	}
	r.items[types.NamespacedName{Name: item.Name, Namespace: item.Namespace}] = item.DeepCopy()
}

func TestManagerLoadInitial(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorConfigName,
			Namespace: defaultNamespace,
		},
		Data: map[string]string{"value": liveConfigValue},
	}

	client := newMemoryConfigMapReader(cm)

	applied := false
	manager, err := NewManager[sampleConfig](Options[sampleConfig]{
		Client:        client,
		ConfigMapKey:  types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
		DefaultConfig: defaultSampleConfig,
		ParseConfigMap: func(cm *corev1.ConfigMap) (*sampleConfig, error) {
			return parseSampleConfig(cm)
		},
		CloneConfig: cloneSampleConfig,
		ApplyConfig: func(*sampleConfig) { applied = true },
	})
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if err := manager.LoadInitial(context.Background()); err != nil {
		t.Fatalf("LoadInitial failed: %v", err)
	}

	cfg := manager.CurrentConfig()
	if cfg.Value != liveConfigValue {
		t.Fatalf("expected config value %q, got %q", liveConfigValue, cfg.Value)
	}
	if !applied {
		t.Fatalf("expected ApplyConfig to be invoked")
	}

	manager.ResetToDefault()
	cfg = manager.CurrentConfig()
	if cfg.Value != defaultConfigValue {
		t.Fatalf("expected default config after reset, got %q", cfg.Value)
	}
}

func TestNewManagerReturnsValidationErrors(t *testing.T) {
	_, err := NewManager[sampleConfig](Options[sampleConfig]{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestManagerCallbacksAndCurrentConfigReceiveClones(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorConfigName,
			Namespace: defaultNamespace,
		},
		Data: map[string]string{"value": liveConfigValue},
	}

	client := newMemoryConfigMapReader(cm)

	manager, err := NewManager[sampleConfig](Options[sampleConfig]{
		Client:        client,
		ConfigMapKey:  types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
		DefaultConfig: defaultSampleConfig,
		ParseConfigMap: func(cm *corev1.ConfigMap) (*sampleConfig, error) {
			return parseSampleConfig(cm)
		},
		CloneConfig: cloneSampleConfig,
		ApplyConfig: func(cfg *sampleConfig) {
			cfg.Value = "mutated-by-apply"
		},
		OnConfigApplied: func(_ ReloadReason, cfg *sampleConfig) {
			cfg.Value = "mutated-by-callback"
		},
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	if err := manager.LoadInitial(context.Background()); err != nil {
		t.Fatalf("LoadInitial failed: %v", err)
	}

	cfg := manager.CurrentConfig()
	if cfg.Value != liveConfigValue {
		t.Fatalf("expected stored config to remain %q, got %q", liveConfigValue, cfg.Value)
	}

	cfg.Value = "mutated-by-caller"
	if manager.CurrentConfig().Value != liveConfigValue {
		t.Fatalf("expected CurrentConfig to return a clone")
	}
}

func TestManagerReconcileSkipsUnchangedConfig(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorConfigName,
			Namespace: defaultNamespace,
		},
		Data: map[string]string{"value": liveConfigValue},
	}

	client := newMemoryConfigMapReader(cm)

	applied := 0
	manager, err := NewManager[sampleConfig](Options[sampleConfig]{
		Client:        client,
		ConfigMapKey:  types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
		DefaultConfig: defaultSampleConfig,
		ParseConfigMap: func(cm *corev1.ConfigMap) (*sampleConfig, error) {
			return parseSampleConfig(cm)
		},
		CloneConfig: cloneSampleConfig,
		ApplyConfig: func(*sampleConfig) {
			applied++
		},
	})
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if err := manager.LoadInitial(context.Background()); err != nil {
		t.Fatalf("LoadInitial failed: %v", err)
	}
	if applied != 1 {
		t.Fatalf("expected first apply during LoadInitial, got %d", applied)
	}

	if _, err := manager.Reconcile(context.Background(), Request{
		NamespacedName: types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
	}); err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}
	if applied != 1 {
		t.Fatalf("expected unchanged reconcile to skip apply, got %d", applied)
	}
}

func TestMatchesConfigMapHandlesNilObject(t *testing.T) {
	manager := &Manager[sampleConfig]{
		opts: Options[sampleConfig]{
			ConfigMapKey: types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
		},
	}
	if manager.MatchesConfigMap(nil) {
		t.Fatalf("expected nil object to be ignored")
	}
}

func TestManagerReconcileNotFoundResetsToDefault(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorConfigName,
			Namespace: defaultNamespace,
		},
		Data: map[string]string{"value": liveConfigValue},
	}

	client := newMemoryConfigMapReader(cm)

	applied := 0
	manager, err := NewManager[sampleConfig](Options[sampleConfig]{
		Client:        client,
		ConfigMapKey:  types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
		DefaultConfig: defaultSampleConfig,
		ParseConfigMap: func(cm *corev1.ConfigMap) (*sampleConfig, error) {
			return parseSampleConfig(cm)
		},
		CloneConfig: cloneSampleConfig,
		ApplyConfig: func(*sampleConfig) {
			applied++
		},
	})
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if err := manager.LoadInitial(context.Background()); err != nil {
		t.Fatalf("LoadInitial failed: %v", err)
	}
	client.delete(cm)

	if _, err := manager.Reconcile(context.Background(), Request{
		NamespacedName: types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
	}); err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}

	if got := manager.CurrentConfig().Value; got != defaultConfigValue {
		t.Fatalf("expected reset to default after not found, got %q", got)
	}
	if applied != 2 {
		t.Fatalf("expected apply on load and fallback reset, got %d", applied)
	}
}
