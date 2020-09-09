package common_errors

import (
	"fmt"
)

//
// ConfigNotFound returns when the specified config File is missing.
//
type ConfigNotFound struct {
	File string
	Err  error
}

func (cnf ConfigNotFound) Error() string {
	return fmt.Sprintf(
		"config File does not exist - %s, Err: %s",
		cnf.File, cnf.Err)
}

//
// ConfigIOError represents error associated with IO interaction.
//
type ConfigIOError struct {
	File string
	Err  error
}

func (cie ConfigIOError) Error() string {
	return fmt.Sprintf(
		"failed with error (%s) to read config File from %s",
		cie.Err, cie.File)
}

//
// ConfigParsingError represents error associated with YAML parsing.
//
type ConfigParsingError struct {
	File string
	Err  error
}

func (cpe ConfigParsingError) Error() string {
	return fmt.Sprintf(
		"failed with error (%s) to parse YAML from File %s",
		cpe.Err, cpe.File)
}

//
// TemplateParsingError represents a template parsing error.
//
type TemplateParsingError struct {
	Err      error
	Template string
}

func (tpe TemplateParsingError) Error() string {
	return fmt.Sprintf(
		"failed to parse a template error - %s with template - %s",
		tpe.Template,
		tpe.Err)
}

//
// ParamNotFoundError represents an error when a desire key doesn't exist.
//
type KeyNotFoundError struct {
	Key string
}

func (knfe KeyNotFoundError) Error() string {
	return fmt.Sprintf("KeyNotFound:%s", knfe.Key)
}
