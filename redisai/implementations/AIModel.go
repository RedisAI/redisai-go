package implementations

import "io/ioutil"

//AIModel holds the model proprties
type AIModel struct {
	backend      string
	device       string
	blob         []byte
	inputs       []string
	outputs      []string
	tag          string
	batchSize    int64
	minBatchSize int64
}

func (m *AIModel) Outputs() []string {
	return m.outputs
}

func (m *AIModel) SetOutputs(outputs []string) {
	m.outputs = outputs
}

func (m *AIModel) Inputs() []string {
	return m.inputs
}

func (m *AIModel) SetInputs(inputs []string) {
	m.inputs = inputs
}

func (m *AIModel) Blob() []byte {
	return m.blob
}

func (m *AIModel) SetBlob(blob []byte) {
	m.blob = blob
}

func (m *AIModel) Device() string {
	return m.device
}

func (m *AIModel) SetDevice(device string) {
	m.device = device
}

func (m *AIModel) Backend() string {
	return m.backend
}

func (m *AIModel) SetBackend(backend string) {
	m.backend = backend
}

func (m *AIModel) Tag() string {
	return m.tag
}

func (m *AIModel) SetTag(tag string) {
	m.tag = tag
}

func (m *AIModel) BatchSize() int64 {
	return m.batchSize
}

func (m *AIModel) SetBatchSize(batchSize int64) {
	m.batchSize = batchSize
}

func (m *AIModel) MinBatchSize() int64 {
	return m.minBatchSize
}

func (m *AIModel) SetMinBatchSize(minBatchSize int64) {
	m.minBatchSize = minBatchSize
}

func NewModel(backend string, device string) *AIModel {
	return &AIModel{backend: backend, device: device}
}

func NewEmptyModel() *AIModel {
	return &AIModel{}
}

func (m *AIModel) SetBlobFromFile(path string) (err error) {
	var data []byte
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	m.SetBlob(data)
	return
}
