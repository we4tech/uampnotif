package receivers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/we4tech/uampnotif/pkg/common_errors"
)

type parser struct{}

//
// NewParser constructs an instance.
//
func NewParser() Parser {
	return &parser{}
}

//
// Read takes the configYamlFile and converts into Spec after
// successful parsing struct.
//
func (dcp *parser) Read(configYamlFile string) (*Spec, error) {
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

func (dcp *parser) ReadBytes(configYaml []byte) (*Spec, error) {
	return dcp.readInternal(configYaml, "config.yaml-string")
}

func (dcp *parser) readInternal(
	fileData []byte,
	configYamlFile string) (*Spec, error) {
	var integration = &Spec{}

	if err := yaml.Unmarshal(fileData, integration); err != nil {
		return nil, common_errors.ConfigParsingError{File: configYamlFile, Err: err}
	}

	return integration, nil
}
