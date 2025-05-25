package yamlvector

import (
	"errors"
	"unicode"

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
	if err := vec.skipws(); err != nil {
		return nil, err
	}

	r, w, err := vec.ReadRuneAt(int(vec.pos))
	if err != nil {
		return nil, err
	}
	switch r {
	case '-':
		// multiline literal
		// todo implement me
	case ':':
		vec.t.typ = tokenColon
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		return &vec.t, nil
	case ',':
		vec.t.typ = tokenComma
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		return &vec.t, nil
	case '[':
		vec.t.typ = tokenLBracket
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		return &vec.t, nil
	case ']':
		vec.t.typ = tokenRBracket
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		return &vec.t, nil
	case '{':
		vec.t.typ = tokenLBrace
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		return &vec.t, nil
	case '}':
		vec.t.typ = tokenRBrace
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		return &vec.t, nil
	case '#':
		vec.t.typ = tokenComment
		hi := vec.readComment()
		vec.t.setlo(vec.pos).sethi(hi)
		return &vec.t, nil
	case '&':
		vec.t.typ = tokenAnchor
		hi := vec.readAnchor()
		vec.t.setlo(vec.pos).sethi(hi)
		return &vec.t, nil
	case '*':
		vec.t.typ = tokenAlias
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		return &vec.t, nil
	case '!':
		vec.t.typ = tokenTag
		hi := vec.readTag()
		vec.t.setlo(vec.pos).sethi(hi)
		return &vec.t, nil
	case '%':
		vec.t.typ = tokenDirective
		hi := vec.readDirective()
		vec.t.setlo(vec.pos).sethi(hi)
		return &vec.t, nil
	case '"', '\'':
		vec.t.typ = tokenString
		hi := vec.readString()
		vec.t.setlo(vec.pos).sethi(hi)
		return &vec.t, nil
	default:
		if unicode.IsLetter(r) || r == '+' || r == '.' {
			vec.t.typ = tokenNumber
			hi := vec.readNumber()
			vec.t.setlo(vec.pos).sethi(hi)
			return &vec.t, nil
		}

		if unicode.IsLetter(r) {
			typ, hi := vec.readKeyword()
			vec.t.typ = typ
			vec.t.setlo(vec.pos).sethi(hi)
			return &vec.t, nil
		}
	}
	return nil, nil
}

func (vec *Vector) parseString() error {
	// todo implement me
	return nil
}

func (vec *Vector) skipws() error {
	for {
		r, w, err := vec.ReadRuneAt(int(vec.pos))
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			break
		}
		vec.pos += uint(w)
	}
	return nil
}

func (vec *Vector) skipl() {
	i := bytealg.IndexByteAtBytes(vec.Src(), '\n', int(vec.pos))
	if i < 0 {
		return
	}
	vec.pos += uint(i + 1)
}
