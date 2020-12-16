package configs

//
// Spec composes all relevant configuration together.
//
type Spec struct {
	Name    string
	Id      string
	Request request
}

//
// NewSpec returns an instance of integration struct based on the
// passed yamlConfig.
//
func NewSpec(yamlConfig []byte) (*Spec, error) {
	dcp := NewParser()

	return dcp.ReadBytes(yamlConfig)
}
