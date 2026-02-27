package yamlvector

import (
	"math"
	"testing"

	"github.com/koykov/vector"
)

func TestScalar(t *testing.T) {
	vec := NewVector()
	t.Run("comment", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
	})
	t.Run("scalar_null", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNull)
	})
	t.Run("scalar_Null", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNull)
	})
	t.Run("scalar_NULL", func(t *testing.T) {
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
	t.Run("scalar_Bool", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeBool)
		assertBool(t, vec, "", true)
	})
	t.Run("scalar_BOOL", func(t *testing.T) {
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
	t.Run("scalar_number_neg", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", -123456)
	})
	t.Run("scalar_number_float", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", 3.1415)
	})
	t.Run("scalar_number_exponential", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", 1200000000)
	})
	t.Run("scalar_number_exponential_neg", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", -0.00000314)
	})
	t.Run("scalar_number_fake", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
		assertString(t, vec, "", "123456#789")
	})
	t.Run("scalar_number_inf", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", math.Inf(1))
	})
	t.Run("scalar_number_Inf", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", math.Inf(1))
	})
	t.Run("scalar_number_INF", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNumber)
		assertNumber(t, vec, "", math.Inf(1))
	})
	t.Run("scalar_string", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
		assertString(t, vec, "", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec eu magna nec felis ullamcorper blandit. Aliquam laoreet sodales massa sit amet porta. Cras mattis ornare faucibus.")
	})
	t.Run("scalar_string_folded_block", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
		assertString(t, vec, "", "Wrapped text will be folded into a single paragraph \nBlank lines denote paragraph breaks")
	})
	t.Run("scalar_string_literal_block", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
		assertString(t, vec, "", "There was a young fellow of Warwick\nWho had reason for feeling euphoric\nFor he could, by election\nHave triune erection\nIonic, Corinthian, and Doric")
	})
	t.Run("scalar_string_escape", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeString)
		assertString(t, vec, "", "this is \\\"escaped\\\" string\\n")
	})
}

func TestObject(t *testing.T) {
	// vec := NewVector()
	// t.Run("object_multi_line", func(t *testing.T) {
	// 	vec = assertParse(t, vec, nil, 0)
	// 	assertType(t, vec, "", vector.TypeObject)
	// })
}

func BenchmarkScalar(b *testing.B) {
	b.Run("scalar_null", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeNull) }) })
	b.Run("scalar_null_canonical", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeNull) }) })
	b.Run("scalar_null_none", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeNull) }) })
	b.Run("scalar_bool", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeBool) }) })
	b.Run("scalar_bool_low", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeBool) }) })
	b.Run("scalar_bool_on", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeBool) }) })
	b.Run("scalar_number", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeNumber) }) })
	b.Run("scalar_number_float", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeNumber) }) })
	b.Run("scalar_string", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeString) }) })
	b.Run("scalar_string_no_fmt", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeString) }) })
	b.Run("scalar_string_keep_fmt", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeString) }) })
	b.Run("scalar_string_escape", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeString) }) })
}
