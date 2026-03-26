package operatorconfig

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

func TestManagerLoadInitial(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 scheme: %v", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorConfigName,
			Namespace: defaultNamespace,
		},
		Data: map[string]string{"value": liveConfigValue},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cm).Build()

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
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 scheme: %v", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorConfigName,
			Namespace: defaultNamespace,
		},
		Data: map[string]string{"value": liveConfigValue},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cm).Build()

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
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 scheme: %v", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorConfigName,
			Namespace: defaultNamespace,
		},
		Data: map[string]string{"value": liveConfigValue},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cm).Build()

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

	if _, err := manager.Reconcile(context.Background(), reconcile.Request{
		NamespacedName: types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
	}); err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}
	if applied != 1 {
		t.Fatalf("expected unchanged reconcile to skip apply, got %d", applied)
	}
}

func TestConfigMapPredicateHandlesNilEvents(t *testing.T) {
	manager := &Manager[sampleConfig]{
		opts: Options[sampleConfig]{
			ConfigMapKey: types.NamespacedName{Name: operatorConfigName, Namespace: defaultNamespace},
		},
	}
	pred := manager.configMapPredicate()
	if pred.Create(event.CreateEvent{}) {
		t.Fatalf("expected nil create event to be ignored")
	}
	if pred.Update(event.UpdateEvent{}) {
		t.Fatalf("expected nil update event to be ignored")
	}
	if pred.Delete(event.DeleteEvent{}) {
		t.Fatalf("expected nil delete event to be ignored")
	}
}

func TestManagerReconcileNotFoundResetsToDefault(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 scheme: %v", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorConfigName,
			Namespace: defaultNamespace,
		},
		Data: map[string]string{"value": liveConfigValue},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(cm).Build()

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
	if err := client.Delete(context.Background(), cm); err != nil {
		t.Fatalf("delete configmap: %v", err)
	}

	if _, err := manager.Reconcile(context.Background(), reconcile.Request{
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
