package dice

import (
	"math"
	"strconv"
)

const (
	Sum = "+"
	Mul = "*"
	Sub = "-"
	Div = "/"
	Exp = "^"
)

type Operation struct {
	l    *Operation
	data string
	r    *Operation
}

func (op Operation) Val() float64 {
	v, err := strconv.ParseFloat(op.data, 64)
	if err != nil {
		return 0
	}
	return v
}

func (op *Operation) Eval() float64 {
	switch op.data {
	case Div, Mul, Sub, Sum:
		return operate(op.l.Eval(), op.r.Eval(), op.data)
	}
	return op.Val()
}

func operate(val1, val2 float64, op string) float64 {
	switch op {
	case Div:
		return val1 / val2
	case Sum:
		return val1 + val2
	case Sub:
		return val1 - val2
	case Mul:
		return val1 * val2
	case Exp:
		return math.Pow(val1, val2)
	default:
		return -1
	}
}
