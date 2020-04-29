package implementations

// TensorInterface is an interface that represents the skeleton of a tensor ( n-dimensional array of numerical data )
// needed to map it to a RedisAI Model with the proper operations
type AiTensor struct {
	typestr string

	// the size - in each dimension - of the tensor.
	shape []int

	data interface{}
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

func (t *AiTensor) TypeString() string {
	return t.typestr
}

func (t *AiTensor) SetTypeString(typestr string) {
	t.typestr = typestr
}

func (m *AiTensor) Shape() []int {
	return m.shape
}

func (m *AiTensor) SetShape(shape []int) {
	m.shape = shape
}

func NewAiTensorWithTypeShape(typestr string, shape []int) *AiTensor {
	return &AiTensor{typestr: typestr, shape: shape}
}

func NewAiTensorWithData(typestr string, shape []int, data interface{}) *AiTensor {
	tensor := NewAiTensorWithTypeShape(typestr, shape)
	tensor.SetData(data)
	return tensor
}

func (m *AiTensor) SetData(data interface{}) {
	m.data = data
}

func (m *AiTensor) Data() interface{} {
	return m.data
}
