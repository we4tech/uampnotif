package integrations

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/WeConnect/hello-tools/uampnotif/pkg/common_errors"
)

type DefaultConfigParser struct{}

//
// NewDefaultConfigParser constructs an instance.
//
func NewDefaultConfigParser() DefaultConfigParser {
	return DefaultConfigParser{}
}

//
// Read takes the configYamlFile and converts into IntegrationSpec after
// successful parsing struct.
//
func (dcp DefaultConfigParser) Read(configYamlFile string) (*IntegrationSpec, error) {
	_, err := os.Stat(configYamlFile)

	if os.IsNotExist(err) {
		return nil, common_errors.ConfigNotFound{File: configYamlFile, Err: err}
	}

	fileData, err := ioutil.ReadFile(configYamlFile)
	if err != nil {
		return nil, common_errors.ConfigIOError{File: configYamlFile, Err: err}
	}

	return dcp.readInternal(fileData, configYamlFile)
}

func (dcp DefaultConfigParser) readInternal(
	fileData []byte,
	configYamlFile string) (*IntegrationSpec, error) {
	var integration = &IntegrationSpec{}

	if err := yaml.Unmarshal(fileData, integration); err != nil {
		return nil, common_errors.ConfigParsingError{File: configYamlFile, Err: err}
	}

	return integration, nil
}
