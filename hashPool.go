package Repository

import (
	"encoding/hex"
	"hash"
	"sync"
)

type hashObject struct {
	hash.Hash
	buffer []byte
}

func newHashObject(algorithm int) *hashObject {
	return &hashObject{
		Hash:   applyAlgorithm(algorithm),
		buffer: allocateBuffer(algorithm),
	}
}

type Hash struct {
	pool sync.Pool //reduce the cost of creating a new hash object
}

func (h *Hash) Sum(content Content) string {
	obj := h.pool.Get().(*hashObject)
	defer obj.Reset()
	defer h.pool.Put(obj)

	obj.Write(content)
	obj.buffer = obj.Sum(nil)
	return hex.EncodeToString(obj.buffer)
}

func NewHash(algId int) *Hash {
	return &Hash{
		pool: sync.Pool{
			New: func() interface{} {
				return newHashObject(algId)
			},
		},
	}
}
