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
	t.Run("scalar_bool", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeBool)
		assertBool(t, vec, "", true)
	})
	t.Run("scalar_bool_low", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeBool)
		assertBool(t, vec, "", false)
	})
	t.Run("scalar_bool_on", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeBool)
		assertBool(t, vec, "", true)
	})
}

func BenchmarkScalar(b *testing.B) {
	b.Run("scalar_null", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeNull)
		})
	})
	b.Run("scalar_null_canonical", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeNull)
		})
	})
	b.Run("scalar_null_none", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeNull)
		})
	})
	b.Run("scalar_bool", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeBool)
		})
	})
	b.Run("scalar_bool_low", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeBool)
		})
	})
	b.Run("scalar_bool_on", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeBool)
		})
	})
}
