package yamlvector

import (
	"errors"
	"unicode/utf8"

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
	root, i := vec.AcquireNodeWithType(0, vector.TypeObject)

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
	// srcp := vec.SrcAddr()
	n := len(src)
	_ = src[n-1]

	for offset < n {
		var ind indent
		ind, vec.indw = vec.indentDW(src, offset, n)
		offset += vec.indw
		if src[offset] == '\t' {
			return offset, ErrBadIndent
		}
		if ind == indentUp {
			return offset, nil
		}

		r, w := utf8.DecodeRune(src[offset:])
		if r == utf8.RuneError {
			return offset, ErrBadUTF8
		}

		switch {
		case node.Type() == vector.TypeUnknown && r != '-':
			node.SetType(vector.TypeObject)
		case node.Type() == vector.TypeUnknown && r == '-':
			node.SetType(vector.TypeArray)
		}

		// todo implement rune based approach

		offset += w
	}

	return offset, err
}

func (vec *Vector) parseGeneric1(depth, offset int, node *vector.Node) (int, error) {
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
	case src[offset] == '"':
		// escaped string
		node.SetType(vector.TypeStr)
		node.Value().SetAddr(srcp, n).SetOffset(offset + 1)
		e := bytealg.IndexByteAtBytes(src, '"', offset+1)
		if e < 0 {
			return n, vector.ErrUnexpEOS
		}
		node.Value().SetBit(flagEscapedString, true) // Always mark string as escaped to avoid double indexing.
		if src[e-1] != '\\' {
			node.Value().SetLen(e - offset - 1)
			offset = e + 1
		} else {
			for i := e; i < n; {
				i = bytealg.IndexByteAtBytes(src, '"', i+1)
				if i < 0 {
					e = n - 1
					break
				}
				e = i
				if src[e-1] != '\\' {
					break
				}
			}
			node.Value().SetLen(e - offset - 1)
			offset = e + 1
		}
	case src[offset] == '|':
		// string block
		i := eot(src, offset)
		node.SetType(vector.TypeString)
		node.Value().InitRaw(srcp, offset, i-offset)
		offset = i
	case src[offset] == '>':
		// foldable string block
		i := eot(src, offset)
		node.SetType(vector.TypeString)
		node.Value().InitRaw(srcp, offset, i-offset).
			SetBit(flagFoldBlock, true)
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
