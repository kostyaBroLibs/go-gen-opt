package testcases

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
	_         struct{}
	ConfigInt int
}

type OriginalObjectOption func(object *OriginalObject)

func WithOptionInt(optionInt int) OriginalObjectOption {
	return func(o *OriginalObject) {
		o.optionInt = optionInt
	}
}

func NewOriginalObject(
	config Config,
	options ...OriginalObjectOption,
) (*OriginalObject, error) {
	oo := &OriginalObject{
		configInt: config.ConfigInt,
		optionInt: 0,
	}

	for _, option := range options {
		option(oo)
	}

	if oo.configInt == 0 {
		return nil, ErrConfigEmpty("configInt")
	}

	return oo, nil
}
