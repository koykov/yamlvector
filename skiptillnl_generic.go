//go:build (!amd64 && !arm64 && !ppc64le && !riscv64) || appengine || !gc || purego

package yamlvector

func skipTillNl(b []byte) int { return skipTillNlGeneric(b) }
