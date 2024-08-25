package yamlvector

import (
	"errors"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

var errBadInit = errors.New("bad vector initialization, use yamlvector.NewVector() or yamlvector.Acquire()")

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
	_ = srcp
	n := len(src)
	_ = src[n-1]
	switch {
	case vec.ensureTrue(src, &offset):
		// todo implement me
	case vec.ensureFalse(src, &offset):
		// todo implement me
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

func (vec *Vector) ensureTrue(src []byte, poffset *int) (ok bool) {
	bin1 := bin(src, *poffset, 1)
	bin2 := bin(src, *poffset, 2)
	bin3 := bin(src, *poffset, 3)
	bin4 := bin(src, *poffset, 4)
	switch {
	case bin1 == binBoolTrue[0] || bin1 == binBoolTrue[1]:
		*poffset += 1
		ok = true
	case bin2 == binBoolTrue[2] || bin2 == binBoolTrue[3] || bin2 == binBoolTrue[4]:
		*poffset += 2
		ok = true
	case bin3 == binBoolTrue[5] || bin3 == binBoolTrue[6] || bin3 == binBoolTrue[7]:
		*poffset += 3
		ok = true
	case bin4 == binBoolTrue[8] || bin4 == binBoolTrue[9] || bin4 == binBoolTrue[10]:
		*poffset += 4
		ok = true
	}
	return
}

func (vec *Vector) ensureFalse(src []byte, poffset *int) (ok bool) {
	bin1 := bin(src, *poffset, 1)
	bin2 := bin(src, *poffset, 2)
	bin3 := bin(src, *poffset, 3)
	bin5 := bin(src, *poffset, 4)
	switch {
	case bin1 == binBoolFalse[0] || bin1 == binBoolFalse[1]:
		*poffset += 1
		ok = true
	case bin2 == binBoolFalse[2] || bin2 == binBoolFalse[3] || bin2 == binBoolFalse[4]:
		*poffset += 2
		ok = true
	case bin3 == binBoolFalse[5] || bin3 == binBoolFalse[6] || bin3 == binBoolFalse[7]:
		*poffset += 3
		ok = true
	case bin5 == binBoolFalse[8] || bin5 == binBoolFalse[9] || bin5 == binBoolFalse[10]:
		*poffset += 5
		ok = true
	}
	return
}
