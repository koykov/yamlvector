package yamlvector

import (
	"errors"
	"unicode"

	"github.com/koykov/bytealg"
	"github.com/koykov/simd/skipline"
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

func (vec *Vector) parseGeneric(depth, offset int, node *vector.Node) (_ int, err error) {
	// todo implement me
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

func (vec *Vector) nextToken() (*token, error) {
	if _, err := vec.skipws(); err != nil {
		return nil, err
	}

	r, w, err := vec.ReadRuneAt(int(vec.pos))
	if err != nil {
		return nil, err
	}
	vec.inccp(w)
	r1, _, err1 := vec.ReadRuneAt(int(vec.pos))
	if err1 != nil {
		return nil, err1
	}
	switch {
	case r == '-' && r1 == ' ':
		// multiline literal
		vec.t.typ = tokenDash
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.inccp(1)
		return &vec.t, nil
	case r == ':':
		vec.t.typ = tokenColon
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.inccp(1)
		return &vec.t, nil
	case r == ',':
		vec.t.typ = tokenComma
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.inccp(1)
		return &vec.t, nil
	case r == '[':
		vec.t.typ = tokenLBracket
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.inccp(1)
		return &vec.t, nil
	case r == ']':
		vec.t.typ = tokenRBracket
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.inccp(1)
		return &vec.t, nil
	case r == '{':
		vec.t.typ = tokenLBrace
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.inccp(1)
		return &vec.t, nil
	case r == '}':
		vec.t.typ = tokenRBrace
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.inccp(1)
		return &vec.t, nil
	case r == '#':
		vec.t.typ = tokenComment
		if _, err = vec.skipws(); err != nil {
			return nil, err
		}
		hi, err2 := vec.readComment()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.inccp(int(hi - vec.pos))
		return &vec.t, nil
	case r == '&':
		vec.t.typ = tokenAnchor
		hi, err2 := vec.readAnchor()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.inccp(int(hi - vec.pos))
		return &vec.t, nil
	case r == '*':
		vec.t.typ = tokenAlias
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.inccp(1)
		return &vec.t, nil
	case r == '!':
		vec.t.typ = tokenTag
		hi, err2 := vec.readTag()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.inccp(int(hi - vec.pos))
		return &vec.t, nil
	case r == '%':
		vec.t.typ = tokenDirective
		hi, err2 := vec.readDirective()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.inccp(int(hi - vec.pos))
		return &vec.t, nil
	case r == '"' || r == '\'':
		vec.t.typ = tokenString
		hi, err2 := vec.readString()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.inccp(int(hi - vec.pos))
		return &vec.t, nil
	case unicode.IsDigit(r) || r == '-' || r == '+' || (r == '.' && unicode.IsDigit(r1)):
		vec.t.typ = tokenNumber
		hi, err2 := vec.readNumber()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.inccp(int(hi - vec.pos))
		return &vec.t, nil
	case unicode.IsLetter(r):
		typ, hi, err2 := vec.readKeyword()
		if err2 != nil {
			return nil, err2
		}
		vec.t.typ = typ
		vec.t.setlo(vec.pos).sethi(hi)
		vec.inccp(int(hi - vec.pos))
		return &vec.t, nil
	case r == '\n' || r == '\r' || (r == '\n' && r1 == '\r'):
		vec.line++
		vec.col = 0
	default:
		return nil, vector.ErrUnexpId
	}
	return nil, vector.ErrUnexpId
}

func (vec *Vector) parseString() error {
	// todo implement me
	return nil
}

func (vec *Vector) readComment() (uint64, error) {
	_, i := skipline.Index2(vec.Src()[vec.pos:])
	vec.pos = vec.pos + uint64(i)
	vec.line++
	vec.col = 0
	return uint64(i), nil
}

func (vec *Vector) readAnchor() (uint64, error) {
	// todo implement me
	return 0, nil
}

func (vec *Vector) readTag() (uint64, error) {
	// todo implement me
	return 0, nil
}

func (vec *Vector) readDirective() (uint64, error) {
	// todo implement me
	return 0, nil
}

func (vec *Vector) readString() (uint64, error) {
	// todo implement me
	return 0, nil
}

func (vec *Vector) readNumber() (uint64, error) {
	// todo implement me
	return 0, nil
}

func (vec *Vector) readKeyword() (ttoken, uint64, error) {
	// todo implement me
	return tokenNull, 0, nil
}
