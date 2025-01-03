package yamlvector

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
	init bool
	indw int
}

func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

func (vec *Vector) ParseStr(s string) error {
	return vec.parse(byteconv.S2B(s), false)
}

func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(byteconv.S2B(s), true)
}

func NewVector() *Vector {
	vec := &Vector{init: true}
	// todo implement helper.
	vec.Helper = nil
	return vec
}
