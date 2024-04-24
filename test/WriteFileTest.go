package main

import (
	"Repository"
	"Repository/test/random"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

type test struct {
	name    string
	content []byte
}

func storeFileSync(raw []test) {

	hash := func(obj *test) string {
		hash := sha1.New()
		hash.Write(obj.content)
		return hex.EncodeToString(hash.Sum(nil))
	}
	for _, obj := range raw {
		hashed := hash(&obj)
		//os.Create(fmt.Sprintf("./resources/%s/%s", hashed[:2], hashed[2:]))
		os.WriteFile(fmt.Sprintf("./resources/%s/%s", hashed[:2], hashed[2:]), obj.content, 0644)
	}
}

func main() {

	//int to string

	count := 10
	raw := make([]test, 0, count)

	r := Repository.New(

		Repository.WithAlgorithm(Repository.SHA256),
		Repository.WithDirLength(2),
	)
	//defer r.Close()

	{
		now := time.Now()
		for i := 0; i < count; i++ {
			raw = append(raw, test{name: random.RandomName(), content: random.Content()})
		}
		fmt.Printf("Time to generate %d random files:  %s\n", count, time.Since(now).String())
	}

	{
		now := time.Now()
		for i := 0; i < count; i++ {
			r.StoreFile(raw[i].name, raw[i].content)
		}

		fmt.Printf("Async sender write %d file cost : %s\n", count, time.Since(now).String())
	}

	//{
	//	now := time.Now()
	//	for i := 0; i < count; i++ {
	//		go r.StoreFile(raw[i].name, raw[i].content)
	//	}
	//	fmt.Printf(" sender running in 1 thread write %d file cost : %s\n", count, time.Since(now).String())
	//	time.Sleep(10 * time.Second)
	//}

	//{
	//	now := time.Now()
	//	storeFileSync(raw)
	//	fmt.Printf("Sync write %d file cost : %s\n", count, time.Since(now).String())
	//}
	time.Sleep(10 * time.Second)
	r.Close()
}
