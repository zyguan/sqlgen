package util

import "github.com/zyguan/sqlgen/exprgen"

type TransformContext struct {
	Cols []Column
	ReplaceChildIdx int
}

type TransformRule interface {
	OneStep(node Expr, ctx TransformContext) []Expr
}

var rules []TransformRule

type ConstantToColumn struct {}

func (c *ConstantToColumn) OneStep(node Expr, ctx TransformContext) []Expr {
	var result []Expr
	switch _ := node.(type) {
	case *Func:
	case *Column:
	case *Constant:
		for _, col := range ctx.Cols {
			result = append(result, col)
		}
	}
	return result
}

type ColumnToConstant struct {}

func (c *ColumnToConstant) OneStep(node Expr, ctx TransformContext) []Expr {
	var result []Expr
	switch _ := node.(type) {
	case *Func:
	case *Constant:
	case *Column:
		result = append(result, exprgen.GenConstant(TypeMask(node.RetType())))
	}
	return result
}

type ReplaceChildToConstant struct {}

func (r *ReplaceChildToConstant) OneStep(node Expr, ctx TransformContext) []Expr {
	var result []Expr
	switch e := node.(type) {
	case *Constant:
	case *Column:
	case *Func:
		if len(e.children) > ctx.ReplaceChildIdx {
			newNode := e.Clone().(*Func)
			newNode.children[ctx.ReplaceChildIdx] = exprgen.GenConstant(TypeMask(newNode.children[ctx.ReplaceChildIdx].RetType()))
		}
	}
	return result
}

