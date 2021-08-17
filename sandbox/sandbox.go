package main

import (
	"bytes"
	"fmt"
	"github.com/stoewer/go-strcase"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/inspector"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

type astStructForGen struct {
	typeSpec    *ast.TypeSpec
	structType  *ast.StructType
	packageName string
}

func main() {

	path := os.Getenv("GOFILE")
	if path == "" {
		log.Fatalf("GOFILE env variable must be set")
	}

	objsForGen, err := pullObjectsForGenerating(path)
	if err != nil {
		log.Fatalf("pull objects for generating error: %s", err.Error())
	}

	if len(objsForGen) == 0 {
		log.Println("objects for which need to generate code was not found")

		return
	}

}

func pullObjectsForGenerating(fileFullPath string) ([]*astStructForGen, error) {
	astInFile, err := parser.ParseFile(
		token.NewFileSet(),
		fileFullPath,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		return nil, fmt.Errorf("can not to parse file, error: %s", err.Error())
	}

	i := inspector.New([]*ast.File{astInFile})
	iFilter := []ast.Node{
		&ast.GenDecl{},
	}

	var genTasks []*astStructForGen

	packageName := astInFile.Name.Name
	i.Nodes(iFilter, func(node ast.Node, push bool) (proceed bool) {
		genDecl := node.(*ast.GenDecl)
		if genDecl.Doc == nil {
			return false
		}

		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return false
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false
		}

		for _, comment := range genDecl.Doc.List {
			switch comment.Text {
			case "//go:genopt":
				genTasks = append(genTasks, &astStructForGen{
					packageName: packageName,
					typeSpec:    typeSpec,
					structType:  structType,
				})
			default:
				return
			}
		}

		return false
	})

	astOutFile := &ast.File{
		Name: astInFile.Name,
	}

	for _, task := range genTasks {
		err = task.generate(astOutFile)
		if err != nil {
			log.Fatalf("generate error: %s", err.Error())
		}
	}

	outFile, err := os.Create(
		strings.TrimSuffix(fileFullPath, ".go") +
			"_gen.go",
	)
	if err != nil {
		log.Fatalf("create file error: %s", err.Error())
	}

	err = printer.Fprint(outFile, token.NewFileSet(), astOutFile)
	if err != nil {
		log.Fatalf("print file error: %s", err.Error())
	}

	return genTasks, nil
}

// Block with templates.
// CC - CamelCase
// LCC - lowerCamelCase
//nolint:gochecknoglobals
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
		{{ .OptionNameLCC }}.{{ .OptionNameLCC }} = {{ .OptionNameLCC }}
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

func (g *astStructForGen) generate(outFile *ast.File) error {
	// required, err := g.requiredField()
	// if err != nil {
	// 	return fmt.Errorf("can not to generate file, error: %w", err)
	// }
	//
	// params := struct {
	// 	ObjectNameLCC string
	// 	ObjectNameCC  string
	// }{
	// 	ObjectNameCC:  g.typeSpec.Name.Name,
	// 	ObjectNameLCC: exprToString(required.Type),
	// }

	buf := new(bytes.Buffer)
	err := error(nil)

	if err = g.writePackage(buf); err != nil {
		return fmt.Errorf("write package error: %w", err)
	}

	if err = g.writeOptionType(buf); err != nil {
		return fmt.Errorf("write options error: %w", err)
	}

	if err = g.writeInit(buf); err != nil {
		return fmt.Errorf("write init error: %w", err)
	}

	templateAst, err := g.generateAST(buf.Bytes())

	for _, decl := range templateAst.Decls {
		outFile.Decls = append(outFile.Decls, decl)
	}

	return nil
}

func (g *astStructForGen) writePackage(buf io.Writer) error {
	err := templatePackage.Execute(buf, struct {
		PackageName string
	}{
		PackageName: g.packageName,
	})
	if err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	return nil
}

func (g *astStructForGen) writeOptionType(buf io.Writer) error {
	// for _, field := range g.structType.Fields.List {
	//
	// }
	err := templateOptionFunc.Execute(buf, struct {
		ObjectNameCC string
	}{
		ObjectNameCC: g.typeSpec.Name.Name,
	})
	if err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	return nil
}

func (g *astStructForGen) writeInit(buf io.Writer) error {
	err := templateInit.Execute(buf, struct {
		ObjectNameLCC string
		ObjectNameCC  string
	}{
		ObjectNameCC: g.typeSpec.Name.Name,
		ObjectNameLCC: strcase.LowerCamelCase(
			g.typeSpec.Name.Name,
		),
	})
	if err != nil {
		return fmt.Errorf("execute template error: %w", err)
	}

	return nil
}

func (g *astStructForGen) generateAST(listing []byte) (*ast.File, error) {
	templateAst, err := parser.ParseFile(
		token.NewFileSet(),
		"",
		listing,
		parser.ParseComments,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can not to create ast for template, error: %w", err,
		)
	}

	return templateAst, nil
}

func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), expr)
	if err != nil {
		log.Fatalf(
			"can not to print type to string, error: %s",
			err.Error(),
		)
	}

	return buf.String()
}

// func generateFileWithOptions(gen *astStructForGen) (*ast.File, error) {
// 	out := new(ast.File)
// 	out.Package
// 	for _, field := range gen.structType.Fields.List {
// 		field.Type.Pos()
// 	}
// }
//
// func (s *astStructForGen) providePackageName() {
//
// }
//

func (g *astStructForGen) requiredField() (*ast.Field, error) {
	for _, field := range g.structType.Fields.List {
		if !strings.Contains(field.Tag.Value, "required") {
			continue
		}

		return field, nil
	}

	return nil, fmt.Errorf("required field not found")
}
