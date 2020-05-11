package implementations

import "reflect"

// TensorInterface is an interface that represents the skeleton of a tensor ( n-dimensional array of numerical data )
// needed to map it to a RedisAI Model with the proper operations
type AITensor struct {
	// the size - in each dimension - of the tensor.
	shape []int64

	data interface{}
}

func (t *AITensor) Dtype() reflect.Type {
	return reflect.TypeOf(t.data)
}

func NewAiTensor() *AITensor {
	return &AITensor{}
}

func (t *AITensor) NumDims() int64 {
	return int64(len(t.Shape()))
}

func (t *AITensor) Len() int64 {
	var result  int64 = 0
	for _, v := range t.shape {
		result += v
	}
	return result
}

func (m *AITensor) Shape() []int64 {
	return m.shape
}

func (m *AITensor) SetShape(shape []int64) {
	m.shape = shape
}

func NewAiTensorWithShape(shape []int64) *AITensor {
	return &AITensor{shape: shape}
}

func NewAiTensorWithData(typestr string, shape []int64, data interface{}) *AITensor {
	tensor := NewAiTensorWithShape(shape)
	tensor.SetData(data)
	return tensor
}

func (m *AITensor) SetData(data interface{}) {
	m.data = data
}

func (m *AITensor) Data() interface{} {
	return m.data
}
