package yamlvector

func (vec *Vector) skipws() (c int, err error) {
	for {
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
}
