package cpf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 200; i++ {
		g := Generate()
		require.Len(t, g, 11)
		require.True(t, Valid(g), "Generate produced an invalid CPF: %s", g)
		seen[g] = true
	}
	// With 200 draws the output must be varied (very likely > 1 distinct value).
	assert.GreaterOrEqual(t, len(seen), 10, "Generate should produce varied numbers")
	assert.False(t, seen["11111111111"], "Generate must never produce repeated-digit CPFs")
}
