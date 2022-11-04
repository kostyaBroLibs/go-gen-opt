package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

const tag = "go:genopt"

// StructInfo contains information about struct
// for which must be generated init function with config and options.
type StructInfo struct {
	Documentation string
	Name          string
	Fields        []FieldInfo
}

// FieldInfo contains information about field of struct.
type FieldInfo struct {
	Documentation string
	Name          string
	Type          string
	IsRequired    bool
	DefaultValue  string
}

// InfoAboutStructsInFile parses file and returns information about structs,
// for which must be generated init function with config and options.
func InfoAboutStructsInFile(filePath string) ([]StructInfo, error) {
	filePayload, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't read file %q: %w", filePath, err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, filePayload, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("can't parse file %q: %w", filePath, err)
	}

	var output []StructInfo

	for _, decl := range file.Decls {
		typeDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if !isStructForGeneration(typeDecl) {
			continue
		}

		structInfo := StructInfo{
			Documentation: typeDecl.Doc.Text(),
			Name:          typeDecl.Specs[0].(*ast.TypeSpec).Name.Name,
		}

		for _, field := range typeDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
			fieldInfo := FieldInfo{
				Documentation: field.Doc.Text(),
				Name:          field.Names[0].Name,
				Type:          field.Type.(*ast.Ident).Name,
			}

			if field.Tag != nil {
				tagValue := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
				optTag := tagValue.Get("opt")
				properties := strings.Split(optTag, " ")
				for _, property := range properties {
					if property == "required" {
						fieldInfo.IsRequired = true
					} else if strings.HasPrefix(property, "default=") {
						fieldInfo.DefaultValue = property[8:]
					}
				}
			}

			structInfo.Fields = append(structInfo.Fields, fieldInfo)
		}

		output = append(output, structInfo)
	}

	return output, nil
}

func isStructForGeneration(typeDecl *ast.GenDecl) bool {
	for _, docLine := range typeDecl.Doc.List {
		if strings.HasPrefix(docLine.Text, "//"+tag) ||
			strings.HasPrefix(docLine.Text, "// "+tag) {
			return true
		}
	}
	return false
}
