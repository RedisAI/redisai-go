package implementations

import "reflect"

// TensorInterface is an interface that represents the skeleton of a tensor ( n-dimensional array of numerical data )
// needed to map it to a RedisAI Model with the proper operations
type AITensor struct {
	// the size - in each dimension - of the tensor.
	shape []int

	data interface{}
}

func (t *AITensor) Dtype() reflect.Type {
	return reflect.TypeOf(t.data)
}

func NewAiTensor() *AITensor {
	return &AITensor{}
}

func (t *AITensor) NumDims() int {
	return len(t.Shape())
}

func (t *AITensor) Len() int {
	result := 0
	for _, v := range t.shape {
		result += v
	}
	return result
}

func (m *AITensor) Shape() []int {
	return m.shape
}

func (m *AITensor) SetShape(shape []int) {
	m.shape = shape
}

func NewAiTensorWithShape(shape []int) *AITensor {
	return &AITensor{shape: shape}
}

func NewAiTensorWithData(typestr string, shape []int, data interface{}) *AITensor {
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
