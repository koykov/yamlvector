package yamlvector

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
	init bool
	indw int
	ind  int
	inds indent

	lineFC, isFold, isLit, IsRawFold bool
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
	vec.ind = 0
	vec.lineFC = false
	vec.isFold = false
	vec.isLit = false
	vec.IsRawFold = false
}

func (vec *Vector) isDoc() bool {
	return vec.isFold || vec.isLit || vec.IsRawFold
}

func (vec *Vector) indent(r rune) {
	if vec.lineFC && (r == '\n' || r == '\r') && vec.isDoc() {
		return
	}
	if vec.lineFC && r == ' ' {
		vec.ind++
		return
	}
	if !vec.lineFC {
		vec.inds = indentEqual
		return
	}
	// todo update indent level
	// todo update indent state
	vec.lineFC = false
}
