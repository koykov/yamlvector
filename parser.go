package yamlvector

import (
	"errors"

	"github.com/koykov/vector"
)

var errBadInit = errors.New("bad vector initialization, use yamlvector.NewVector() or yamlvector.Acquire()")

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	_, _ = s, copy
	if !vec.init {
		err = errBadInit
		return
	}
	return vector.ErrNotImplement
}

func (vec *Vector) parseGeneric(depth, offset int, node *vector.Node) (int, error) {
	_, _ = depth, node
	// todo implement me
	return offset, nil
}

func (vec *Vector) parseObject(depth, offset int, node *vector.Node) (int, error) {
	_, _ = depth, node
	// todo implement me
	return offset, nil
}

func (vec *Vector) parseArray(depth, offset int, node *vector.Node) (int, error) {
	_, _ = depth, node
	// todo implement me
	return offset, nil
}
