package yamlvector

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
	"github.com/koykov/simd/skipline"
	"github.com/koykov/vector"
)

var errBadInit = errors.New("bad vector initialization, use yamlvector.NewVector() or yamlvector.Acquire()")

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	if !vec.CheckBit(vector.FlagInit) {
		err = errBadInit
		return
	}

	s = bytealg.TrimBytesFmt4(s)
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}
	vec.pos = 0

	for {
		vec.skipBlank()
		if vec.eof() {
			break
		}
		if vec.isDocMarker() {
			vec.skipln()
			continue
		}
		start := vec.pos
		indent := vec.lineIndent()

		// Create root node and register it.
		root, i := vec.AcquireNode(0)
		if err = vec.parseGeneric(0, root, indent); err != nil {
			vec.SetErrOffset(int(vec.pos))
			return err
		}
		vec.ReleaseNode(i, root)

		// Guard against non-advancing parsers.
		if vec.pos == start && !vec.eof() {
			vec.SetErrOffset(int(vec.pos))
			return vector.ErrUnexpId
		}
	}
	return
}

func (vec *Vector) parseGeneric(depth int, node *vector.Node, indent int) error {
	vec.skipBlank()
	if vec.eof() {
		node.SetType(vector.TypeNull)
		return nil
	}
	if vec.lineIndent() < indent {
		node.SetType(vector.TypeNull)
		return nil
	}
	switch c := vec.Src()[vec.pos]; {
	case c == '[':
		return vec.parseFlowArray(depth, node)
	case c == '{':
		return vec.parseFlowObject(depth, node)
	case c == '&':
		if err := vec.readAnchorNode(node); err != nil {
			return err
		}
		vec.skipWS()
		if vec.eol() {
			vec.skipln()
			return vec.parseNested(depth, node, indent)
		}
		return vec.parseInline(depth, node, indent)
	case c == '!':
		vec.readTagName()
		vec.skipWS()
		if vec.eol() {
			vec.skipln()
			return vec.parseNested(depth, node, indent)
		}
		return vec.parseInline(depth, node, indent)
	case c == '*':
		return vec.parseAlias(node)
	case c == '"' || c == '\'':
		return vec.parseQuoted(node)
	case c == '|' || c == '>':
		return vec.parseBlockScalar(node)
	case c == '-' && vec.isDash():
		return vec.parseBlockSeq(depth, node, indent)
	default:
		if vec.isMapStart(int(vec.pos)) {
			return vec.parseBlockMap(depth, node, indent)
		}
		return vec.parsePlainScalar(node)
	}
}

// parseBlockMap parses a block mapping whose keys are aligned at given indentation.
func (vec *Vector) parseBlockMap(depth int, node *vector.Node, indent int) error {
	node.SetType(vector.TypeObject)
	node.SetOffset(vec.Index.Len(depth + 1))
	for first := true; ; first = false {
		if !first {
			vec.skipBlank()
			if vec.eof() || vec.lineIndent() != indent {
				break
			}
		}
		if vec.eof() {
			break
		}
		if vec.Src()[vec.pos] == '#' {
			vec.skipln()
			continue
		}
		child, i := vec.AcquireChildWithType(node, depth+1, vector.TypeUnknown)
		if err := vec.readKey(child); err != nil {
			vec.ReleaseNode(i, child)
			return err
		}
		vec.skipWS()
		if vec.eof() || vec.Src()[vec.pos] != ':' {
			vec.ReleaseNode(i, child)
			return vector.ErrUnexpId
		}
		vec.pos++
		vec.skipWS()
		var err error
		if vec.eol() {
			vec.skipln()
			err = vec.parseNested(depth+1, child, indent)
		} else {
			err = vec.parseInline(depth+1, child, indent)
		}
		vec.ReleaseNode(i, child)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseBlockSeq parses a block sequence (dash items) aligned at given indentation.
func (vec *Vector) parseBlockSeq(depth int, node *vector.Node, indent int) error {
	node.SetType(vector.TypeArray)
	node.SetOffset(vec.Index.Len(depth + 1))
	for {
		vec.skipBlank()
		if vec.eof() || vec.lineIndent() != indent || !vec.isDash() {
			break
		}
		vec.pos++
		child, i := vec.AcquireChildWithType(node, depth+1, vector.TypeUnknown)
		vec.skipWS()
		var err error
		if vec.eol() {
			vec.skipln()
			err = vec.parseNested(depth+1, child, indent)
		} else if col := int(vec.pos) - vec.lineStart(); vec.isMapStart(int(vec.pos)) {
			err = vec.parseBlockMap(depth+1, child, col)
		} else {
			err = vec.parseInline(depth+1, child, indent)
		}
		vec.ReleaseNode(i, child)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseNested parses a value placed on the next line(s) with deeper indentation.
func (vec *Vector) parseNested(depth int, node *vector.Node, parentIndent int) error {
	vec.skipBlank()
	if vec.eof() {
		node.SetType(vector.TypeNull)
		return nil
	}
	if ind := vec.lineIndent(); ind > parentIndent {
		return vec.parseGeneric(depth, node, ind)
	} else if ind == parentIndent && vec.isDash() {
		// Block sequence may be placed at the same indentation as its key.
		return vec.parseGeneric(depth, node, ind)
	}
	node.SetType(vector.TypeNull)
	return nil
}

// parseInline parses a value that starts on the current line.
func (vec *Vector) parseInline(depth int, node *vector.Node, parentIndent int) error {
	vec.skipWS()
	switch c := vec.Src()[vec.pos]; {
	case c == '[':
		return vec.parseFlowArray(depth, node)
	case c == '{':
		return vec.parseFlowObject(depth, node)
	case c == '|' || c == '>':
		return vec.parseBlockScalar(node)
	case c == '*':
		return vec.parseAlias(node)
	case c == '"' || c == '\'':
		return vec.parseQuoted(node)
	case c == '&':
		if err := vec.readAnchorNode(node); err != nil {
			return err
		}
		vec.skipWS()
		if vec.eol() {
			vec.skipln()
			return vec.parseNested(depth, node, parentIndent)
		}
		return vec.parseInline(depth, node, parentIndent)
	default:
		return vec.parsePlainScalar(node)
	}
}

// parseFlowArray parses a flow sequence: [a, b, c].
func (vec *Vector) parseFlowArray(depth int, node *vector.Node) error {
	node.SetType(vector.TypeArray)
	node.SetOffset(vec.Index.Len(depth + 1))
	vec.pos++
	for {
		vec.skipWSn()
		if vec.eof() {
			return vector.ErrUnexpEOF
		}
		if vec.Src()[vec.pos] == ']' {
			vec.pos++
			return nil
		}
		child, i := vec.AcquireChildWithType(node, depth+1, vector.TypeUnknown)
		err := vec.parseFlowValue(depth+1, child)
		vec.ReleaseNode(i, child)
		if err != nil {
			return err
		}
		vec.skipWSn()
		if vec.eof() {
			return vector.ErrUnexpEOF
		}
		switch vec.Src()[vec.pos] {
		case ',':
			vec.pos++
		case ']':
			vec.pos++
			return nil
		default:
			return vector.ErrUnexpId
		}
	}
}

// parseFlowObject parses a flow mapping: {a: b, c: d}.
func (vec *Vector) parseFlowObject(depth int, node *vector.Node) error {
	node.SetType(vector.TypeObject)
	node.SetOffset(vec.Index.Len(depth + 1))
	vec.pos++
	for {
		vec.skipWSn()
		if vec.eof() {
			return vector.ErrUnexpEOF
		}
		if vec.Src()[vec.pos] == '}' {
			vec.pos++
			return nil
		}
		child, i := vec.AcquireChildWithType(node, depth+1, vector.TypeUnknown)
		err := vec.readFlowKey(child)
		if err == nil {
			vec.skipWSn()
			if vec.eof() || vec.Src()[vec.pos] != ':' {
				err = vector.ErrUnexpId
			} else {
				vec.pos++
				vec.skipWSn()
				err = vec.parseFlowValue(depth+1, child)
			}
		}
		vec.ReleaseNode(i, child)
		if err != nil {
			return err
		}
		vec.skipWSn()
		if vec.eof() {
			return vector.ErrUnexpEOF
		}
		switch vec.Src()[vec.pos] {
		case ',':
			vec.pos++
		case '}':
			vec.pos++
			return nil
		default:
			return vector.ErrUnexpId
		}
	}
}

// parseFlowValue parses a single flow value (any type).
func (vec *Vector) parseFlowValue(depth int, node *vector.Node) error {
	vec.skipWSn()
	if vec.eof() {
		return vector.ErrUnexpEOF
	}
	switch c := vec.Src()[vec.pos]; {
	case c == '"' || c == '\'':
		return vec.parseQuoted(node)
	case c == '[':
		return vec.parseFlowArray(depth, node)
	case c == '{':
		return vec.parseFlowObject(depth, node)
	case c == '*':
		return vec.parseAlias(node)
	case c == '&':
		if err := vec.readAnchorNode(node); err != nil {
			return err
		}
		vec.skipWSn()
		return vec.parseFlowValue(depth, node)
	default:
		return vec.parseFlowScalar(node)
	}
}

// parseQuoted reads a quoted string node.
func (vec *Vector) parseQuoted(node *vector.Node) error {
	t, err := vec.nextToken()
	if err != nil {
		return err
	}
	node.SetType(vector.TypeString)
	node.Value().SetAddr(vec.SrcAddr(), vec.SrcLen()).SetOffset(int(t.lo)).SetLen(int(t.hi - t.lo))
	node.Value().SetBit(flagEscapedString, true)
	return nil
}

// parseBlockScalar reads a literal (|) or folded (>) block scalar node.
func (vec *Vector) parseBlockScalar(node *vector.Node) error {
	t, err := vec.nextToken()
	if err != nil {
		return err
	}
	node.SetType(vector.TypeString)
	node.Value().SetAddr(vec.SrcAddr(), vec.SrcLen()).SetOffset(int(t.lo)).SetLen(int(t.hi - t.lo))
	return nil
}

// parseAlias reads an alias (*name) node and links it to the anchored node.
func (vec *Vector) parseAlias(node *vector.Node) error {
	vec.pos++
	start := vec.pos
	vec.readAnchorName()
	name := byteconv.B2S(vec.Src()[start:vec.pos])
	if idx, ok := vec.anchors[name]; ok {
		vec.cloneNode(node, vec.NodeAt(idx))
	} else {
		node.SetType(vector.TypeNull)
	}
	vec.skipl()
	return nil
}

// cloneNode makes node a shallow clone of target keeping node's key.
func (vec *Vector) cloneNode(node, target *vector.Node) {
	key := *node.Key()
	*node = *target
	*node.Key() = key
}

// parsePlainScalar reads a plain (unquoted) scalar node.
func (vec *Vector) parsePlainScalar(node *vector.Node) error {
	srcp := vec.SrcAddr()
	n := vec.SrcLen()
	start := int(vec.pos)
	for !vec.eof() {
		c := vec.Src()[vec.pos]
		if c == '\n' || c == '\r' {
			break
		}
		if c == '#' && vec.pos > uint64(start) {
			if pc := vec.Src()[vec.pos-1]; pc == ' ' || pc == '\t' {
				break
			}
		}
		vec.pos++
	}
	end := int(vec.pos)
	for end > start && (vec.Src()[end-1] == ' ' || vec.Src()[end-1] == '\t') {
		end--
	}
	vec.checktyp(node, srcp, n, start, end)
	return nil
}

// parseFlowScalar reads a plain scalar in the flow context.
func (vec *Vector) parseFlowScalar(node *vector.Node) error {
	srcp := vec.SrcAddr()
	n := vec.SrcLen()
	start := int(vec.pos)
	for !vec.eof() {
		c := vec.Src()[vec.pos]
		if c == ',' || c == ']' || c == '}' || c == '\n' || c == '\r' {
			break
		}
		if c == ':' {
			j := int(vec.pos) + 1
			if j >= vec.SrcLen() || vec.Src()[j] == ' ' || vec.Src()[j] == ',' || vec.Src()[j] == ']' || vec.Src()[j] == '}' {
				break
			}
		}
		if c == '#' && vec.pos > uint64(start) {
			if pc := vec.Src()[vec.pos-1]; pc == ' ' || pc == '\t' {
				break
			}
		}
		vec.pos++
	}
	end := int(vec.pos)
	for end > start && (vec.Src()[end-1] == ' ' || vec.Src()[end-1] == '\t') {
		end--
	}
	vec.checktyp(node, srcp, n, start, end)
	return nil
}

// checktyp sets node type and value span for a plain scalar.
func (vec *Vector) checktyp(node *vector.Node, srcp uintptr, n, lo, hi int) {
	b := vec.Src()[lo:hi]
	switch {
	case bytes.Equal(b, bnull) || bytes.Equal(b, bNull) || bytes.Equal(b, bNULL) ||
		bytes.Equal(b, bNone) || bytes.Equal(b, bTilda):
		node.SetType(vector.TypeNull)
	case bytes.Equal(b, btrue) || bytes.Equal(b, bTrue) || bytes.Equal(b, bTRUE) ||
		bytes.Equal(b, bOn) || bytes.Equal(b, bon) || bytes.Equal(b, bON):
		node.SetType(vector.TypeBool)
		node.Value().SetAddr(srcp, n).SetOffset(lo).SetLen(hi-lo).SetBit(vector.FlagExtraBool, true)
	case bytes.Equal(b, bfalse) || bytes.Equal(b, bFalse) || bytes.Equal(b, bFALSE) ||
		bytes.Equal(b, bOff) || bytes.Equal(b, boff) || bytes.Equal(b, bOFF):
		node.SetType(vector.TypeBool)
		node.Value().SetAddr(srcp, n).SetOffset(lo).SetLen(hi - lo)
	case bytes.Equal(b, binf) || bytes.Equal(b, bInf) || bytes.Equal(b, bINF):
		node.SetType(vector.TypeNumber)
		node.Value().SetAddr(pbnums, 10).SetOffset(0).SetLen(3)
	case bytes.Equal(b, bninf) || bytes.Equal(b, bnInf) || bytes.Equal(b, bnINF):
		node.SetType(vector.TypeNumber)
		node.Value().SetAddr(pbnums, 10).SetOffset(3).SetLen(4)
	case bytes.Equal(b, bnan) || bytes.Equal(b, bNaN) || bytes.Equal(b, bNAN):
		node.SetType(vector.TypeNumber)
		node.Value().SetAddr(pbnums, 10).SetOffset(7).SetLen(3)
	default:
		if isNumeric(b) {
			node.SetType(vector.TypeNumber)
		} else {
			node.SetType(vector.TypeString)
		}
		node.Value().SetAddr(srcp, n).SetOffset(lo).SetLen(hi - lo)
	}
}

// readKey fills child key from the source at current position (block context).
func (vec *Vector) readKey(child *vector.Node) error {
	srcp := vec.SrcAddr()
	n := vec.SrcLen()
	if c := vec.Src()[vec.pos]; c == '"' || c == '\'' {
		t, err := vec.nextToken()
		if err != nil {
			return err
		}
		child.Key().SetAddr(srcp, n).SetOffset(int(t.lo)).SetLen(int(t.hi - t.lo))
		child.Key().SetBit(flagEscapedString, true)
		return nil
	}
	start := int(vec.pos)
	for !vec.eof() {
		c := vec.Src()[vec.pos]
		if c == '\n' || c == '\r' {
			break
		}
		if c == ':' && vec.isSepAt(int(vec.pos)+1) {
			break
		}
		vec.pos++
	}
	end := int(vec.pos)
	for end > start && (vec.Src()[end-1] == ' ' || vec.Src()[end-1] == '\t') {
		end--
	}
	child.Key().SetAddr(srcp, n).SetOffset(start).SetLen(end - start)
	child.Key().SetBit(flagEscapedString, true)
	return nil
}

// readFlowKey fills child key from the source at current position (flow context).
func (vec *Vector) readFlowKey(child *vector.Node) error {
	srcp := vec.SrcAddr()
	n := vec.SrcLen()
	if c := vec.Src()[vec.pos]; c == '"' || c == '\'' {
		t, err := vec.nextToken()
		if err != nil {
			return err
		}
		child.Key().SetAddr(srcp, n).SetOffset(int(t.lo)).SetLen(int(t.hi - t.lo))
		child.Key().SetBit(flagEscapedString, true)
		return nil
	}
	start := int(vec.pos)
	for !vec.eof() {
		c := vec.Src()[vec.pos]
		if c == ':' || c == ',' || c == '}' || c == '\n' || c == '\r' {
			break
		}
		vec.pos++
	}
	end := int(vec.pos)
	for end > start && (vec.Src()[end-1] == ' ' || vec.Src()[end-1] == '\t') {
		end--
	}
	child.Key().SetAddr(srcp, n).SetOffset(start).SetLen(end - start)
	child.Key().SetBit(flagEscapedString, true)
	return nil
}

// readAnchorNode reads an anchor (&name) and registers it for the given node.
func (vec *Vector) readAnchorNode(node *vector.Node) error {
	vec.pos++
	start := vec.pos
	vec.readAnchorName()
	if vec.pos == start {
		return vector.ErrUnexpId
	}
	if vec.anchors == nil {
		vec.anchors = make(map[string]int)
	}
	vec.anchors[byteconv.B2S(vec.Src()[start:vec.pos])] = node.Index()
	return nil
}

// readAnchorName advances position through anchor/alias name characters.
func (vec *Vector) readAnchorName() {
	for !vec.eof() {
		r, w := utf8.DecodeRune(vec.Src()[vec.pos:])
		if !isAnchorChar(r) {
			break
		}
		vec.pos += uint64(w)
	}
}

// readTagName advances position through a tag (!!str, !custom) token.
func (vec *Vector) readTagName() {
	vec.pos++
	for !vec.eof() {
		c := vec.Src()[vec.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			break
		}
		vec.pos++
	}
}

// isMapStart reports whether the current line contains a block mapping key separator.
func (vec *Vector) isMapStart(p int) bool {
	s := vec.Src()
	quote := byte(0)
	depth := 0
	for i := p; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			if c == '\\' && quote == '"' {
				i++
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '\n', '\r':
			return false
		case '#':
			if i > p && (s[i-1] == ' ' || s[i-1] == '\t') {
				return false
			}
		case '"', '\'':
			quote = c
		case '[', '{':
			depth++
		case ']', '}':
			depth--
		case ':':
			if depth == 0 && vec.isSepAt(i+1) {
				return true
			}
		}
	}
	return false
}

// isSepAt reports whether position i is a separator (space, tab, EOL or end of source).
func (vec *Vector) isSepAt(i int) bool {
	if i >= vec.SrcLen() {
		return true
	}
	switch vec.Src()[i] {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

// isDocMarker reports whether current position starts a document marker (--- or ...).
func (vec *Vector) isDocMarker() bool {
	s := vec.Src()
	i := int(vec.pos)
	if i+3 > len(s) {
		return false
	}
	if !(bytes.Equal(s[i:i+3], bdocStart) || bytes.Equal(s[i:i+3], bdocEnd)) {
		return false
	}
	if i+3 == len(s) {
		return true
	}
	switch s[i+3] {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

// isDash reports whether current position starts a block sequence item.
func (vec *Vector) isDash() bool {
	if vec.Src()[vec.pos] != '-' {
		return false
	}
	i := int(vec.pos) + 1
	if i >= vec.SrcLen() {
		return true
	}
	switch vec.Src()[i] {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

// eof reports whether position reached the end of source.
func (vec *Vector) eof() bool {
	return vec.pos >= uint64(vec.SrcLen())
}

// eol reports whether current position is a line end or an inline comment.
func (vec *Vector) eol() bool {
	if vec.eof() {
		return true
	}
	switch vec.Src()[vec.pos] {
	case '\n', '\r', '#':
		return true
	}
	return false
}

// skipWS skips spaces and tabs.
func (vec *Vector) skipWS() {
	for !vec.eof() {
		switch vec.Src()[vec.pos] {
		case ' ', '\t':
			vec.pos++
		default:
			return
		}
	}
}

// skipWSn skips spaces and line breaks inside flow collections.
func (vec *Vector) skipWSn() {
	for !vec.eof() {
		switch vec.Src()[vec.pos] {
		case ' ', '\t', '\n', '\r':
			vec.pos++
		case '#':
			vec.skipln()
		default:
			return
		}
	}
}

// skipl moves position to nearest EOL (newline is not consumed).
func (vec *Vector) skipl() {
	for !vec.eof() && vec.Src()[vec.pos] != '\n' {
		vec.pos++
	}
}

// skipln moves position to nearest EOL (including newline).
func (vec *Vector) skipln() {
	vec.skipl()
	if !vec.eof() {
		vec.pos++
	}
}

// skipBlank skips spaces, blank lines, comments and document directives.
func (vec *Vector) skipBlank() {
	s := vec.Src()
	for {
		for vec.pos < uint64(len(s)) && (s[vec.pos] == ' ' || s[vec.pos] == '\t') {
			vec.pos++
		}
		if vec.pos >= uint64(len(s)) {
			return
		}
		switch s[vec.pos] {
		case '#':
			vec.skipln()
		case '%':
			vec.skipln()
		case '\r':
			vec.pos++
			if vec.pos < uint64(len(s)) && s[vec.pos] == '\n' {
				vec.pos++
			}
		case '\n':
			vec.pos++
		default:
			return
		}
	}
}

// lineStart returns position of the current line start.
func (vec *Vector) lineStart() int {
	s := vec.Src()
	i := int(vec.pos)
	for i > 0 && s[i-1] != '\n' {
		i--
	}
	return i
}

// lineIndent returns count of leading spaces of the current line.
func (vec *Vector) lineIndent() int {
	s := vec.Src()
	i, c := vec.lineStart(), 0
	for i < len(s) && s[i] == ' ' {
		c++
		i++
	}
	return c
}

// isNumeric reports whether b can be parsed as a number.
func isNumeric(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	_, err := strconv.ParseFloat(byteconv.B2S(b), 64)
	return err == nil
}

var (
	bdocStart = []byte("---")
	bdocEnd   = []byte("...")
)

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
	vec.offmove(uint64(w))
	r1, _, err1 := vec.ReadRuneAt(int(vec.pos))
	if err1 != nil && err1 != io.ErrUnexpectedEOF {
		return nil, err1
	}
	switch {
	case r == '-' && r1 == ' ':
		// multiline literal
		vec.t.typ = tokenDash
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.offmove(1)
		return &vec.t, nil
	case r == ':':
		vec.t.typ = tokenColon
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.offmove(1)
		return &vec.t, nil
	case r == ',':
		vec.t.typ = tokenComma
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.offmove(1)
		return &vec.t, nil
	case r == '[':
		vec.t.typ = tokenLBracket
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.offmove(1)
		return &vec.t, nil
	case r == ']':
		vec.t.typ = tokenRBracket
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.offmove(1)
		return &vec.t, nil
	case r == '{':
		vec.t.typ = tokenLBrace
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.offmove(1)
		return &vec.t, nil
	case r == '}':
		vec.t.typ = tokenRBrace
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.offmove(1)
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
		vec.offmove(hi - vec.pos)
		return &vec.t, nil
	case r == '&':
		vec.t.typ = tokenAnchor
		hi, err2 := vec.readAnchor()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.offmove(hi - vec.pos)
		return &vec.t, nil
	case r == '*':
		vec.t.typ = tokenAlias
		vec.t.setlo(vec.pos).sethi(vec.pos + 1)
		vec.offmove(1)
		return &vec.t, nil
	case r == '!':
		vec.t.typ = tokenTag
		hi, err2 := vec.readTag()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.offmove(hi - vec.pos)
		return &vec.t, nil
	case r == '%':
		vec.t.typ = tokenDirective
		hi, err2 := vec.readDirective()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.offmove(hi - vec.pos)
		return &vec.t, nil
	case r == '"' || r == '\'':
		vec.t.typ = tokenString
		lo := vec.pos
		hi, err2 := vec.readString(byte(r))
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(lo).sethi(hi)
		return &vec.t, nil
	case unicode.IsDigit(r) || r == '-' || r == '+' || (r == '.' && unicode.IsDigit(r1)):
		if r == '-' && r1 == '.' {
			// possible negative infinity
			r2, _, err2 := vec.ReadRuneAt(int(vec.pos + 1))
			if err2 == nil && (r2 == 'i' || r2 == 'I') {
				vec.pos--
				typ, hi, err2 := vec.readKeyword()
				if err2 != nil {
					return nil, err2
				}
				vec.t.setlo(vec.pos).sethi(hi)
				vec.offmove(hi - vec.pos)
				vec.t.typ = typ
				return &vec.t, nil
			}
		}
		vec.t.typ = tokenNumber
		hi, nan, err2 := vec.readNumber()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.offmove(hi - vec.pos)
		if nan {
			vec.t.typ = tokenString
		}
		return &vec.t, nil
	case r == '.' && (r1 == 'i' || r1 == 'I' || r1 == 'n' || r1 == 'N'):
		vec.pos--
		vec.t.typ = tokenNumber
		typ, hi, err2 := vec.readKeyword()
		if err2 != nil {
			return nil, err2
		}
		vec.t.setlo(vec.pos).sethi(hi)
		vec.offmove(hi - vec.pos)
		vec.t.typ = typ
		return &vec.t, nil
	case unicode.IsLetter(r) || r == '~':
		vec.pos--
		off := vec.pos
		typ, hi, err2 := vec.readKeyword()
		if err2 != nil {
			return nil, err2
		}
		vec.t.typ = typ
		vec.t.setlo(off).sethi(hi)
		vec.offmove(hi - vec.pos)
		return &vec.t, nil
	case r == '>' || r == '|':
		vec.t.typ = tokenString
		off := vec.pos
		i, j := skipline.Index2(vec.Src()[off:])
		if i == -1 {
			return nil, vector.ErrUnexpId
		}
		vec.pos = off + uint64(j)
		off = vec.pos
		eow, err := vec.skipws()
		if err != nil {
			return nil, err
		}
		pad := uint64(eow)
		for {
			i, j = skipline.Index2(vec.Src()[vec.pos:])
			if i == -1 {
				vec.pos = uint64(vec.SrcLen())
				break
			}
			if c := vec.SrcAt(int(vec.pos)); c == '\n' || c == 'r' {
				vec.pos++
				continue
			}
			vec.pos += uint64(j)
			if r == '>' {
				vec.Src()[vec.pos-1] = ' '
			}
			if eow, err = vec.skipws(); err != nil {
				return nil, err
			}
			pad1 := uint64(eow)
			if pad1 > pad {
				return nil, ErrBadIndent
			}
			if pad1 < pad || vec.pos == uint64(vec.SrcLen()) {
				break
			}
		}
		vec.t.setlo(off).sethi(vec.pos)
		return &vec.t, nil
	case r == '\r' || (r == '\n' && r1 == '\r'):
		vec.pos++
	default:
		return nil, vector.ErrUnexpId
	}
	return nil, vector.ErrUnexpId
}

func (vec *Vector) readComment() (uint64, error) {
	b := vec.Src()[vec.pos:]
	i, j := skipline.Index2(b)
	if i == -1 {
		j = len(b)
	}
	vec.pos = vec.pos + uint64(j)
	return vec.pos, nil
}

func (vec *Vector) readAnchor() (uint64, error) {
	for !vec.eof() {
		switch vec.Src()[vec.pos] {
		case ' ', '\t', '\n', '\r':
			return vec.pos, nil
		}
		vec.pos++
	}
	return vec.pos, nil
}

func (vec *Vector) readTag() (uint64, error) {
	for !vec.eof() {
		switch vec.Src()[vec.pos] {
		case ' ', '\t', '\n', '\r':
			return vec.pos, nil
		}
		vec.pos++
	}
	return vec.pos, nil
}

func (vec *Vector) readDirective() (uint64, error) {
	for !vec.eof() {
		switch vec.Src()[vec.pos] {
		case '\n', '\r':
			return vec.pos, nil
		}
		vec.pos++
	}
	return vec.pos, nil
}

func (vec *Vector) readString(b byte) (uint64, error) {
	p := vec.Src()
	i := bytealg.IndexByteAtBytes(p, b, int(vec.pos))
	if i < 0 {
		return 0, vector.ErrUnexpEOF
	}
	if p[i-1] != '\\' {
		vec.pos = uint64(i + 1)
		return uint64(i), nil
	} else {
		for j := i; j < len(p); {
			j = bytealg.IndexByteAtBytes(p, b, j+1)
			if i < 0 {
				return 0, vector.ErrUnexpEOF
			}
			i = j
			if p[j-1] != '\\' {
				break
			}
		}
	}
	vec.pos = uint64(i + 1)
	return vec.pos - 1, nil
}

func (vec *Vector) readNumber() (uint64, bool, error) {
	p := vec.Src()
	pl := uint64(len(p))
	var i uint64
	vec.pos--
	for i = vec.pos; i < pl; i++ {
		if !unicode.IsDigit(rune(p[i])) && p[i] != '.' && p[i] != 'e' && p[i] != 'E' && p[i] != '-' && p[i] != '+' {
			j := bytealg.IndexByteAtBytes(p, '\n', int(i))
			if j < 0 {
				j = int(pl)
			}
			i = uint64(j)
			return i, true, nil
		}
	}
	if i == pl {
		return pl, false, nil
	}
	return i, false, nil
}

func (vec *Vector) readKeyword() (ttoken, uint64, error) {
	off := vec.pos
	i, j := skipline.Index2(vec.Src()[off:])
	if i == -1 {
		i = vec.SrcLen()
		j = i
	}
	v := vec.Src()[off : off+uint64(i)]
	if k := bytealg.IndexByteAtBytes(v, ':', 0); k > 0 {
		vec.pos = uint64(k)
		return tokenString, vec.pos, nil
	}
	vec.pos = uint64(j)
	switch {
	case bytes.Equal(v, bnull) || bytes.Equal(v, bNull) || bytes.Equal(v, bNULL), bytes.Equal(v, bNone) || bytes.Equal(v, bTilda):
		return tokenNull, vec.pos, nil
	case bytes.Equal(v, btrue) || bytes.Equal(v, bTrue) || bytes.Equal(v, bTRUE) || bytes.Equal(v, bOn),
		bytes.Equal(v, bfalse) || bytes.Equal(v, bFalse) || bytes.Equal(v, bFALSE) || bytes.Equal(v, bOff):
		return tokenBool, vec.pos, nil
	case bytes.Equal(v, binf) || bytes.Equal(v, bInf) || bytes.Equal(v, bINF):
		return tokenInf, vec.pos, nil
	case bytes.Equal(v, bninf) || bytes.Equal(v, bnInf) || bytes.Equal(v, bnINF):
		return tokenNInf, vec.pos, nil
	case bytes.Equal(v, bnan) || bytes.Equal(v, bNaN) || bytes.Equal(v, bNAN):
		return tokenNaN, vec.pos, nil
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
	bon    = []byte("on")
	bON    = []byte("ON")
	boff   = []byte("off")
	bOFF   = []byte("OFF")
	binf   = []byte(".inf")
	bInf   = []byte(".Inf")
	bINF   = []byte(".INF")
	bninf  = []byte("-.inf")
	bnInf  = []byte("-.Inf")
	bnINF  = []byte("-.INF")
	bnan   = []byte(".nan")
	bNaN   = []byte(".NaN")
	bNAN   = []byte(".NAN")

	bnums  = []byte("Inf-InfNaN")
	pbnums uintptr
)

func init() {
	h := *(*byteconv.SliceHeader)(unsafe.Pointer(&bnums))
	pbnums = h.Data
}
