package implementations

//AIScript holds the script proprties
type AIScript struct {
	device      string
	tag         string
	source      string
	entryPoints []string
}

//Device return the device type
func (m *AIScript) Device() string {
	return m.device
}

//SetDevice set the device type
func (m *AIScript) SetDevice(device string) {
	m.device = device
}

//Tag return the tag
func (m *AIScript) Tag() string {
	return m.tag
}

//SetTag set the tag
func (m *AIScript) SetTag(tag string) {
	m.tag = tag
}

//Source return the source
func (m *AIScript) Source() string {
	return m.source
}

//SetSource set the source
func (m *AIScript) SetSource(source string) {
	m.source = source
}

//EntryPoints return the entry points
func (m *AIScript) EntryPoints() []string {
	return m.entryPoints
}

//SetEntryPoints set the entry points
func (m *AIScript) SetEntryPoints(entryPoinys []string) {
	m.entryPoints = entryPoinys
}

//NewScript create a new AIScript with the given device
func NewScript(device string) *AIScript {
	return &AIScript{device: device}
}

//NewEmptycript create empty AIScript
func NewEmptycript() *AIScript {
	return &AIScript{}
}
