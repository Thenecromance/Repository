package Buffer

import "testing"

func BenchmarkBuffer_New(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewBuffer(testSize)
	}
	b.StopTimer()

}

func BenchmarkBuffer_AppendSingle(b *testing.B) {
	buf := NewBuffer(testSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Append(i)
	}
	b.StopTimer()
}

func BenchmarkBuffer_AppendGroup(b *testing.B) {
	buf := NewBuffer(testSize)
	a := []obj{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Append(a...)
	}
	b.StopTimer()
}
