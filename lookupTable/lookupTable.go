package lookupTable

type object = any // this should be a pointer to the fileDescriptor or something else which can record the file information
type LookUpTable struct {
	table map[string]object
	cache []object // for cached all the objects when the table can not hold all the objects
}

func New() *LookUpTable {
	return &LookUpTable{}
}
