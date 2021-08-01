package test_cases

//go:genopt
type originalObject struct {
	optionInt int
	configInt int `opt:required`
}
