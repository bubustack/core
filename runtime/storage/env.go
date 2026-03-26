package storage

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/bubustack/core/contracts"
)

// Config represents the storage provider configuration required to wire env vars,
// secrets, and volumes into controller-managed workloads.
type Config struct {
	S3   *S3Config
	File *FileConfig
}

// S3Config captures the parameters needed to talk to an S3-compatible backend.
type S3Config struct {
	Bucket       string
	Region       string
	Endpoint     string
	SecretName   string
	UsePathStyle bool
}

// FileConfig configures a filesystem-backed storage target.
type FileConfig struct {
	Path            string
	VolumeClaimName string
	EmptyDir        *corev1.EmptyDirVolumeSource
}

// ApplyEnv applies storage-related env vars, volumes, and mounts to the provided
// container/pod spec. The helper is provider-aware and does nothing when cfg is nil.
func ApplyEnv(cfg *Config, podSpec *corev1.PodSpec, container *corev1.Container, timeout string) {
	_ = ApplyEnvE(cfg, podSpec, container, timeout)
}

// ApplyEnvE applies storage-related env vars, volumes, and mounts to the provided
// container/pod spec, returning configuration errors to callers that want to fail fast.
func ApplyEnvE(cfg *Config, podSpec *corev1.PodSpec, container *corev1.Container, timeout string) error {
	if cfg == nil || container == nil {
		return nil
	}
	if cfg.S3 != nil && cfg.File != nil {
		return fmt.Errorf("storage config must set only one provider")
	}
	timeout = strings.TrimSpace(timeout)
	if timeout != "" {
		if _, err := time.ParseDuration(timeout); err != nil {
			return fmt.Errorf("invalid storage timeout %q: %w", timeout, err)
		}
	}

	switch {
	case cfg.S3 != nil:
		applyS3Env(container, cfg.S3, timeout)
	case cfg.File != nil:
		applyFileEnv(container, podSpec, cfg.File, timeout)
	}
	return nil
}

func applyS3Env(container *corev1.Container, cfg *S3Config, timeout string) {
	upsertEnv(container, corev1.EnvVar{Name: contracts.StorageProviderEnv, Value: "s3"})
	if cfg.Bucket != "" {
		upsertEnv(container, corev1.EnvVar{Name: contracts.StorageS3BucketEnv, Value: cfg.Bucket})
	}
	if cfg.Region != "" {
		upsertEnv(container, corev1.EnvVar{Name: contracts.StorageS3RegionEnv, Value: cfg.Region})
	}
	if cfg.Endpoint != "" {
		upsertEnv(container, corev1.EnvVar{Name: contracts.StorageS3EndpointEnv, Value: cfg.Endpoint})
	}
	if cfg.SecretName != "" {
		addSecretEnv(container, cfg.SecretName, "AWS_ACCESS_KEY_ID")
		addSecretEnv(container, cfg.SecretName, "AWS_SECRET_ACCESS_KEY")
		addSecretEnvOptional(container, cfg.SecretName, "AWS_SESSION_TOKEN")
	}
	if cfg.UsePathStyle {
		upsertEnv(container, corev1.EnvVar{
			Name:  contracts.StorageS3ForcePathStyleEnv,
			Value: strconv.FormatBool(true),
		})
	}
	maybeSetStorageTimeout(container, timeout)
}

func applyFileEnv(container *corev1.Container, podSpec *corev1.PodSpec, cfg *FileConfig, timeout string) {
	mountPath := strings.TrimSpace(cfg.Path)
	if mountPath == "" {
		mountPath = "/var/run/bubu/storage"
	}
	const volumeName = "bubu-storage"
	if podSpec != nil {
		ensureStorageVolume(podSpec, volumeName, cfg)
	}
	ensureVolumeMount(container, volumeName, mountPath)
	upsertEnv(container,
		corev1.EnvVar{Name: contracts.StorageProviderEnv, Value: "file"},
		corev1.EnvVar{Name: contracts.StoragePathEnv, Value: mountPath},
	)
	maybeSetStorageTimeout(container, timeout)
}

func ensureStorageVolume(podSpec *corev1.PodSpec, volumeName string, cfg *FileConfig) {
	for i := range podSpec.Volumes {
		if podSpec.Volumes[i].Name == volumeName {
			return
		}
	}

	var volume corev1.Volume
	switch {
	case cfg.VolumeClaimName != "":
		volume = corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: cfg.VolumeClaimName},
			},
		}
	case cfg.EmptyDir != nil:
		volume = corev1.Volume{
			Name:         volumeName,
			VolumeSource: corev1.VolumeSource{EmptyDir: cfg.EmptyDir.DeepCopy()},
		}
	default:
		volume = corev1.Volume{
			Name:         volumeName,
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		}
	}
	podSpec.Volumes = append(podSpec.Volumes, volume)
}

func ensureVolumeMount(container *corev1.Container, volumeName, mountPath string) {
	for i := range container.VolumeMounts {
		if container.VolumeMounts[i].Name == volumeName && container.VolumeMounts[i].MountPath == mountPath {
			return
		}
	}
	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
	})
}

func addSecretEnv(container *corev1.Container, secretName, key string) {
	upsertEnv(container, corev1.EnvVar{
		Name: key,
		ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
			Key:                  key,
		}},
	})
}

func addSecretEnvOptional(container *corev1.Container, secretName, key string) {
	optional := true
	upsertEnv(container, corev1.EnvVar{
		Name: key,
		ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
			Key:                  key,
			Optional:             &optional,
		}},
	})
}

func maybeSetStorageTimeout(container *corev1.Container, timeout string) {
	if strings.TrimSpace(timeout) == "" {
		return
	}
	upsertEnv(container, corev1.EnvVar{Name: contracts.StorageTimeoutEnv, Value: timeout})
}

func upsertEnv(container *corev1.Container, envVars ...corev1.EnvVar) {
	for _, env := range envVars {
		updated := false
		for i := range container.Env {
			if container.Env[i].Name == env.Name {
				container.Env[i] = env
				updated = true
				break
			}
		}
		if !updated {
			container.Env = append(container.Env, env)
		}
	}
}
