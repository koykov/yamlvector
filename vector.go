package yamlvector

import (
	"io"

	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector

	t       token
	pos     uint64
	anchors map[string]int
}

func NewVector() *Vector {
	vec := &Vector{}
	vec.SetBit(vector.FlagInit, true)
	vec.Helper = helper
	return vec
}

func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

func (vec *Vector) ParseString(s string) error {
	return vec.parse(byteconv.S2B(s), false)
}

func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

func (vec *Vector) ParseCopyString(s string) error {
	return vec.parse(byteconv.S2B(s), true)
}

func (vec *Vector) ParseFile(path string) error {
	err := vec.Vector.ParseFile(path)
	if err != vector.ErrNotImplement {
		return err
	}
	return vec.parse(vec.Buf(), false)
}

func (vec *Vector) ParseReader(r io.Reader) error {
	err := vec.Vector.ParseReader(r)
	if err != vector.ErrNotImplement {
		return err
	}
	return vec.parse(vec.Buf(), false)
}

var _ vector.Interface = (*Vector)(nil)

func (vec *Vector) Reset() {
	vec.Vector.Reset()

	vec.t = token{}
	vec.pos = 0
	clear(vec.anchors)
}

func (vec *Vector) offmove(delta uint64) *Vector {
	vec.pos += uint64(delta)
	return vec
}
