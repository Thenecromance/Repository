package Buffer

import (
	"testing"
)

//func BenchmarkDoubleBuffer_New(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		NewDoubleBuffer(testSize)
//	}
//	b.StopTimer()
//
//}
//
//func BenchmarkDoubleBuffer_AppendSingle(b *testing.B) {
//	buf := NewDoubleBuffer(testSize)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		buf.Append(i)
//	}
//	b.StopTimer()
//}
//
//func BenchmarkDoubleBuffer_AppendGroup(b *testing.B) {
//	buf := NewDoubleBuffer(testSize)
//	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		buf.Append(a)
//	}
//	b.StopTimer()
//}

func BenchmarkDoubleBuffer_AppendGroupOverSize(b *testing.B) {
	buf := NewDoubleBuffer(testSize)
	a := make([]obj, testSize+testSize/2)
	for i := 0; i < testSize+testSize/2; i++ {
		a[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Append(a...)
		buf.Clear()
	}
	b.StopTimer()
}
