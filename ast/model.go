package ast

type Expr interface {
}

type NumberExpr struct {
	Value int
}

func (numExpr NumberExpr) Eval() int {
	return numExpr.Value
}

type VariableExpr struct {
	Name string
}

type BinaryExpr struct {
	Op       int
	LHS, RHS Expr
}

type CallExpr struct {
	Callee string
	Args   []Expr
}

type PrototypeExpr struct {
	Name string
	Args []string
}

type FunctionExpr struct {
	Proto PrototypeExpr
	Body  Expr
}
