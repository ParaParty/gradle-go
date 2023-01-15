package main

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
	"gradle-go-generated/parser"
	"math"
	"strconv"
)

type F64Stack struct {
	stack *stack.Stack
}

func (s F64Stack) Push(value float64) {
	s.stack.Push(value)
}

func (s F64Stack) Pop() float64 {
	value, ok := s.stack.Peek().(float64)
	if !ok {
		panic("fail to cast to f64")
	}
	s.stack.Pop()
	return value
}

func (s F64Stack) Pop2() (float64, float64) {
	return s.Pop(), s.Pop()
}

func NewF64Stack() F64Stack {
	return F64Stack{stack: stack.New()}
}

type Inst uint

// https://wizardforcel.gitbooks.io/antlr4-short-course/content/calculator-listener.html
const (
	LDV Inst = iota // 变量入栈
	LDC             // 常量入栈
	DEF             // 栈顶一个元素存入指定变量
	ADD             // 栈顶两个元素出栈，求和后入栈
	SUB             // 栈顶两个元素出栈，求差后入栈
	MUL             // 栈顶两个元素出栈，求积后入栈
	MOD             // 栈顶两个元素出栈，求模后入栈
	DIV             // 栈顶两个元素出栈，求商后入栈
	POW             // 栈顶两个元素出栈，乘方后入栈
	RET             // 栈顶一个元素出栈，计算结束
)

type InstV struct {
	Inst  Inst
	Value float64
	Id    string
}

type InstSet struct {
	instSet      []InstV
	ResolveIdent func(id string) float64
	SetIdent     func(id string, value float64)
}

func (i *InstSet) Clear() {
	i.instSet = make([]InstV, 0)
}

func (i *InstSet) Evaluate() float64 {
	reg := NewF64Stack()

	for _, item := range i.instSet {
		switch item.Inst {
		case LDV:
			{
				reg.Push(i.ResolveIdent(item.Id))
			}
		case LDC:
			{
				reg.Push(item.Value)
			}
		case DEF:
			{
				value := reg.Pop()
				i.SetIdent(item.Id, value)
			}
		case ADD:
			{
				v1, v2 := reg.Pop2()
				reg.Push(v1 + v2)
			}
		case SUB:
			{
				v1, v2 := reg.Pop2()
				reg.Push(v1 - v2)
			}
		case MUL:
			{
				v1, v2 := reg.Pop2()
				reg.Push(v1 * v2)
			}
		case DIV:
			{
				v1, v2 := reg.Pop2()
				reg.Push(v1 / v2)
			}
		case MOD:
			{
				rhs, lhs := reg.Pop2()
				lhsInt := int64(lhs)
				rhsInt := int64(rhs)
				if math.Abs(float64(lhsInt)-lhs) < 1e-6 && math.Abs(float64(rhsInt)-rhs) < 1e-6 {
					reg.Push(float64(lhsInt % rhsInt))
				}

				panic("the value near mod '%' operator must be integers")
			}
		case POW:
			{
				v1, v2 := reg.Pop2()
				reg.Push(math.Pow(v1, v2))
			}
		case RET:
			{
				fmt.Printf("%v\n", reg.Pop())
			}
		}

	}

	return 0
}

func (i *InstSet) AddInst(inst InstV) {
	i.instSet = append(i.instSet, inst)
}

func (i *InstSet) Print() {
	for _, item := range i.instSet {
		fmt.Printf("%5v %5v %5v\n", item.Inst, item.Value, item.Id)
	}
}

type CalcListener struct {
	*parser.BaseCalcListener

	memory  map[string]float64
	instSet InstSet
}

func (s *CalcListener) ExitCal_stat(ctx *parser.Cal_statContext) {
	s.instSet.AddInst(InstV{Inst: RET})
}

func (s *CalcListener) ExitPrint(ctx *parser.PrintContext) {}

func (s *CalcListener) ExitAssign(ctx *parser.AssignContext) {
	id := ctx.ID().GetText()
	s.instSet.AddInst(InstV{Inst: DEF, Id: id})
}

func (s *CalcListener) ExitMd_expr(ctx *parser.Md_exprContext) {
	if ctx.GetOp().GetTokenType() == parser.CalcParserMUL {
		s.instSet.AddInst(InstV{Inst: MUL})
	} else if ctx.GetOp().GetTokenType() == parser.CalcParserDIV {
		s.instSet.AddInst(InstV{Inst: DIV})
	} else {
		s.instSet.AddInst(InstV{Inst: MOD})
	}
}

func (s *CalcListener) ExitAs_expr(ctx *parser.As_exprContext) {
	if ctx.GetOp().GetTokenType() == parser.CalcParserADD {
		s.instSet.AddInst(InstV{Inst: ADD})
	} else {
		s.instSet.AddInst(InstV{Inst: SUB})
	}
}

func (s *CalcListener) ExitPow_expr(ctx *parser.Pow_exprContext) {
	s.instSet.AddInst(InstV{Inst: POW})
}

func (s *CalcListener) ExitId(ctx *parser.IdContext) {
	s.instSet.AddInst(InstV{Inst: LDV, Id: ctx.ID().GetText()})

}

func (s *CalcListener) ExitNumber(ctx *parser.NumberContext) {
	childrenCount := ctx.GetChildCount()
	value, err := strconv.ParseFloat(ctx.NUM().GetText(), 64)
	if err != nil {
		panic(err)
	}
	if childrenCount == 2 && ctx.GetSign().GetTokenType() == parser.CalcParserSUB {
		s.instSet.AddInst(InstV{Inst: LDC, Value: -value})
	} else {
		s.instSet.AddInst(InstV{Inst: LDC, Value: value})
	}
}

func CreateCalcListener() *CalcListener {
	ret := &CalcListener{
		memory:  make(map[string]float64, 0),
		instSet: InstSet{},
	}

	ret.instSet.ResolveIdent = func(id string) float64 {
		ret, ok := ret.memory[id]
		if !ok {
			panic(fmt.Sprintf("identity %v not defined", id))
		}
		return ret
	}

	ret.instSet.SetIdent = func(id string, value float64) {
		ret.memory[id] = value
	}

	return ret
}
