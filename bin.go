package yamlvector

import "unsafe"

func binSafe(src []byte, offset, size int) uint64 {
	n := len(src)
	if offset+size > n {
		return 0
	}
	_ = src[n-1]
	switch size {
	case 1:
		return uint64(src[offset])
	case 2:
		return uint64(src[offset+0]) |
			uint64(src[offset+1])<<8
	case 3:
		return uint64(src[offset+0]) |
			uint64(src[offset+1])<<8 |
			uint64(src[offset+2])<<16
	case 4:
		return uint64(src[offset+0]) |
			uint64(src[offset+1])<<8 |
			uint64(src[offset+2])<<16 |
			uint64(src[offset+3])<<24
	case 5:
		return uint64(src[offset+0]) |
			uint64(src[offset+1])<<8 |
			uint64(src[offset+2])<<16 |
			uint64(src[offset+3])<<24 |
			uint64(src[offset+4])<<32
	case 6:
		return uint64(src[offset+0]) |
			uint64(src[offset+1])<<8 |
			uint64(src[offset+2])<<16 |
			uint64(src[offset+3])<<24 |
			uint64(src[offset+4])<<32 |
			uint64(src[offset+5])<<40
	case 7:
		return uint64(src[offset+0]) |
			uint64(src[offset+1])<<8 |
			uint64(src[offset+2])<<16 |
			uint64(src[offset+3])<<24 |
			uint64(src[offset+4])<<32 |
			uint64(src[offset+5])<<40 |
			uint64(src[offset+6])<<48
	case 8:
		return uint64(src[offset+0]) |
			uint64(src[offset+1])<<8 |
			uint64(src[offset+2])<<16 |
			uint64(src[offset+3])<<24 |
			uint64(src[offset+4])<<32 |
			uint64(src[offset+5])<<40 |
			uint64(src[offset+6])<<48 |
			uint64(src[offset+7])<<56
	default:
		return 0
	}
}

func binUnsafe(src []byte, offset, size int) uint64 {
	n := len(src)
	if offset+size >= n {
		return 0
	}
	_ = src[n-1]
	binSrc := src[offset : offset+size]
	b := *(*uint64)(unsafe.Pointer(&binSrc[0]))
	return b & binMasks[size-1]
}

var binMasks = [8]uint64{
	0x00000000000000FF,
	0x000000000000FFFF,
	0x0000000000FFFFFF,
	0x00000000FFFFFFFF,
	0x000000FFFFFFFFFF,
	0x0000FFFFFFFFFFFF,
	0x00FFFFFFFFFFFFFF,
	0xFFFFFFFFFFFFFFFF,
}
