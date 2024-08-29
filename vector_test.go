package yamlvector

import (
	"testing"

	"github.com/koykov/vector"
)

func TestScalar(t *testing.T) {
	vec := NewVector()
	t.Run("scalar_null", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNull)
	})
	t.Run("scalar_null_canonical", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNull)
	})
	t.Run("scalar_null_none", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNull)
	})
}
