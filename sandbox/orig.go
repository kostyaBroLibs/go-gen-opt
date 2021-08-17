package main

//go:generate go-gen-opt

//go:genopt
type OriginalObject struct {
	optionInt int `genopt:"required"`
}
