package Repository

import (
	"Repository/test/random"
	"log"
	"math/rand"

	"testing"
)

const (
	count = 10000
)

type fileObj struct {
	name    string
	content []byte
}

var (
	raw  []fileObj
	repo *Repository
)

func init() {
	raw = make([]fileObj, 0, count)
	for i := 0; i < count; i++ {
		raw = append(raw, fileObj{name: random.RandomName(), content: random.Content()})
	}
	log.Println("Finish to generate random files")
	repo = New(WithAlgorithm(SHA1), WithDirLength(2), WithResourceDir("../../resources"))

}

func BenchmarkRepository_StoreFile(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.StoreFile(raw[i%count].name, raw[i%count].content)
	}
	b.StopTimer()

}

func BenchmarkRepository_StoreAsync(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			repo.StoreFile(raw[rand.Intn(count)].name, raw[rand.Intn(count)].content)
		}
	})
	b.StopTimer()
}
