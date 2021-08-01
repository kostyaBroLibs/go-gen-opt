package test_cases

import (
	"errors"
	"fmt"
)

var (
	ErrConfigEmpty = func(configName string) error {
		return fmt.Errorf("%w: %s must not been empty", ErrWrongConfig, configName)
	}
	ErrWrongConfig = errors.New("wrong config")
)

type Config struct {
	_         struct{}
	OptionInt int
}

type OriginalObjectOption func(object *originalObject)

func WithOptionInt(optionInt int) OriginalObjectOption {
	return func(o *originalObject) {
		o.optionInt = optionInt
	}
}

func NewOriginalObject(
	config Config,
	options ...OriginalObjectOption,
) (*originalObject, error) {
	oo := &originalObject{
		optionInt: config.OptionInt,
	}

	for _, option := range options {
		option(oo)
	}

	if oo.configInt == 0 {
		return nil, ErrConfigEmpty("configInt")
	}

	return oo, nil
}
