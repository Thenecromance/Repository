package Repository

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

const (
	SHA1 = iota
	SHA256
	SHA224
	SHA384
	SHA512
	SHA512_224
	SHA512_256
	MD5
)

func applyAlgorithm(algId int) hash.Hash {
	sha1.New()
	switch algId {
	case SHA1:
		return sha1.New()
	case SHA256:
		return sha256.New()
	case SHA224:
		return sha256.New224()
	case SHA384:
		return sha512.New384()
	case SHA512:
		return sha512.New()
	case SHA512_224:
		return sha512.New512_224()
	case SHA512_256:
		return sha512.New512_256()
	case MD5:
		return md5.New()
	default:
		return sha1.New()

	}
}

func allocateBuffer(algId int) []byte {
	switch algId {
	case SHA1:
		return make([]byte, sha1.Size)
	case SHA256:
		return make([]byte, sha256.Size)
	case SHA224:
		return make([]byte, sha256.Size224)
	case SHA384:
		return make([]byte, sha512.Size384)
	case SHA512:
		return make([]byte, sha512.Size)
	case SHA512_224:
		return make([]byte, sha512.Size224)
	case SHA512_256:
		return make([]byte, sha512.Size256)
	case MD5:
		return make([]byte, md5.Size)
	default:
		return make([]byte, sha1.Size)
	}
}
