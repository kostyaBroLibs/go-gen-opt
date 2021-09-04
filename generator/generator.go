package genarator

import (
	"bytes"
	"fmt"
	"github.com/stoewer/go-strcase"
	"go-gen-opt/model"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
)

type Config struct {
	Input model.ParsedStruct
}

type Generator struct {
	input model.ParsedStruct
}

func NewGenerator(config *Config) *Generator {
	return &Generator{
		input: config.Input,
	}
}

func (g *Generator) generate(outFile *ast.File) error {
	buf := new(bytes.Buffer)
	err := error(nil)

	if err = g.writePackage(buf); err != nil {
		return fmt.Errorf("write package error: %w", err)
	}

	if err = g.writeOptionType(buf); err != nil {
		return fmt.Errorf("write options error: %w", err)
	}

	if err = g.writeOptionFuncs(buf); err != nil {
		return fmt.Errorf("write options funcs error: %w", err)
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

func (g *Generator) writePackage(buf io.Writer) error {
	err := templatePackage.Execute(buf, struct {
		PackageName string
	}{
		PackageName: g.input.PackageName,
	})
	if err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	return nil
}

func (g *Generator) writeOptionType(buf io.Writer) error {
	err := templateOptionFunc.Execute(buf, struct {
		ObjectNameCC string
	}{
		ObjectNameCC: strcase.UpperCamelCase(g.input.TypeSpec.Name.Name),
	})
	if err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	return nil
}

func (g *Generator) writeOptionFuncs(buf io.Writer) error {
	for _, field := range g.input.StructType.Fields.List {
		err := templateWithFunc.Execute(buf, struct {
			OptionNameCC  string
			OptionNameOC  string
			OptionNameLCC string
			OptionType    string
			ObjectNameCC  string
			ObjectNameLCC string
		}{
			ObjectNameCC:  strcase.UpperCamelCase(g.input.TypeSpec.Name.Name),
			ObjectNameLCC: strcase.LowerCamelCase(g.input.TypeSpec.Name.Name),
			OptionType:    exprToString(field.Type),
			OptionNameCC:  strcase.UpperCamelCase(field.Names[0].Name),
			OptionNameOC:  field.Names[0].Name,
			OptionNameLCC: strcase.LowerCamelCase(field.Names[0].Name),
		})
		if err != nil {
			return fmt.Errorf("template execute error: %w", err)
		}
	}

	return nil
}

func (g *Generator) writeInit(buf io.Writer) error {
	err := templateInit.Execute(buf, struct {
		ObjectNameLCC string
		ObjectNameCC  string
	}{
		ObjectNameCC: g.input.TypeSpec.Name.Name,
		ObjectNameLCC: strcase.LowerCamelCase(
			g.input.TypeSpec.Name.Name,
		),
	})
	if err != nil {
		return fmt.Errorf("execute template error: %w", err)
	}

	return nil
}

func (g *Generator) generateAST(listing []byte) (*ast.File, error) {
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
