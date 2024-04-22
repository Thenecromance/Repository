package lookupTable

import (
	"sync"
	"testing"
)

func BenchmarkSyncMap_Store(b *testing.B) {
	var m sync.Map

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
	b.StopTimer()
}

func BenchmarkSyncMap_Load(b *testing.B) {
	var m sync.Map
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Load(i)
	}
	b.StopTimer()
}
