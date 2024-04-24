package Buffer

type obj = any

type IBuffer interface {
	Append(data ...obj)
	Get() []obj
	Empty() bool
	HasData() bool
	Full() bool
	Clear()
	Size() int
}
