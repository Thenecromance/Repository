package random

type RandFile struct {
	Name    string
	Content []byte
}

func init() {

}

func Random(arr []RandFile, n int) {
	arr = make([]RandFile, n)
	for i := 0; i < n; i++ {
		go func(i int) {
			arr[i].Name = RandomName()
			arr[i].Content = Content()
		}(i)
	}
}
