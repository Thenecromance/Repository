package Buffer

const testSize = 100000

//
//func swapByReset(src, dst []int) {
//	src = src[:0]
//	src = append(src, dst...)
//}
//func swapByCopy(src, dst []int) {
//	copy(dst, src)
//	src = src[:0]
//}
//
//// by using pointer to speed up the swap
//func swapByPtrSwap(src, dst *[]int) {
//	*src, *dst = *dst, *src
//	*src = (*src)[:0] //cleanup the src
//
//}
//
//func BenchmarkSwap_SwapReset(b *testing.B) {
//	src := make([]int, testSize)
//	dst := make([]int, testSize)
//	for i := 0; i < testSize; i++ {
//		src[i] = i
//	}
//	var trigger bool
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		if !trigger {
//			swapByReset(src, dst)
//		} else {
//			swapByReset(dst, src)
//		}
//		trigger = !trigger
//	}
//	b.StopTimer()
//}
//
//func BenchmarkSwap_Copy(b *testing.B) {
//	src := make([]int, testSize)
//	dst := make([]int, testSize)
//	for i := 0; i < testSize; i++ {
//		src[i] = i
//	}
//	var trigger bool
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		if !trigger {
//			swapByCopy(src, dst)
//		} else {
//			swapByCopy(dst, src)
//		}
//		trigger = !trigger
//
//	}
//	b.StopTimer()
//}
//
//func BenchmarkSwap_PtrSwap(b *testing.B) {
//	src := make([]int, testSize)
//	dst := make([]int, testSize)
//	for i := 0; i < testSize; i++ {
//		src[i] = i
//	}
//	trigger := false
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		if !trigger {
//			swapByPtrSwap(&src, &dst)
//		} else {
//			swapByPtrSwap(&dst, &src)
//		}
//		trigger = !trigger
//
//	}
//	b.StopTimer()
//}
//
//func Benchmark_DataByAppend(b *testing.B) {
//	src := make([]int, testSize)
//	dst := make([]int, testSize)
//	for i := 0; i < testSize; i++ {
//		src[i] = i
//	}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		src = append(src, dst...)
//	}
//	b.StopTimer()
//}
//
//func Benchmark_DataByLoop(b *testing.B) {
//	src := make([]int, testSize)
//	dst := make([]int, testSize)
//	for i := 0; i < testSize; i++ {
//		src[i] = i
//	}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		for _, v := range dst {
//			src = append(src, v)
//		}
//	}
//	b.StopTimer()
//}
//
//// the fastest way to copy data
//func Benchmark_DataByCopyValue(b *testing.B) {
//	src := make([]int, testSize)
//	dst := make([]int, testSize)
//	for i := 0; i < testSize; i++ {
//		src[i] = i
//	}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		//src = append(src, dst...)
//		copy(dst, src)
//	}
//	b.StopTimer()
//}
//
//// then use this way to copy data
//func Benchmark_DataByCopyValueLoop(b *testing.B) {
//	src := make([]int, testSize)
//	dst := make([]int, testSize)
//	for i := 0; i < testSize; i++ {
//		src[i] = i
//	}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		//src = append(src, dst...)
//		for j, _ := range src {
//			copy(dst[j:], src[j:j+1])
//		}
//	}
//	b.StopTimer()
//
//}
