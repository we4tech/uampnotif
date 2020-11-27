package notifications

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/we4tech/uampnotif/pkg/common_errors"
)

//
// configParser represents the instance of the config parser.
//
type configParser struct {
}

//
// NewConfigParser constructs and returns a instance of configParser.
//
func NewConfigParser() ConfigParser {
	return &configParser{}
}

//
// Read from the specified appConfigYaml and converts into a Config.
//
func (dcp *configParser) Read(appConfigYaml string) (*Config, error) {
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

func (dcp *configParser) readInternal(
	fileData []byte,
	appConfigYaml string) (*Config, error) {

	var notifiers = &Config{}
	if err := yaml.Unmarshal(fileData, notifiers); err != nil {
		return nil, common_errors.ConfigParsingError{File: appConfigYaml, Err: err}
	}

	return notifiers, nil
}
