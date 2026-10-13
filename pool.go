package yamlvector

import (
	"sync"

	"github.com/koykov/vector"
)

// Pool represents YAML vectors pool.
type Pool struct {
	p sync.Pool
}

var (
	// P is a default instance of the pool.
	// Just call yamlvector.Acquire() and yamlvector.Release().
	P Pool
	// Suppress go vet warnings.
	_, _ = Acquire, Release
)

// Get old vector from the pool or create new one.
func (p *Pool) Get() *Vector {
	v := p.p.Get()
	if v != nil {
		if vec, ok := v.(*Vector); ok {
			vec.SetBit(vector.FlagInit, true)
			vec.Helper = helper
			return vec
		}
	}
	return NewVector()
}

// Put vector back to the pool.
func (p *Pool) Put(vec *Vector) {
	vec.Reset()
	p.p.Put(vec)
}

// Acquire returns vector from default pool instance.
func Acquire() *Vector {
	return P.Get()
}

// Release puts vector back to default pool instance.
func Release(vec *Vector) {
	P.Put(vec)
}
