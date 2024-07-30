package yamlvector

import "unsafe"

// Table based approach of fmt skip.
func skipFmt(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	_ = skipTable[255]
	if n-offset > 512 {
		offset, _ = skipFmtBin8(src, n, offset)
	}
	for ; skipTable[src[offset]]; offset++ {
	}
	return offset, offset == n
}

// Binary based approach of fmt skip.
func skipFmtBin8(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	_ = skipTable[255]
	if *(*uint64)(unsafe.Pointer(&src[offset])) == binNlSpace7 {
		offset += 8
		for offset < n && *(*uint64)(unsafe.Pointer(&src[offset])) == binSpace8 {
			offset += 8
		}
	}
	return offset, false
}

var (
	skipTable   = [256]bool{}
	binNlSpace7 uint64
	binSpace8   uint64
)

func init() {
	skipTable[' '] = true
	skipTable['\t'] = true
	skipTable['\n'] = true
	skipTable['\t'] = true

	binNlSpace7Bytes, binSpace8Bytes := []byte("\n       "), []byte("        ")
	binNlSpace7, binSpace8 = *(*uint64)(unsafe.Pointer(&binNlSpace7Bytes[0])), *(*uint64)(unsafe.Pointer(&binSpace8Bytes[0]))
}
