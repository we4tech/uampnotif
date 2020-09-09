package notifiers

//
// Notifiers keeps a default settings and a list of notifiers.
//
type Notifiers struct {
	DefaultSettings Setting `yaml:"default_settings"`
	Notifiers       []Notifier
}

//
// Setting represents a default or a scoped setting for a or a group of
// integrations.
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
	Params   *Params
	Settings *Setting
}

//
// AppConfigParser interface to provide a way to interact with different
// implementations.
//
type AppConfigParser interface {
	//
	// Read a configuration yaml and convert into a Notifiers object.
	//
	Read(configYamlFile string) (*Notifiers, error)
}
