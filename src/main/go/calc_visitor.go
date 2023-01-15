package main

import (
	"fmt"
	"gradle-go-generated/parser"
	"math"
	"strconv"
)

type CalcVisitor struct {
	*parser.BaseCalcVisitor

	memory map[string]float64
}

func (v *CalcVisitor) VisitProg(ctx *parser.ProgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *CalcVisitor) VisitCal_stat(ctx *parser.Cal_statContext) interface{} {
	value := v.VisitChildren(ctx.Expr())
	ret, ok := value.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}
	return ret
}

func (v *CalcVisitor) VisitAssign(ctx *parser.AssignContext) interface{} {
	id := ctx.ID().GetText()
	value := v.VisitChildren(ctx.Expr())
	ret, ok := value.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}
	v.memory[id] = ret
	return ret
}

func (v *CalcVisitor) VisitPrint(ctx *parser.PrintContext) interface{} {
	value := v.VisitChildren(ctx.Expr())
	ret, ok := value.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}
	fmt.Printf("%v\n", ret)
	return ret
}

func (v *CalcVisitor) VisitNumber(ctx *parser.NumberContext) interface{} {
	childrenCount := ctx.GetChildCount()
	value, err := strconv.ParseFloat(ctx.NUM().GetText(), 64)
	if err != nil {
		return err
	}
	if childrenCount == 2 && ctx.GetSign().GetTokenType() == parser.CalcParserSUB {
		return -value
	}
	return value
}

func (v *CalcVisitor) VisitParens(ctx *parser.ParensContext) interface{} {
	return v.VisitChildren(ctx.Expr())
}

func (v *CalcVisitor) VisitAs_expr(ctx *parser.As_exprContext) interface{} {
	lhsValue := v.VisitChildren(ctx.Expr(0))
	lhs, ok := lhsValue.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}

	rhsValue := v.VisitChildren(ctx.Expr(1))
	rhs, ok := rhsValue.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}

	if ctx.GetOp().GetTokenType() == parser.CalcParserADD {
		return lhs + rhs
	} else {
		return lhs - rhs
	}
}

func (v *CalcVisitor) VisitId(ctx *parser.IdContext) interface{} {
	id := ctx.ID().GetText()
	value, ok := v.memory[id]
	if ok {
		return value
	} else {
		fmt.Printf("undefined identifier %v\n", id)
		return 0
	}
}

func (v *CalcVisitor) VisitMd_expr(ctx *parser.Md_exprContext) interface{} {
	lhsValue := v.VisitChildren(ctx.Expr(0))
	lhs, ok := lhsValue.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}

	rhsValue := v.VisitChildren(ctx.Expr(1))
	rhs, ok := rhsValue.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}

	if ctx.GetOp().GetTokenType() == parser.CalcParserMUL {
		return lhs * rhs
	}

	if ctx.GetOp().GetTokenType() == parser.CalcParserDIV {
		if math.Abs(rhs) < 1e-6 {
			panic("divided by zero")
		}
		return lhs / rhs
	}

	lhsInt := int64(lhs)
	rhsInt := int64(rhs)
	if math.Abs(float64(lhsInt)-lhs) < 1e-6 && math.Abs(float64(rhsInt)-rhs) < 1e-6 {
		return float64(lhsInt % rhsInt)
	}

	panic("the value near mod '%' operator must be integers")
}

func (v *CalcVisitor) VisitPow_expr(ctx *parser.Pow_exprContext) interface{} {
	truthValue := v.VisitChildren(ctx.Expr(0))
	truth, ok := truthValue.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}

	powerValue := v.VisitChildren(ctx.Expr(1))
	power, ok := powerValue.(float64)
	if !ok {
		return fmt.Errorf("cast to f64 failed")
	}

	return math.Pow(truth, power)
}
