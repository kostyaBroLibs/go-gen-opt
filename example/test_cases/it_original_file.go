package testcases

//go:genopt
type OriginalObject struct {
	optionInt int
	configInt int `opt:"required"`
}
