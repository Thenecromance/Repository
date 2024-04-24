package shelf

import (
	"bufio"
	"os"
	"sync"
)

// Shelf is the directory's abstraction
type Shelf struct {
	name  string
	locks []sync.RWMutex // locks group can change to use lru cache to manage the locks also it can limit the number of locks

	quit        chan struct{}
	newItems    chan item
	deleteItems chan string // just passed the file name to delete
}

type item struct {
	file    string
	content []byte
}

func (s *Shelf) GetName() string {
	return s.name
}

func (s *Shelf) NewItems(filename string, content []byte) {
	s.newItems <- item{ // 1 alloc
		file:    filename,
		content: content,
	}
}

func (s *Shelf) GetItems(name string) []byte {
	file, err := os.Open(s.name + "/" + name)
	if err != nil {
		return nil
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	content, err := reader.ReadBytes('\n')
	if err != nil {
		return nil
	}
	return content
}

func (s *Shelf) UpdateItems(oldName, newName string, content []byte) {
	// todo: this methods need to communicate with other storage's methods

}

func (s *Shelf) DeleteItems(name string) {
	s.deleteItems <- name
}

func (s *Shelf) storeItem(item *item) {
	//log.Print(s.name + "/" + item.file)
	file, err := os.Create(s.name + "/" + item.file)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bufio.NewWriterSize(file, len(item.content))

	_, err = file.Write(item.content)
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
		case <-s.deleteItems: // delete the file
		}
	}
}

func (s *Shelf) createDirectory() {
	err := os.Mkdir(s.name, os.ModePerm)
	if err != nil {
		return
	}
}

func New(root, id string) *Shelf {
	s := &Shelf{
		name:        root + "/" + id,
		quit:        make(chan struct{}, 1),
		newItems:    make(chan item),
		deleteItems: make(chan string),
	}
	s.createDirectory()
	go s.run()
	return s
}
