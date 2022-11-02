package testdata

// OriginalObject is the object
// for which must be generated init function with config and options.
//
//go:genopt
type OriginalObject struct {
	optionalField            int
	optionalFieldWithDefault int `opt:"default=42"`
	requiredField            int `opt:"required"`
}
