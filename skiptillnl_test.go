package yamlvector

import (
	"strconv"
	"testing"
)

type stnlStage struct {
	data []byte
	pos  int
}

var stnlStages []stnlStage

func init() {
	for i := 1; i < 1e10; i *= 10 {
		data := make([]byte, i)
		data = append(data, '\n')
		stnlStages = append(stnlStages, stnlStage{data: data, pos: len(data) - 1})
	}
}

func TestSkipTillNl(t *testing.T) {
	for _, st := range stnlStages {
		pos := skipTillNlAVX512(st.data)
		if pos != st.pos {
			t.Errorf("got %d, want %d", pos, st.pos)
		}
	}
}

func BenchmarkSkipTillNl(b *testing.B) {
	for _, st := range stnlStages {
		b.Run(strconv.Itoa(len(st.data)), func(b *testing.B) {
			b.Run("generic", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(st.data)))
				for i := 0; i < b.N; i++ {
					skipTillNlGeneric(st.data)
				}
			})
			b.Run("sse2", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(st.data)))
				for i := 0; i < b.N; i++ {
					skipTillNlSSE2(st.data)
				}
			})
			b.Run("avx2", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(st.data)))
				for i := 0; i < b.N; i++ {
					skipTillNlSSE2(st.data)
				}
			})
			b.Run("avx512", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(st.data)))
				for i := 0; i < b.N; i++ {
					skipTillNlSSE2(st.data)
				}
			})
		})
	}
}
