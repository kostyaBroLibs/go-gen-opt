package testdata

import (
	"context"
	"fmt"
)

type Config struct {
	_ struct{}

	Required string
	Host     string
	Port     int
}

type Option func(*Object)

func WithOptional(optional int) Option {
	return func(object *Object) {
		object.optional = optional
	}
}

func NewObject(
	ctx context.Context,
	config *Config,
	options ...Option,
) (*Object, error) {
	conn, err := createConnection(config.Host, config.Port)
	if err != nil {
		return nil, fmt.Errorf("createConnection error: %w", err)
	}

	object := &Object{
		ctx:      ctx,
		required: config.Required,
		conn:     conn,
		optional: 5,
	}

	for _, option := range options {
		option(object)
	}

	return object, nil
}
