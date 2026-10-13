package backend

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"strconv"
	"strings"
)

type CalcResult struct {
	Expression string `json:"expression"`
	Result     string `json:"result"`
	Success    bool   `json:"success"`
}

func Evaluate(expr string) CalcResult {
	cleaned := strings.TrimSpace(expr)
	if cleaned == "" {
		return CalcResult{Expression: expr, Success: false}
	}

	// Must contain at least one math operator or known math function to trigger
	hasOperator := strings.ContainsAny(cleaned, "+-*/%^")
	hasFunc := strings.Contains(cleaned, "sqrt(") || strings.Contains(cleaned, "abs(") ||
		strings.Contains(cleaned, "sin(") || strings.Contains(cleaned, "cos(") ||
		strings.Contains(cleaned, "tan(") || strings.Contains(cleaned, "log(")

	if !hasOperator && !hasFunc {
		return CalcResult{Expression: expr, Success: false}
	}

	// Replace ^ with pow syntax or handle
	cleaned = strings.ReplaceAll(cleaned, "^", "**")

	val, err := evalExprString(cleaned)
	if err != nil {
		return CalcResult{Expression: expr, Success: false}
	}

	// Format result
	var resStr string
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return CalcResult{Expression: expr, Success: false}
	}

	if val == math.Trunc(val) {
		resStr = fmt.Sprintf("%.0f", val)
	} else {
		resStr = fmt.Sprintf("%.6f", val)
		resStr = strings.TrimRight(strings.TrimRight(resStr, "0"), ".")
	}

	return CalcResult{
		Expression: expr,
		Result:     resStr,
		Success:    true,
	}
}

func evalExprString(s string) (float64, error) {
	node, err := parser.ParseExpr(s)
	if err != nil {
		return 0, err
	}
	return evalNode(node)
}

func evalNode(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.BasicLit:
		if n.Kind == token.INT || n.Kind == token.FLOAT {
			return strconv.ParseFloat(n.Value, 64)
		}
		return 0, fmt.Errorf("unsupported literal: %s", n.Value)

	case *ast.ParenExpr:
		return evalNode(n.X)

	case *ast.UnaryExpr:
		val, err := evalNode(n.X)
		if err != nil {
			return 0, err
		}
		if n.Op == token.SUB {
			return -val, nil
		}
		if n.Op == token.ADD {
			return val, nil
		}
		return 0, fmt.Errorf("unsupported unary op: %s", n.Op)

	case *ast.BinaryExpr:
		left, err := evalNode(n.X)
		if err != nil {
			return 0, err
		}
		right, err := evalNode(n.Y)
		if err != nil {
			return 0, err
		}

		switch n.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		case token.REM:
			if right == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			return math.Mod(left, right), nil
		}

	case *ast.CallExpr:
		ident, ok := n.Fun.(*ast.Ident)
		if !ok {
			return 0, fmt.Errorf("unsupported call")
		}
		args := make([]float64, len(n.Args))
		for i, arg := range n.Args {
			v, err := evalNode(arg)
			if err != nil {
				return 0, err
			}
			args[i] = v
		}

		switch strings.ToLower(ident.Name) {
		case "sqrt":
			if len(args) != 1 {
				return 0, fmt.Errorf("sqrt requires 1 arg")
			}
			return math.Sqrt(args[0]), nil
		case "abs":
			if len(args) != 1 {
				return 0, fmt.Errorf("abs requires 1 arg")
			}
			return math.Abs(args[0]), nil
		case "sin":
			if len(args) != 1 {
				return 0, fmt.Errorf("sin requires 1 arg")
			}
			return math.Sin(args[0]), nil
		case "cos":
			if len(args) != 1 {
				return 0, fmt.Errorf("cos requires 1 arg")
			}
			return math.Cos(args[0]), nil
		case "pow":
			if len(args) != 2 {
				return 0, fmt.Errorf("pow requires 2 args")
			}
			return math.Pow(args[0], args[1]), nil
		}

	case *ast.Ident:
		switch strings.ToLower(n.Name) {
		case "pi":
			return math.Pi, nil
		case "e":
			return math.E, nil
		}
	}

	return 0, fmt.Errorf("unsupported expression")
}
