package integrations

//
// IntegrationSpec composes all relevant configuration together.
//
type IntegrationSpec struct {
	Name    string
	Id      string
	Request request
}

//
// NewIntegrationSpec returns an instance of integration struct based on the
// passed yamlConfig.
//
func NewIntegrationSpec(yamlConfig []byte) (*IntegrationSpec, error) {
	dcp := NewDefaultConfigParser()

	return dcp.readInternal(yamlConfig, "internal-buffer")
}
