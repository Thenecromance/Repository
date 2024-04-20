package Repository

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
)

type fileDescriptor struct {
	Uid         string `json:"uid"` // the only way to identify a file
	ContentHash string `json:"ContentHash"`
	FileName    string `json:"fileName"`
}

func (fd *fileDescriptor) Hash() string {
	return fd.ContentHash
}

func (fd *fileDescriptor) GetDirectory(len int) string {
	return fd.ContentHash[:len]
}

func (fd *fileDescriptor) GetStoreName(len int) string {
	return fd.ContentHash[len:]
}

func (fd *fileDescriptor) GetFullPath(dir string, len int) string {
	return dir + "/" + fd.GetDirectory(len) + "/" + fd.GetStoreName(len)
}

func (fd *fileDescriptor) GetId(len int) uint32 {
	id, err := strconv.ParseUint(fd.GetDirectory(len), 16, len*4)
	if err != nil {
		panic(err)
	}
	return uint32(id)
}

func newDescriptor(fileName string) *fileDescriptor {
	return &fileDescriptor{
		Uid:      uuid.NewV4().String(),
		FileName: fileName,
	}

}
