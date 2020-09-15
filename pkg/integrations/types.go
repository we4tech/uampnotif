package integrations

import "github.com/we4tech/uampnotif/pkg/templates"

//
// ParsedRequest represents how to construct a request to the integrated system.
//
type ParsedRequest interface {
	Url(ctx *templates.TemplateContext) (string, error)
	Body(ctx *templates.TemplateContext) (string, error)
}

//
// SearchableParams provides an interface to search over a data structure.
//
type Searchable interface {
	//
	// Find a value by the specified needle
	//
	Find(needle string) string

	//
	// Exists lookup for a specific item, returns true if exists.
	//
	Exists(needle string) bool
}

//
// IterableParams provides an interface to iterate over a data structure.
//
type IterableParams interface {
	//
	// ForEach receives a Param whenever underlying implementation iterates.
	//
	ForEach(cb func(i int, p *Param))

	//
	// IsEmpty returns true if the underlying data structure doesn't have any
	// value.
	//
	IsEmpty() bool
}

//
// IterableHeaders provides an interface for iterating over the underlying
// data structure.
//
type IterableHeaders interface {
	//
	// ForEach receives a ParsedHeader whenever underlying implementation iterates.
	//
	ForEach(cb func(i int, h ParsedHeader) (bool, error)) error

	//
	// IsEmpty returns true if the underlying data structure doesn't have any
	// value.
	//
	IsEmpty() bool
}

//
// ParsedHeader represents a specific header with name and value template.
//
type ParsedHeader interface {
	//
	// GetValue returns a string after applying go-template parser. Error is
	// raised whenever template parser ends with an error.
	//
	GetValue(ctx *templates.TemplateContext) (string, error)

	GetName() string
}

//
// ConfigParser parses YAML configuration file and converts to
// Spec struct.
//
// Following errors are raised in case of the following scenarios.
//   - "Invalid yaml config" if failed to parse YAML file.
//   - "Invalid config but " if valid yaml but required attributes are missing.
//
type ConfigParser interface {
	//
	// Read a configuration YAML and converts into an Spec object.
	//
	Read(configYamlFile string) (Spec, error)
}

//
// Service represents the actual implementation based on Spec
// object.
//
type Service interface {
	//
	// Execute the underlying integration request with appropriate context.
	//
	Execute(i Spec, context templates.TemplateContext) (bool, error)
}
