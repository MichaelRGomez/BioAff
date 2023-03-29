// BIOAFF/backend/internal/validator/validator.go
package validator

import (
	"net/url"
	"regexp"
)

var (
	//regex for valid email
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	//regex for valid phone number
	PhoneRX = regexp.MustCompile(`^\+?\(?[0-9]{3}\)?\s?-\s?[0-9]{3}\s?-\s?[0-9]{4}$`)
)

// creating a type, that'll wrap out validation errors map
type Validator struct {
	Errors map[string]string
}

// New() - creates a new instance of a Validator
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid() - will check the "Errors map" for any entires
// if any entries are found, there's an error
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// In() - checks if the supplied element can be found within a provided list of elements
func In(element string, list ...string) bool {
	for i := range list {
		if element == list[i] {
			return true
		}
	}
	return false
}

// Matches() - will return true if a string value, matches a specific regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// ValidWebsite() - checks if a string value is a valid web URL
func ValidWebsite(website string) bool {
	_, err := url.ParseRequestURI(website)
	return err == nil
}

// AddError() adds an error entry to the Errors map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check() - performs validation checks and classes the AddError method in turn if an error entry needs to be added
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Unique() - checks that there are no repeating values in the slice
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
