package yamlvector

var funcAMD64 func([]byte) int

func init() {
	funcAMD64 = skipTillNlSSE2
}

func skipTillNl(b []byte) int {
	return funcAMD64(b)
}

//go:noescape
func skipTillNlSSE2(b []byte) int

//go:noescape
func skipTillNlAVX2(b []byte) int

//go:noescape
func skipTillNlAVX512(b []byte) int
