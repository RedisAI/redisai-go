package redisai


// Interface represents an n-dimensional array of numerical data.
type Interface interface {
	// Len returns the number of elements in the tensor.
	Len() int

	// Shape returns the size - in each dimension - of the tensor.
	Shape() []int64

	// NumDims returns the number of dimensions of the tensor.
	NumDims() int

	DataType() string

	Data() []interface{}

	Blob() []byte

}
