package implementations

import "io/ioutil"

type AiModel struct {
	backend string
	device  string
	blob    []byte
	inputs  []string
	outputs []string
}

func (m *AiModel) Outputs() []string {
	return m.outputs
}

func (m *AiModel) SetOutputs(outputs []string) {
	m.outputs = outputs
}

func (m *AiModel) Inputs() []string {
	return m.inputs
}

func (m *AiModel) SetInputs(inputs []string) {
	m.inputs = inputs
}

func (m *AiModel) Blob() []byte {
	return m.blob
}

func (m *AiModel) SetBlob(blob []byte) {
	m.blob = blob
}

func (m *AiModel) Device() string {
	return m.device
}

func (m *AiModel) SetDevice(device string) {
	m.device = device
}

func (m *AiModel) Backend() string {
	return m.backend
}

func (m *AiModel) SetBackend(backend string) {
	m.backend = backend
}

func NewModel(backend string, device string) *AiModel {
	return &AiModel{backend: backend, device: device}
}

func NewEmptyModel() *AiModel {
	return &AiModel{}
}


func (m *AiModel) SetBlobFromFile(path string) (err error) {
	var data []byte
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	m.SetBlob(data)
	return
}
