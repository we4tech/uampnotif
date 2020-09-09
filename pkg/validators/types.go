package validators

//
// ## WHY?:
//
// The integration-validator validates a parsed struct to ensure all required
// attributes exist.
//

//
// Validator provides Validate and GetErrors methods to initiate a validation
// process.
//
type Validator interface {
	//
	// GetErrors returns a map of errors with corresponding field.
	//
	GetErrors() ValidationErrors

	//
	// Validate returns a boolean state to indicate valid or non-valid state.
	//
	Validate() bool
}

//
// ValidationErrors provides an interface to interact with validation specific
// errors.
//
type ValidationErrors interface {
	//
	// HasError returns true if an error exist for the specific field.
	//
	HasError(field string) bool

	//
	// GetError returns the error string for the specific field.
	//
	GetError(field string) string
}
