package yamlvector

import (
	"bytes"
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

	// Create root node and register it.
	root, i := vec.AcquireNodeWithType(0, vector.TypeObject)

	// Parse source data.
	if err = vec.parseGeneric(0, root); err != nil {
		vec.SetErrOffset(int(vec.pos))
		return err
	}
	vec.ReleaseNode(i, root)

	// Check unparsed tail.
	if vec.pos < uint64(vec.SrcLen()) {
		vec.SetErrOffset(int(vec.pos))
		return vector.ErrUnparsedTail
	}

	return
}

func (vec *Vector) parseGeneric(depth int, node *vector.Node) error {
	for {
		t, err := vec.nextToken()
		if err != nil {
			return err
		}
		switch t.typ {
		case tokenComment:
			// do nothing
		case tokenEOF:
			return nil
		case tokenDash:
			err = vec.parseObject(depth, node)
		case tokenColon:
			err = vec.parseArray(depth, node)
		case tokenComma:
			err = vec.parseGeneric(depth, node)
		default:
			return vector.ErrUnexpId
		}
	}
}

func (vec *Vector) parseObject(depth int, node *vector.Node) error {
	_, _ = depth, node
	// todo implement me
	return nil
}

func (vec *Vector) parseArray(depth int, node *vector.Node) error {
	_, _ = depth, node
	// todo implement me
	return nil
}

func (vec *Vector) nextToken() (*token, error) {
	if _, err := vec.skipws(); err != nil {
		return nil, err
	}

	if vec.pos >= uint64(vec.SrcLen()) {
		vec.t.typ = tokenEOF
		return &vec.t, nil
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
	b := vec.Src()[vec.pos:]
	i, j := skipline.Index2(b)
	if i == -1 {
		j = len(b)
	}
	vec.pos = vec.pos + uint64(j)
	vec.line++
	vec.col = 0
	return vec.pos, nil
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
	off := vec.pos
	i, j := skipline.Index2(vec.Src()[off:])
	if i == -1 {
		i = vec.SrcLen()
		j = i
	}
	vec.pos = uint64(j)
	v := vec.Src()[off:i]
	switch {
	case bytes.Equal(v, bnull) || bytes.Equal(v, bNull) || bytes.Equal(v, bNULL), bytes.Equal(v, bNone) || bytes.Equal(v, bTilda):
		return tokenNull, vec.pos, nil
	case bytes.Equal(v, btrue) || bytes.Equal(v, bTrue) || bytes.Equal(v, bTRUE) || bytes.Equal(v, bOn),
		bytes.Equal(v, bfalse) || bytes.Equal(v, bFalse) || bytes.Equal(v, bFALSE) || bytes.Equal(v, bOff):
		return tokenBool, vec.pos, nil
	default:
		return tokenString, vec.pos, nil
	}
}

var (
	bnull  = []byte("null")
	bNull  = []byte("Null")
	bNULL  = []byte("NULL")
	bNone  = []byte("None")
	bTilda = []byte("~")
	btrue  = []byte("true")
	bTrue  = []byte("True")
	bTRUE  = []byte("TRUE")
	bfalse = []byte("false")
	bFalse = []byte("False")
	bFALSE = []byte("FALSE")
	bOn    = []byte("On")
	bOff   = []byte("Off")
)
