package codegen

import (
	"fmt"
	"simenf05/flc/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func NewGenerator() Generator {
	m := ir.NewModule()
	return Generator{
		module: m,
	}
}

func (g *Generator) Println() {
	fmt.Println(g.module)
}

func (g *Generator) NumberExprCodegen(numExpr ast.NumberExpr) value.Value {
	return constant.NewInt(types.I8, 0)
}

func (g *Generator) BinaryExprCodegen(binExpr ast.BinaryExpr) value.Value {
	left := g.Codegen(binExpr.Left)
	right := g.Codegen(binExpr.Right)

	if left == nil || right == nil {
		return nil
	}

	switch binExpr.Op {
	case ast.Addition:
		return ir.NewAdd(left, right)
	case ast.Subtraction:
		return ir.NewSub(left, right)
	case ast.Multiplication:
		return ir.NewMul(left, right)
	case ast.Division:
		return ir.NewSDiv(left, right)
	default:
		return nil
	}
}

func (g *Generator) CallExprCodegen(callExpr ast.CallExpr) value.Value {
	for _, callee := range g.module.Funcs {
		if callee.Name() != callExpr.Callee {
			continue
		}

		var args []value.Value

		for _, expr := range callExpr.Args {
			arg := g.Codegen(expr)
			args = append(args, arg)
		}

		return ir.NewCall(callee, args...)

	}
	return nil
}

func (g *Generator) VariableExprCodegen(varExpr ast.VariableExpr) value.Value {
	return nil
}

func (g *Generator) ProtoExprCodegen(protoExpr ast.PrototypeExpr) value.Value {

	newFunc := block.NewFunc(protoExpr.Name, types.I8)
	return newFunc
}

func (g *Generator) FuncExprCodegen(funcExpr ast.FunctionExpr) value.Value {
	newFunc := g.Codegen(funcExpr.Proto)
	newBody := g.Codegen(funcExpr.Body)

	irFunc, ok := newFunc.(*ir.IFunc)
	if !ok {
		return nil
	}

	return irFunc
}

func (g *Generator) Codegen(expr ast.Expr, block *ir.Block) value.Value {

	fmt.Println(expr)

	switch e := expr.(type) {
	case ast.BinaryExpr:
		return g.BinaryExprCodegen(e)
	case ast.CallExpr:
		return g.CallExprCodegen(e)
	case ast.NumberExpr:
		return g.NumberExprCodegen(e)
	case ast.VariableExpr:
		return g.VariableExprCodegen(e)
	case ast.FunctionExpr:
		return g.FuncExprCodegen(e)
	case ast.PrototypeExpr:
		return g.ProtoExprCodegen(e)
	}

	return nil
}
