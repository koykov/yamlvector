package yamlvector

import (
	"errors"

	"github.com/koykov/vector"
)

var errBadInit = errors.New("bad vector initialization, use jsonvector.NewVector() or jsonvector.Acquire()")

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	_, _ = s, copy
	if !vec.init {
		err = errBadInit
		return
	}
	return vector.ErrNotImplement
}
