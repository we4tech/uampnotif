package notifcfg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/we4tech/uampnotif/pkg/common_errors"
)

//
// parser represents the instance of the config parser.
//
type parser struct {
}

//
// NewParser constructs and returns a instance of parser.
//
func NewParser() Parser {
	return &parser{}
}

//
// Read from the specified appConfigYaml and converts into a Config.
//
func (dcp *parser) Read(appConfigYaml string) (*Config, error) {
	_, err := os.Stat(appConfigYaml)

	if os.IsNotExist(err) {
		return nil, common_errors.ConfigNotFound{File: appConfigYaml, Err: err}
	}

	fileData, err := ioutil.ReadFile(appConfigYaml)

	if err != nil {
		return nil, common_errors.ConfigIOError{File: appConfigYaml, Err: err}
	}

	return dcp.readInternal(fileData, appConfigYaml)
}

func (dcp *parser) readInternal(
	fileData []byte,
	appConfigYaml string) (*Config, error) {

	var cfg = &Config{}
	if err := yaml.Unmarshal(fileData, cfg); err != nil {
		return nil, common_errors.ConfigParsingError{File: appConfigYaml, Err: err}
	}

	return cfg, nil
}
