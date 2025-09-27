package yamlvector

func (vec *Vector) skipws() (c int, err error) {
	for vec.pos < uint64(vec.SrcLen()) {
		switch vec.SrcAt(int(vec.pos)) {
		case ' ':
			vec.pos++
			c++
		case '\t':
			err = ErrBadIndent
			return
		default:
			return
		}
	}
	return
}
