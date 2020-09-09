package notifiers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/WeConnect/hello-tools/uampnotif/pkg/common_errors"
)

//
// DefaultConfigParser represents the instance of the config parser.
//
type DefaultConfigParser struct {
}

//
// NewDefaultConfigParser constructs and returns a instance of DefaultConfigParser.
//
func NewDefaultConfigParser() DefaultConfigParser {
	return DefaultConfigParser{}
}

//
// Read from the specified appConfigYaml and converts into a Notifiers.
//
func (dcp DefaultConfigParser) Read(appConfigYaml string) (*Notifiers, error) {
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

func (dcp DefaultConfigParser) readInternal(
	fileData []byte,
	appConfigYaml string) (*Notifiers, error) {

	var notifiers = &Notifiers{}
	if err := yaml.Unmarshal(fileData, notifiers); err != nil {
		return nil, common_errors.ConfigParsingError{File: appConfigYaml, Err: err}
	}

	return notifiers, nil
}
