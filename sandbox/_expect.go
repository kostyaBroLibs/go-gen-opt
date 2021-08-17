package main

type OriginalObjectOption func(object *OriginalObject)

func WithOptionInt(optionInt int) OriginalObjectOption {
	return func(o *OriginalObject) {
		o.optionInt = optionInt
	}
}

func NewOriginalObject(
	options ...OriginalObjectOption,
) (*OriginalObject, error) {
	originalObject := &OriginalObject{}

	for _, option := range options {
		option(originalObject)
	}

	return originalObject, nil
}

