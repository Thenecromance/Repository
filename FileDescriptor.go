package Repository

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

type fileDescriptor struct {
	Uid         string `json:"uid"` // the only way to identify a file
	ContentHash string `json:"ContentHash"`
	FileName    string `json:"fileName"`
}

// Hash return the hash of the file
func (fd *fileDescriptor) Hash() string {
	return fd.ContentHash
}

func (fd *fileDescriptor) GetDirectory(_len int) string {
	return fd.ContentHash[:_len]
}

func (fd *fileDescriptor) GetStoreName(_len int) string {
	return fd.ContentHash[_len:]
}

func (fd *fileDescriptor) GetFullPath(dir string, len_ int) string {
	var builder strings.Builder
	{
		builder.Grow(len(dir) + 2 + 64)
		builder.WriteString(dir)
		builder.WriteByte('/')
		builder.WriteString(fd.GetDirectory(len_))
		builder.WriteByte('/')
		builder.WriteString(fd.GetStoreName(len_))
	}
	return builder.String()
	//return dir + "/" + fd.GetDirectory(len_) + "/" + fd.GetStoreName(len_)
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
