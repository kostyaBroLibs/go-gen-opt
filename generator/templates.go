package genarator

import (
	"text/template"
)

var (
	templatePackage = template.Must(template.New("").Parse(`
package {{ .PackageName }}
`))
	templateOptionFunc = template.Must(template.New("").Parse(`
type {{ .ObjectNameCC }}Option func(object *{{ .ObjectNameCC }})
`))
	templateWithFunc = template.Must(template.New("").Parse(`
func With{{ .OptionNameCC }}(
	{{ .OptionNameLCC }} {{ .OptionType }},
) {{ .ObjectNameCC }}Option {
	return func({{ .ObjectNameLCC }} *{{ .ObjectNameCC}}) {
		{{ .ObjectNameLCC }}.{{ .OptionNameOC }} = {{ .OptionNameLCC }}
	}
}
`))
	templateInit = template.Must(template.New("").Parse(`
func New{{ .ObjectNameCC }}(
	options ...{{ .ObjectNameCC }}Option,
) (*{{ .ObjectNameCC }}, error) {
	{{ .ObjectNameLCC }} := &{{ .ObjectNameCC }}{}

	for _, option := range options {
		option({{ .ObjectNameLCC }})
	}

	return {{ .ObjectNameLCC }}, nil
}
`))
)
