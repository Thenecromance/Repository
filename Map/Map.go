package Map

import (
	"sync"
)

// Todo : there will be a more faster way to store the data info instead of using the sync.Map

type Map struct {
	mtx sync.Mutex
	//atomic.Pointer[]
	dirty any
}

func (m *Map) Load(key any) (value any, exist bool) {

}

func (m *Map) Store(key, value any) {

}
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool) {

	return
}
func (m *Map) LoadAndDelete(key any) (value any, loaded bool) {
	return
}
func (m *Map) Delete(any) {

}
func (m *Map) Swap(key, value any) (previous any, loaded bool) {

}
func (m *Map) CompareAndSwap(key, old, new any) (swapped bool) {
	return
}

func (m *Map) CompareAndDelete(key, old any) (deleted bool) {
	return
}
func (m *Map) Range(func(key, value any) (shouldContinue bool)) {

}
