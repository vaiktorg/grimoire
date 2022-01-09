package dice

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type T byte

const (
	Unknown T = iota
	OpType
	ValType
	ScopeType
	DiceType
)

type MacroProcessor struct {
	proc Operation
	dice *Dice
	exp  string
}

func (m *MacroProcessor) Eval(exp string) float64 {
	//rgx := regexp.MustCompile(`(?:\$D\d{1,3}|[0-9]{1,3}\.\d{1,4}|\d+)|[-+/*]|\(|\)`)
	rgx := regexp.MustCompile(`\$D\d+|[()]|[/+*-]|\d+\.\d+|\d+`)
	toks := rgx.FindAllString(exp, -1)

	depth := 0
	for _, t := range toks {
		switch is(t) {
		case ScopeType:
			fmt.Printf(t+"\t%v_%v\n", ScopeType, depth)
		case OpType:
			fmt.Printf(t+"\t%v\n", OpType)
		case ValType:
			fmt.Printf(t+"\t%v\n", ValType)
		case DiceType:
			fmt.Printf(t+"\t%v\n", DiceType)
		default:
			return -1
		}
	}

	return 1.2
}

func is(tok string) T {
	if strings.Contains(tok, "$D") {
		return DiceType
	}
	if strings.ContainsAny(tok, "+/*^-") {
		return OpType
	}
	if strings.ContainsAny(tok, "()") {
		return ScopeType
	}
	if strings.ContainsAny(tok, "1234567890.") {
		return ValType
	}
	return Unknown
}
func (m *MacroProcessor) diceVal(tok string) float64 {
	v, _ := strconv.Atoi(tok[0:2])
	return float64(m.dice.getDiceVal(v))
}
