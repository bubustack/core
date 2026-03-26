package storage

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/bubustack/core/contracts"
)

func TestApplyEnvERejectsAmbiguousProvider(t *testing.T) {
	err := ApplyEnvE(&Config{
		S3:   &S3Config{},
		File: &FileConfig{},
	}, &corev1.PodSpec{}, &corev1.Container{}, "5s")
	if err == nil {
		t.Fatalf("expected ambiguous provider error")
	}
}

func TestApplyEnvERejectsInvalidTimeout(t *testing.T) {
	err := ApplyEnvE(&Config{
		File: &FileConfig{},
	}, &corev1.PodSpec{}, &corev1.Container{}, "not-a-duration")
	if err == nil {
		t.Fatalf("expected invalid timeout error")
	}
}

func TestApplyEnvEFileConfigAppliesDefaults(t *testing.T) {
	podSpec := &corev1.PodSpec{}
	container := &corev1.Container{}

	err := ApplyEnvE(&Config{
		File: &FileConfig{},
	}, podSpec, container, "15s")
	if err != nil {
		t.Fatalf("ApplyEnvE failed: %v", err)
	}

	if len(podSpec.Volumes) != 1 {
		t.Fatalf("expected one volume, got %d", len(podSpec.Volumes))
	}
	if len(container.VolumeMounts) != 1 || container.VolumeMounts[0].MountPath != "/var/run/bubu/storage" {
		t.Fatalf("unexpected volume mounts: %#v", container.VolumeMounts)
	}
	assertContainerEnvValue(t, container, contracts.StorageProviderEnv, "file")
	assertContainerEnvValue(t, container, contracts.StoragePathEnv, "/var/run/bubu/storage")
	assertContainerEnvValue(t, container, contracts.StorageTimeoutEnv, "15s")
}

func assertContainerEnvValue(t *testing.T, container *corev1.Container, name, expected string) {
	t.Helper()
	for _, env := range container.Env {
		if env.Name == name {
			if env.Value != expected {
				t.Fatalf("expected %s=%s, got %s", name, expected, env.Value)
			}
			return
		}
	}
	t.Fatalf("expected env %s to be present", name)
}
