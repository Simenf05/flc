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

func (g *Generator) NumberExprCodegen(numExpr ast.NumberExpr, block *ir.Block) value.Value {
	return constant.NewInt(types.I8, int64(numExpr.Value))
}

func (g *Generator) BinaryExprCodegen(binExpr ast.BinaryExpr, block *ir.Block) value.Value {
	left := g.Codegen(binExpr.Left, block)
	right := g.Codegen(binExpr.Right, block)

	if left == nil || right == nil {
		return nil
	}

	switch binExpr.Op {
	case ast.Addition:
		return block.NewAdd(left, right)
	case ast.Subtraction:
		return block.NewSub(left, right)
	case ast.Multiplication:
		return block.NewMul(left, right)
	case ast.Division:
		return block.NewSDiv(left, right)
	default:
		return nil
	}
}

func (g *Generator) CallExprCodegen(callExpr ast.CallExpr, block *ir.Block) value.Value {
	for _, callee := range g.module.Funcs {
		if callee.Name() != callExpr.Callee {
			continue
		}

		var args []value.Value

		for _, expr := range callExpr.Args {
			arg := g.Codegen(expr, block)
			args = append(args, arg)
		}

		return block.NewCall(callee, args...)

	}
	return nil
}

func (g *Generator) VariableExprCodegen(varExpr ast.VariableExpr, block *ir.Block) value.Value {
	return nil
}

func (g *Generator) ReturnExprCodegen(retExpr ast.ReturnExpr, block *ir.Block) value.Value {
	var value value.Value
	if retExpr.Value != nil {
		value = g.Codegen(retExpr.Value, block)
		if value == nil {
			return nil
		}
	}

	block.NewRet(value)
	return block
}

func (g *Generator) ProtoExprHandle(protoExpr ast.PrototypeExpr) *ir.Func {
	newFunc := g.module.NewFunc(protoExpr.Name, types.I8)
	return newFunc
}

func (g *Generator) FuncExprCodegen(funcExpr ast.FunctionExpr) value.Value {
	newFunc := g.ProtoExprHandle(funcExpr.Proto)
	block := newFunc.NewBlock("")

	value := g.Codegen(funcExpr.Body, block)
	if value == nil {
		return nil
	}

	return newFunc
}

func (g *Generator) Codegen(expr ast.Expr, block *ir.Block) value.Value {
	switch e := expr.(type) {
	case ast.ReturnExpr:
		return g.ReturnExprCodegen(e, block)
	case ast.BinaryExpr:
		return g.BinaryExprCodegen(e, block)
	case ast.CallExpr:
		return g.CallExprCodegen(e, block)
	case ast.NumberExpr:
		return g.NumberExprCodegen(e, block)
	case ast.VariableExpr:
		return g.VariableExprCodegen(e, block)
	}

	return nil
}
