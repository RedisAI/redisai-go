package implementations

type AIScript struct {
	device      string
	tag         string
	source      string
	entryPoints []string
}

func (m *AIScript) Device() string {
	return m.device
}

func (m *AIScript) SetDevice(device string) {
	m.device = device
}

func (m *AIScript) Tag() string {
	return m.tag
}

func (m *AIScript) SetTag(tag string) {
	m.tag = tag
}

func (m *AIScript) Source() string {
	return m.source
}

func (m *AIScript) SetSource(source string) {
	m.source = source
}

func (m *AIScript) EntryPoints() []string {
	return m.entryPoints
}

func (m *AIScript) SetEntryPoints(entryPoinys []string) {
	m.entryPoints = entryPoinys
}

func NewScript(device string) *AIScript {
	return &AIScript{device: device}
}

func NewEmptycript() *AIScript {
	return &AIScript{}
}
