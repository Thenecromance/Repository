package Repository

import (
	"Repository/shelf"
	"encoding/json"
	"fmt"

	"os"
	"strconv"
	"sync"
)

type Content = []byte

type object struct {
	ptr     *fileDescriptor
	content Content
}

type Repository struct {
	ResourceTable sync.Map

	DirLength   int
	ResourceDir string
	hash        *Hash
	shelves     []*shelf.Shelf

	objPool        sync.Pool
	preProcessChan chan *object

	quit chan struct{}
}

func (r *Repository) ShelfCount() int {
	return 1<<(r.DirLength*4) - 1
}

func (r *Repository) StoreFile(fileName string, fileContent []byte) (token string) {

	info := newDescriptor(fileName) //1 alloc

	obj := r.objPool.Get().(*object) // 1 alloc
	// /*&object{}*/
	obj.ptr = info
	obj.content = fileContent

	//
	go func() {
		select {
		case r.preProcessChan <- obj:
		}
	}()
	return info.Uid
}

// GetFile is the method to get the file from the repository
func (r *Repository) GetFile(fileName string) Content {
	hash, ok := r.ResourceTable.Load(fileName)
	if !ok {
		return nil
	}
	f := newDescriptor(fileName)
	f.ContentHash = hash.(string)
	return r.shelves[f.GetId(r.DirLength)].GetItems(f.GetStoreName(r.DirLength))
}

func (r *Repository) loadShelves() {
	if r.DirLength < 1 {
		panic("DirLength must be greater than 0")
	}
	r.shelves = make([]*shelf.Shelf, 0, r.ShelfCount())
	for id := 0; id <= r.ShelfCount(); id++ {
		// format the id to hex string just like 000 to fff (based on the DirLength)
		strId := fmt.Sprintf("%0"+strconv.Itoa(r.DirLength)+"x", id)
		r.shelves = append(r.shelves,
			shelf.New(r.ResourceDir, strId),
		)
	}
}

func (r *Repository) preProcess(obj *object) {
	defer r.objPool.Put(obj)                      //release the object
	obj.ptr.ContentHash = r.hash.Sum(obj.content) // maybe 1 alloc?

	//r.ResourceTable.Store(obj.ptr.Uid, obj.ptr)
	r.ResourceTable.Store(obj.ptr.FileName, obj.ptr.ContentHash)
	// then distribute the file to the shelf
	r.shelves[obj.ptr.GetId(r.DirLength)].NewItems(obj.ptr.GetStoreName(r.DirLength), obj.content)
}

func (r *Repository) run() {
	for {
		select {
		case obj := <-r.preProcessChan:
			{
				r.preProcess(obj)
			}
		case <-r.quit:
			return
		}
	}
}

func (r *Repository) Close() {
	r.quit <- struct{}{}

	for _, s := range r.shelves {
		s.Close()
	}

	bytes, err := json.Marshal(r.ResourceTable)
	if err != nil {
		return
	}
	os.WriteFile("./table.json", bytes, os.ModePerm)
}

func New(opts ...Option) *Repository {
	obj := &Repository{
		ResourceTable:/*make(map[string]*fileDescriptor, 1024)*/ sync.Map{},
		DirLength:   2,
		ResourceDir: "./resources",
		quit:        make(chan struct{}, 1),
		objPool: sync.Pool{
			New: func() interface{} {
				return &object{}
			},
		},
		preProcessChan: make(chan *object, 1000),
	}

	for _, opt := range opts {
		opt.apply(obj)
	}

	// if not set use which hash algorithm, use SHA1 as default
	if obj.hash == nil {
		obj.hash = NewHash(SHA1)
	}

	//create dir
	os.Mkdir(obj.ResourceDir, os.ModePerm)

	obj.loadShelves()

	for i := 0; i < 100; i++ {
		go obj.run()
	}

	return obj
}
