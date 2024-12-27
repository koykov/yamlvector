package yamlvector

import "github.com/koykov/bytealg"

func scanl(src []byte, n, offset int) (pos, semicolon int, eof bool) {
	semicolon = -1
	_ = src[n-1]
	for pos = offset; pos < n; pos++ {
		if src[pos] == ':' {
			semicolon = pos
		}
	}
	eof = pos == n-1
	return
}

func scans(src []byte, b byte, offset int) (pos int, eof bool) {
	for {
		if pos = bytealg.IndexByteAtBytes(src, b, offset); pos == -1 {
			return pos, true
		}
		if src[pos-1] != '\\' {
			break
		}
	}
	return
}
