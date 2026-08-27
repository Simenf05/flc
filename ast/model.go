package ast

type OperatorType rune

const (
	Addition OperatorType = '+'
	Subtraction OperatorType = '-'
	Multiplication OperatorType = '*'
	Division OperatorType = '/'
)

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
	Op          OperatorType
	Left, Right Expr
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
