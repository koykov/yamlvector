package yamlvector

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
	init bool
	indw int

	isFold, isLit, IsRawFold bool
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

func (vec *Vector) Reset() {
	vec.Vector.Reset()
	vec.init = false
	vec.indw = 0
	vec.isFold = false
	vec.isLit = false
	vec.IsRawFold = false
}

func (vec *Vector) isDoc() bool {
	return vec.isFold || vec.isLit || vec.IsRawFold
}
