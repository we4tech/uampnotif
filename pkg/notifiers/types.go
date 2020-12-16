package notifiers

//
// Config keeps a default settings and a list of notifiers.
//
type Config struct {
	DefaultSettings Setting `yaml:"default_settings"`
	Notifiers       []Notifier
}

//
// Setting represents a default or a scoped setting for a or a group of
// configs.
//
type Setting struct {
	Retries          int
	Async            bool
	OnError          string     `yaml:"on_error"`
	OnErrorNotifiers []Notifier `yaml:"on_error_notifiers"`
}

//
// Notifier stores a specific notifier integration specific configuration.
//
type Notifier struct {
	Id       string
	Params   Params
	Settings *Setting

	// Optional:
	Desc string `yaml:"desc,omitempty"`
}

//
// Parser interface to provide a way to interact with different
// implementations.
//
type Parser interface {
	//
	// Read a configuration yaml and convert into a Config object.
	//
	Read(configYamlFile string) (*Config, error)
}
