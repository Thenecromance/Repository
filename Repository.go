package Repository

import (
	"Repository/shelf"
	"fmt"
	"log"

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
	resourceTable sync.Map
	dirLen        int
	rootDir       string

	hash           *Hash
	shelves        map[string]*shelf.Shelf
	objPool        sync.Pool
	preProcessChan chan *object

	quit chan struct{}
}

func (r *Repository) ShelfCount() int {
	return 1<<(r.dirLen*4) - 1
}

func (r *Repository) StoreFile(fileName string, fileContent []byte) {
	obj := r.objPool.Get().(*object) // 1 alloc

	obj.ptr = newDescriptor(fileName) //1 alloc
	// /*&object{}*/
	obj.content = fileContent

	//
	go func() {
		select {
		case r.preProcessChan <- obj:
		}
	}()
}

// GetFile is the method to get the file from the repository
func (r *Repository) GetFile(fileName string) Content {
	hash, ok := r.resourceTable.Load(fileName)
	if !ok {
		return nil
	}
	f := newDescriptor(fileName)
	f.ContentHash = hash.(string)
	//return r.shelves[f.GetId(r.dirLen)].GetItems(f.GetStoreName(r.dirLen))
	return r.shelves[f.GetDirectory(r.dirLen)].GetItems(f.GetStoreName(r.dirLen))
}

// create the shelves
func (r *Repository) createShelves() {
	if r.dirLen < 1 {
		panic("dirLen must be greater than 0")
	}

	//pre-allocate the slice
	r.shelves = make(map[string]*shelf.Shelf, r.ShelfCount())
	for id := 0; id <= r.ShelfCount(); id++ {
		// format the id to hex string just like 000 to fff (based on the dirLen)
		//r.shelves = append(r.shelves,
		//	//"%02x" means 2 characters, 0 padding, base 16
		//	shelf.New(r.rootDir, fmt.Sprintf("%0"+strconv.Itoa(r.dirLen)+"x", id)),
		//)
		shf := shelf.New(r.rootDir, fmt.Sprintf("%0"+strconv.Itoa(r.dirLen)+"x", id))

		r.shelves[shf.GetName()] = shf

	}
}

func (r *Repository) preProcess(obj *object) {
	defer r.objPool.Put(obj) //release the object

	// calculate hash of the content
	obj.ptr.ContentHash = r.hash.Sum(obj.content) // maybe 1 alloc?

	r.resourceTable.Store(obj.ptr.FileName, obj.ptr.ContentHash)
	// then distribute the file to the shelf
	r.shelves[obj.ptr.GetDirectory(r.dirLen)].NewItems(obj.ptr.GetStoreName(r.dirLen), obj.content)
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

}

func New(opts ...Option) *Repository {
	obj := &Repository{
		resourceTable: /*make(map[string]*fileDescriptor, 1024)*/ sync.Map{},
		dirLen:                                                   2,
		rootDir:                                                  "./resources",
		quit:                                                     make(chan struct{}, 1),
		objPool: sync.Pool{
			New: func() any {
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
	err := os.Mkdir(obj.rootDir, os.ModePerm)
	if err != nil {
		log.Println(err)
		return nil
	}

	obj.createShelves()

	for i := 0; i < 100; i++ {
		go obj.run()
	}

	return obj
}
