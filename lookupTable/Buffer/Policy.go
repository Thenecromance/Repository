package Buffer

const (
	// control about when the buffer is full, how to handle the new data
	block  = iota // block the new data until the buffer is not full
	ignore        // just ignore the new data (the data will be lost)
	custom
)

type FullPolicy = int
type Policy struct {
	// this policy control the behavior when the buffer is full and how to handle the new data
	NewDataWhenFull FullPolicy
}
