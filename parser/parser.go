package parser

import (
	"fmt"
	"go-gen-opt/model"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/inspector"
)

type Config struct {
	FileFullPath string
}

type Parser struct {
	// iAst is the AST(abstract syntax tree) for input file.
	iAst *ast.File
}

func NewParser(config Config) (*Parser, error) {
	astFile, err := parser.ParseFile(
		token.NewFileSet(),
		config.FileFullPath,
		nil,
		parser.ParseComments|parser.AllErrors,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can not to parse file '%s', error: %w",
			config.FileFullPath, err,
		)
	}

	return &Parser{
		iAst: astFile,
	}, nil
}

func (p Parser) ParseStructs(tag string) []*model.ParsedStruct {
	i := inspector.New([]*ast.File{p.iAst})
	iFilter := []ast.Node{
		&ast.GenDecl{},
	}
	output := make([]*model.ParsedStruct, 0, 2)
	packageName := p.iAst.Name.Name

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
			case tag:
				output = append(output, &model.ParsedStruct{
					PackageName: packageName,
					TypeSpec:    typeSpec,
					StructType:  structType,
				})
			default:
				return
			}
		}

		return false
	})

	return output
}
