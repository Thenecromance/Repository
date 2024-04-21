package Buffer

import (
	"testing"
)

func BenchmarkDoubleBuffer_New(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewDoubleBuffer(testSize)
	}
	b.StopTimer()

}

func BenchmarkDoubleBuffer_AppendSingle(b *testing.B) {
	buf := NewDoubleBuffer(testSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Append(i)
	}
	b.StopTimer()
}

func BenchmarkDoubleBuffer_AppendGroup(b *testing.B) {
	buf := NewDoubleBuffer(testSize)
	a := createSlice()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Append(a...)
	}
	b.StopTimer()
}

func BenchmarkDoubleBuffer_AppendGroup_OverSized(b *testing.B) {
	buf := NewDoubleBuffer(testSize)
	a := createOverSizedSlice()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Append(a...)
		buf.Clear()
	}
	b.StopTimer()
}
