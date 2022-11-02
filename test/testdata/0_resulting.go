// Code generated by go-gen-opt v0.0.0. DO NOT EDIT.

package testdata

import (
	"errors"
	"fmt"
)

var (
	ErrConfigEmpty = func(configName string) error {
		return fmt.Errorf(
			"%w: %s must not been empty",
			ErrWrongConfig, configName,
		)
	}
	ErrWrongConfig = errors.New("wrong config")
)

type Config struct {
	_             struct{}
	RequiredField int
}

type OriginalObjectOption func(object *OriginalObject)

func WithOptionalField(optionalField int) OriginalObjectOption {
	return func(o *OriginalObject) {
		o.optionalField = optionalField
	}
}

func WithOptionalFieldWithDefault(optionalFieldWithDefault int) OriginalObjectOption {
	return func(o *OriginalObject) {
		o.optionalFieldWithDefault = optionalFieldWithDefault
	}
}

func NewOriginalObject(
	config Config,
	options ...OriginalObjectOption,
) (*OriginalObject, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	oo := &OriginalObject{
		optionalFieldWithDefault: 42,
		requiredField:            config.RequiredField,
	}

	for _, option := range options {
		option(oo)
	}

	return oo, nil
}

func (c Config) Validate() error {
	if c.RequiredField == 0 {
		return ErrConfigEmpty("requiredField")
	}

	return nil
}