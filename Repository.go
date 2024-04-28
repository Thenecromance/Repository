package Repository

import (
	"Repository/shelf"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Content = *[]byte

type object struct {
	ptr     *fileDescriptor
	content Content
}

type Repository struct {
	resourceTable sync.Map
	dirLen        int
	rootDir       string

	hash    *Hash
	shelves map[string]*shelf.Shelf
	objPool sync.Pool

	preProcessChan chan *object
	releaseObjChan chan *object

	quit chan struct{}
}

func (r *Repository) ShelfCount() int {
	return 1<<(r.dirLen*4) - 1
}

// StoreFile is the method to store the file to the repository
func (r *Repository) StoreFile(fileName string, fileContent *[]byte) {

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

	str := hash.(string)

	//return r.shelves[f.GetId(r.dirLen)].GetItems(f.GetStoreName(r.dirLen))
	buffer := r.shelves[str[0:r.dirLen]].GetItems(str[r.dirLen:])
	return &buffer
}

// create the shelves
func (r *Repository) initShelves() {
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

func (r *Repository) storeInShelf(shelfId, Name string, content Content) {
	r.shelves[shelfId].NewItems(Name, content)
}
func (r *Repository) preProcess(obj *object) {
	defer func(ptr *object) {
		ptr.content = nil  // release the content
		ptr.ptr = nil      // release the pointer
		r.objPool.Put(ptr) //release the object
	}(obj)

	// calculate hash of the content
	r.hash.Sum(obj.content, &obj.ptr.ContentHash) // maybe 1 alloc?

	// store the files' info into the resource table,for delete/modify/search only
	r.resourceTable.Store(obj.ptr.FileName, obj.ptr.ContentHash)
	// then distribute the file to the shelf
	r.storeInShelf(obj.ptr.GetDirectory(r.dirLen), obj.ptr.GetStoreName(r.dirLen), obj.content)
}

func (r *Repository) run() {
	for {
		select {
		case obj := <-r.preProcessChan:
			{
				r.preProcess(obj)
			}
		case obj := <-r.releaseObjChan:
			{
				r.objPool.Put(obj)
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
		resourceTable:/*make(map[string]*fileDescriptor, 1024)*/ sync.Map{},
		dirLen:         2,
		rootDir:        "./resources",
		quit:           make(chan struct{}, 1),
		preProcessChan: make(chan *object, 1000),
		releaseObjChan: make(chan *object, 1000),
		objPool: sync.Pool{
			New: func() any {
				return &object{}
			},
		},
	}

	for _, opt := range opts {
		opt.apply(obj)
	}

	// if not set use which hash algorithm, use SHA1 as default
	if obj.hash == nil {
		obj.hash = NewHash(SHA1)
	}

	//create dir
	_ = os.Mkdir(obj.rootDir, os.ModePerm)

	obj.initShelves()

	for i := 0; i < 100; i++ {
		go obj.run()
	}

	return obj
}
