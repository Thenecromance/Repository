package mmap

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"syscall"
)

type MMap struct {
	lock sync.RWMutex // if support write, so need a lock to make sure the data's safety
	data []byte
}

func (m *MMap) Close() error {
	if m.data == nil {
		return nil
	} else if len(m.data) == 0 {
		m.data = nil
		return nil
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	data := m.data
	m.data = nil
	runtime.SetFinalizer(m, nil) // reset the finalizer
	return syscall.Munmap(data)
}

func (m *MMap) Len() int {
	return len(m.data)
}

func (m *MMap) At(i int) byte {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.data[i]
}

func (m *MMap) ReadAt(b []byte, off int64) (n int, err error) {
	if m.data == nil {
		return 0, errors.New("mmap: closed")
	}

	if off < 0 || int64(len(m.data)) < off {
		return 0, fmt.Errorf("mmap: invalid ReadAt offset %d", off)
	}

	m.lock.RLock()
	defer m.lock.RUnlock()

	n = copy(b, m.data[off:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

func (m *MMap) Read() (b []byte, err error) {
	if m.data == nil {
		return nil, errors.New("mmap: closed")
	}

	m.lock.RLock()
	defer m.lock.RUnlock()

	b = make([]byte, len(m.data))
	copy(b, m.data)
	return b, nil
}

func (m *MMap) WriteAt(b []byte, off int64) (n int, err error) {
	if m.data == nil {
		return 0, errors.New("mmap: closed")
	}

	if off < 0 || int64(len(m.data)) < off {
		return 0, fmt.Errorf("mmap: invalid WriteAt offset %d", off)
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	n = copy(m.data[off:], b)
	if n < len(b) {
		return n, io.ErrShortWrite
	}
	return n, nil
}

// for support io.Writer
func (m *MMap) Write(p []byte) (n int, err error) {
	if m.data == nil {
		return 0, errors.New("mmap: closed")
	}

	// mmap not support to extend the size of the file
	if len(p) > len(m.data) {
		return 0, fmt.Errorf("mmap: invalid Write size %d", len(p))
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	n = copy(m.data, p)
	if n < len(p) {
		return n, io.ErrShortWrite
	}

	return n, nil
}

// Open memory-maps the named file for reading.
func Open(filename string) (*MMap, error) {
	//f, err := os.Open(filename)
	f, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := fi.Size()
	if size == 0 {
		// Treat (size == 0) as a special case, avoiding the syscall, since
		// "man 2 mmap" says "the length... must be greater than 0".
		//
		// As we do not call syscall.Mmap, there is no need to call
		// runtime.SetFinalizer to enforce a balancing syscall.Munmap.
		return &MMap{
			data: make([]byte, 0),
		}, nil
	}
	if size < 0 {
		return nil, fmt.Errorf("mmap: file %q has negative size", filename)
	}
	if size != int64(int(size)) {
		return nil, fmt.Errorf("mmap: file %q is too large", filename)
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED) //read and write available
	if err != nil {
		log.Println("mmap: ", err)
		return nil, err
	}
	r := &MMap{
		data: data,
	}
	//if debug {
	//	var p *byte
	//	if len(data) != 0 {
	//		p = &data[0]
	//	}
	//	println("mmap", r, p)
	//}
	runtime.SetFinalizer(r, (*MMap).Close)
	return r, nil
}
