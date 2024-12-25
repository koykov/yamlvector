package yamlvector

import "errors"

var (
	ErrBadUTF8   = errors.New("bad UTF8 sequence")
	ErrBadIndent = errors.New("bad indentation: only spaces allowed")
)
