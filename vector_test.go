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
	t.Run("scalar_number", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", 123456)
	})
	t.Run("scalar_number_float", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", 3.1415)
	})
	t.Run("scalar_string", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
	})
	t.Run("scalar_string_no_fmt", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
	})
	t.Run("scalar_string_keep_fmt", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
	})
	t.Run("scalar_string_escape", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
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
	b.Run("scalar_number", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeNumber)
		})
	})
	b.Run("scalar_number_float", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeNumber)
		})
	})
	b.Run("scalar_string", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeString)
		})
	})
	b.Run("scalar_string_no_fmt", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeString)
		})
	})
	b.Run("scalar_string_keep_fmt", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeString)
		})
	})
	b.Run("scalar_string_escape", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeString)
		})
	})
}
