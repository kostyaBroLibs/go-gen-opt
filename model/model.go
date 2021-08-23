package model

import (
	"go/ast"
)

type ParsedStruct struct {
	TypeSpec    *ast.TypeSpec
	StructType  *ast.StructType
	PackageName string
}
