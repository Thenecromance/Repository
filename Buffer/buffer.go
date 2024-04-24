package Buffer

import (
	"sync"
	"sync/atomic"
)

type Buffer struct {
	noLock bool
	lock   sync.RWMutex
	limit  int
	size   atomic.Int64
	buf    []obj
}

// when the buffer is used in a single goroutine or it is protected by other lock, you can disable the lock
func (b *Buffer) disableLock() {
	b.noLock = true
}

func (b *Buffer) Append(data ...obj) {
	if !b.noLock {
		b.lock.Lock() // write lock
		defer b.lock.Unlock()
	}

	// if the buffer is full, just ignore the new data
	if b.Full() {
		return
	}

	if b.restCount() < len(data) { //means there is no enough space to store the data
		copy(b.buf[b.size.Load():], data[:b.restCount()]) // so just fill the buffer with the rest space
		b.markBufferIsFull()
	} else {
		copy(b.buf[b.size.Load():], data)
		b.size.Add(int64(len(data)))
	}

}

func (b *Buffer) Empty() bool {
	return b.size.Load() == 0
}
func (b *Buffer) HasData() bool {
	return !b.Empty()
}
func (b *Buffer) Full() bool {
	return b.size.Load() == int64(b.limit)
}

func (b *Buffer) markBufferIsFull() {
	b.size.Store(int64(b.limit))
}
func (b *Buffer) Clear() {
	//b.buf = b.buf[:0] // if directly set b.size to 0, it means the buffer is empty, but the data is still in the buffer
	b.size.Store(0)
}

func (b *Buffer) swap(other *Buffer) {
	*b, *other = *other, *b //now other is need to be cleared
}

func (b *Buffer) restCount() int {
	return b.limit - int(b.size.Load())
}

// Get returns all the data in the buffer
func (b *Buffer) Get() []obj {
	if !b.noLock {
		b.lock.RLock() // read lock
		defer b.lock.RUnlock()
	}

	if b.size.Load() == 0 {
		return []obj{}
	}

	length := b.size.Load()
	b.size.Store(0)
	return b.buf[:length]
}

func (b *Buffer) Size() int {
	return b.limit
}

func NewBuffer(bufSize int) *Buffer {
	b := &Buffer{
		limit: bufSize,
		buf:   make([]obj, bufSize),
	}
	b.size.Store(0)
	return b
}
