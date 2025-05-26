package yamlvector

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
	init bool

	t    token
	pos  uint64
	line uint64
	col  uint64
	pcol uint64
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

	vec.pos = 0
	vec.line = 0
	vec.col = 0
	vec.pcol = 0
}

func (vec *Vector) incp(d int) *Vector {
	vec.pos += uint64(d)
	return vec
}

func (vec *Vector) incc(d int) *Vector {
	vec.col += uint64(d)
	return vec
}

func (vec *Vector) inccp(d int) *Vector {
	vec.pos += uint64(d)
	vec.col += uint64(d)
	return vec
}

func (vec *Vector) incl(d int) *Vector {
	vec.line += uint64(d)
	return vec
}
