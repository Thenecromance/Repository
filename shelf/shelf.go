package shelf

import (
	"bufio"
	"os"
	"sync"
)

// Shelf is the directory's abstraction
type Shelf struct {
	dir      string
	name     string
	fullPath string
	locks    []sync.RWMutex // locks group can change to use lru cache to manage the locks also it can limit the number of locks

	quit     chan struct{}
	newItems chan item
}

type item struct {
	file    string
	content *[]byte
}

func (s *Shelf) GetName() string {
	return s.name
}

func (s *Shelf) NewItems(filename string, content *[]byte) {
	s.newItems <- item{ // 1 alloc
		file:    filename,
		content: content,
	}
}

func (s *Shelf) GetItems(name string) []byte {

	bytes, err := os.ReadFile(s.fullPath + "/" + name)
	if err != nil {
		return nil
	}
	return bytes

	//file, err := os.Open(s.fullPath + "/" + name)
	//if err != nil {
	//	return nil
	//}
	//defer file.Close()
	//reader := bufio.NewReader(file)
	//content, err := reader.ReadBytes('\n')
	//if err != nil {
	//	return nil
	//}
	//return content

}

func (s *Shelf) UpdateItems(oldName, newName string, content []byte) {
	// todo: this methods need to communicate with other storage's methods

}

func (s *Shelf) DeleteItems(name string) {

}

func (s *Shelf) storeItem(item *item) {
	//log.Print(s.name + "/" + item.file)
	file, err := os.Create(s.fullPath + "/" + item.file)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bufio.NewWriterSize(file, len(*item.content))

	_, err = file.Write(*item.content)
	if err != nil {
		panic(err)
	}
}

// Close the shelf
func (s *Shelf) Close() {
	s.quit <- struct{}{}
}

func (s *Shelf) run() {

	for {
		select {
		case <-s.quit:
			{
				return
			}
		case item := <-s.newItems: // store the file
			s.storeItem(&item)

		}
	}
}

func (s *Shelf) createDirectory() {
	err := os.Mkdir(s.fullPath, os.ModePerm)
	if err != nil {
		return
	}
}

func New(root, id string) *Shelf {
	s := &Shelf{
		dir:      root,
		name:     id,
		fullPath: root + "/" + id,
		quit:     make(chan struct{}, 1),
		newItems: make(chan item),
	}
	s.createDirectory()
	go s.run()
	return s
}
