package templating

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasStorageRefsInStepVars(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		text string
		vars map[string]any
		want bool
	}{
		{
			name: "compound template with storage ref",
			text: `{{ (index .steps "x").output.body | toJson }}`,
			vars: map[string]any{
				"steps": map[string]any{
					"x": map[string]any{
						"output": map[string]any{
							"$bubuStorageRef": "s3://bucket/key",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "compound template without storage ref",
			text: `{{ (index .steps "x").output.body | toJson }}`,
			vars: map[string]any{
				"steps": map[string]any{
					"x": map[string]any{
						"output": map[string]any{
							"body": "hello",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "no steps reference",
			text: `{{ .inputs.foo | toJson }}`,
			vars: map[string]any{
				"inputs": map[string]any{"foo": "bar"},
			},
			want: false,
		},
		{
			name: "single action template with storage ref",
			text: `{{ (index .steps "x").output.body }}`,
			vars: map[string]any{
				"steps": map[string]any{
					"x": map[string]any{
						"output": map[string]any{
							"$bubuStorageRef": "s3://bucket/key",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "step not in vars",
			text: `{{ (index .steps "missing").output.body | toJson }}`,
			vars: map[string]any{
				"steps": map[string]any{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := HasStorageRefsInStepVars(tt.text, tt.vars)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestHasStorageRefsInStepVarsSupportsAnyMaps(t *testing.T) {
	t.Parallel()

	vars := map[string]any{
		RootSteps: map[any]any{
			"x": map[any]any{
				"output": map[any]any{
					StorageRefKey: "s3://bucket/key",
				},
			},
		},
	}

	require.True(t, HasStorageRefsInStepVars(`{{ (index .steps "x").output.body | toJson }}`, vars))
}

func TestContainsStorageRefHandlesCyclesAndDeepNesting(t *testing.T) {
	t.Parallel()

	root := map[string]any{}
	current := root
	for range 80 {
		next := map[string]any{}
		current["child"] = next
		current = next
	}
	current[StorageRefKey] = "s3://bucket/key"
	root["self"] = root

	require.True(t, containsStorageRef(root))
}
