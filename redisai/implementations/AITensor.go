package implementations

import "reflect"

// TensorInterface is an interface that represents the skeleton of a tensor ( n-dimensional array of numerical data )
// needed to map it to a RedisAI Model with the proper operations
type AiTensor struct {
	// the size - in each dimension - of the tensor.
	shape []int

	data interface{}
}

func (t *AiTensor) Dtype() reflect.Type {
	return reflect.TypeOf(t.data)
}

func NewAiTensor() *AiTensor {
	return &AiTensor{}
}

func (t *AiTensor) NumDims() int {
	return len(t.Shape())
}

func (t *AiTensor) Len() int {
	result := 0
	for _, v := range t.shape {
		result += v
	}
	return result
}

func (m *AiTensor) Shape() []int {
	return m.shape
}

func (m *AiTensor) SetShape(shape []int) {
	m.shape = shape
}

func NewAiTensorWithShape(shape []int) *AiTensor {
	return &AiTensor{shape: shape}
}

func NewAiTensorWithData(typestr string, shape []int, data interface{}) *AiTensor {
	tensor := NewAiTensorWithShape(shape)
	tensor.SetData(data)
	return tensor
}

func (m *AiTensor) SetData(data interface{}) {
	m.data = data
}

func (m *AiTensor) Data() interface{} {
	return m.data
}
