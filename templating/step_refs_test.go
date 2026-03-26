package templating

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractStepReferences(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "index expression",
			input: `{{ (index .steps "fetch-feed").output.body }}`,
			want:  []string{"fetch-feed"},
		},
		{
			name:  "dot access",
			input: `{{ .steps.myStep.output.result }}`,
			want:  []string{"myStep"},
		},
		{
			name:  "multiple refs",
			input: `{{ (index .steps "a").output.x }} and {{ (index .steps "b").output.y }}`,
			want:  []string{"a", "b"},
		},
		{
			name:  "no steps ref",
			input: `{{ .inputs.foo }}`,
			want:  nil,
		},
		{
			name:  "not a template",
			input: `plain text`,
			want:  nil,
		},
		{
			name:  "piped expression",
			input: `{{ (index .steps "x").output.body | toJson }}`,
			want:  []string{"x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ExtractStepReferences(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestExtractStepReferencesWithErrorReturnsParseFailure(t *testing.T) {
	t.Parallel()

	refs, err := ExtractStepReferencesWithError(`{{`)
	require.Error(t, err)
	require.Nil(t, refs)
}
