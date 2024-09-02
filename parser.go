package yamlvector

import (
	"errors"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

var errBadInit = errors.New("bad vector initialization, use yamlvector.NewVector() or yamlvector.Acquire()")

var bBools = []byte("truefalse")

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	if !vec.init {
		err = errBadInit
		return
	}

	s = bytealg.TrimBytesFmt4(s)
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}

	offset := 0
	// Create root node and register it.
	root, i := vec.AcquireNode(0)

	// Parse source data.
	if offset, err = vec.parseGeneric(0, offset, root); err != nil {
		vec.SetErrOffset(offset)
		return err
	}
	vec.ReleaseNode(i, root)

	// Check unparsed tail.
	if offset < vec.SrcLen() {
		vec.SetErrOffset(offset)
		return vector.ErrUnparsedTail
	}

	return
}

func (vec *Vector) parseGeneric(depth, offset int, node *vector.Node) (int, error) {
	var err error
	node.SetOffset(vec.Index.Len(depth))
	src := vec.Src()
	srcp := vec.SrcAddr()
	n := len(src)
	_ = src[n-1]

	var (
		typ vector.Type
		bv  bool
	)

	switch {
	case ensureNullOrBool(src, &offset, &typ, &bv):
		node.SetType(typ)
		if typ == vector.TypeBool {
			if bv {
				node.Value().Init(bBools, 0, 4)
			} else {
				node.Value().Init(bBools, 4, 5)
			}
		}
	case ensureDigit(src[offset]):
		i := offset
		for ensureDigit(src[i]) {
			i++
			if i == n {
				break
			}
		}
		node.SetType(vector.TypeNumber)
		node.Value().InitRaw(srcp, offset, i-offset)
		offset = i
	default:
		// raw string case
		i := eol(src, offset)
		node.SetType(vector.TypeString)
		node.Value().InitRaw(srcp, offset, i-offset)
		offset = i
	}
	return offset, err
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
