package yamlvector

import "unsafe"

func bin(src []byte, offset, size int) uint64 {
	n := len(src)
	if offset+size >= n {
		return 0
	}
	_ = src[n-1]
	binSrc := src[offset : offset+size]
	b := *(*uint64)(unsafe.Pointer(&binSrc[0]))
	return b
}

var (
	bBoolTrue = [11][]byte{
		[]byte("true"),
		[]byte("True"),
		[]byte("TRUE"),
		[]byte("y"),
		[]byte("Y"),
		[]byte("yes"),
		[]byte("Yes"),
		[]byte("YES"),
		[]byte("on"),
		[]byte("On"),
		[]byte("ON"),
	}
	binBoolTrue = [11]uint64{}
	bBoolFalse  = [11][]byte{
		[]byte("false"),
		[]byte("False"),
		[]byte("FALSE"),
		[]byte("n"),
		[]byte("N"),
		[]byte("no"),
		[]byte("No"),
		[]byte("NO"),
		[]byte("off"),
		[]byte("Off"),
		[]byte("OFF"),
	}
	binBoolFalse = [11]uint64{}
)

func init() {
	for i := 0; i < len(bBoolTrue); i++ {
		binBoolTrue[i] = *(*uint64)(unsafe.Pointer(&bBoolTrue[i][0]))
	}
	for i := 0; i < len(bBoolFalse); i++ {
		binBoolFalse[i] = *(*uint64)(unsafe.Pointer(&bBoolFalse[i][0]))
	}
}
