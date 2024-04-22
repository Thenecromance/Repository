package Buffer

import (
	"sync"
	"sync/atomic"
)

const (
	current = 0
	backup  = 1
)

type DoubleBuffer struct {
	lock          sync.RWMutex
	bufferTrigger atomic.Bool
	current       *Buffer
	buf           [2]*Buffer // here could be merged into a single buffer, just make it simple
}

func (d *DoubleBuffer) Append(data ...obj) {
	d.lock.Lock()
	defer d.lock.Unlock()
	sync.Map{}

	// if current buffer is full, swap to use backup buffer
	if d.buf[current].Full() {
		d.swap()
	}

	if d.buf[current].restCount()+d.buf[backup].restCount() < len(data) {
		// if both buffer can't store the data, then just ignore the data
		return
	}

	if d.buf[current].restCount() < len(data) {
		//so if the current buffer only has 10 space , but the data has 20, then the first 10 will be filled in the current buffer
		//and the rest 10 will be filled in the backup buffer
		//then swap the buffer to use the backup buffer
		cnt := d.buf[current].restCount()
		d.buf[current].Append(data[:cnt]...) // first fill the current buffer, then swap it
		d.buf[backup].Append(data[cnt:]...)  // then fill the backup buffer
		d.buf[current].swap(d.buf[backup])
	} else {
		d.buf[current].Append(data...)
	}
}

func (d *DoubleBuffer) append(data obj) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.current.Full() {

	}

	if d.buf[current].restCount() == 0 {
		return
	}

	d.buf[current].Append(data)
}

// return all buffer's record data count
func (d *DoubleBuffer) CachedSize() int {
	return int(d.buf[current].size.Load() + d.buf[backup].size.Load())
}

func (d *DoubleBuffer) swap() {
	defer d.bufferTrigger.Store(!d.bufferTrigger.Load())
	if d.bufferTrigger.Load() {
		d.current = d.buf[current]
	} else {
		d.current = d.buf[backup]
	}
}

func (d *DoubleBuffer) Clear() {
	d.buf[current].Clear()
	d.buf[backup].Clear()
}

// Empty returns true if the buffer is empty
func (d *DoubleBuffer) Empty() bool {
	return d.CachedSize() == 0
}
func (d *DoubleBuffer) HasData() bool {
	return !d.Empty()
}
func (d *DoubleBuffer) Full() bool {
	return d.buf[current].restCount() == 0 && d.buf[backup].restCount() == 0
}
func (d *DoubleBuffer) Get() (data []obj) {
	if d.Empty() {
		return []obj{}
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	// pre allocate the memory
	data = make([]obj, d.CachedSize())
	length := d.buf[current].size.Load()

	copy(data, d.buf[current].Get())
	copy(data[length:], d.buf[backup].Get())
	return
}

func NewDoubleBuffer(bufSize int) *DoubleBuffer {

	db := &DoubleBuffer{
		buf: [2]*Buffer{
			NewBuffer(bufSize),
			NewBuffer(bufSize),
		},
	}
	db.buf[current].disableLock()
	db.buf[backup].disableLock()
	return db
}
