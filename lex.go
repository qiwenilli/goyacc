package goyacc

import (
	"bufio"
	"log"
	"reflect"
	"strconv"
)

var (
	// eq	Equal	city eq 'Redmond'
	// ne	Not equal	city ne 'London'
	// gt	Greater than	price gt 20
	// ge	Greater than or equal	price ge 10
	// lt	Less than	price lt 20
	// le	Less than or equal	price le 100

	comparisonEq = []rune{'e', 'q'}
	comparisonNe = []rune{'n', 'e'}
	comparisonGt = []rune{'g', 't'}
	comparisonGe = []rune{'g', 'e'}
	comparisonLt = []rune{'l', 't'}
	comparisonLe = []rune{'l', 'e'}
)

type line struct {
	input string
	buf   *bufio.Reader
	data  interface{}
}

func (p *line) Lex(lval *yySymType) int {
	for {
		r := p.read()
		if r == 0 {
			return 0
		}
		switch r {
		case ' ':
			p.unread()
			return p.scanSymbol()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			// 负数 -100; 正数 100; 小数 -0.100
			p.unread()
			str := p.scanNumber()
			//
			lval.i = str
			return VAL
		case 39: //'
			str := p.scanString()
			//
			lval.i = str
			return VAL
		default:
			p.unread()
			key := p.scanKey()
			//
			lval.s = key
			return KEY
		}
	}
	return 0
}

func (p *line) Error(s string) {
	log.Println("syntax error: ", s, p.input)
}

func (p *line) read() rune {
	r, _, _ := p.buf.ReadRune()
	return r
}

func (p *line) scanSymbol() int {
	// 开始 & 结尾 都得有一个空格
	matchSymbol := func(symbol []rune) int {
		p.unread()
		for _, s := range symbol {
			r := p.read()
			if r != s {
				return 0
			}
		}
		//以空格结尾： ' eq '
		r := p.read()
		if r != ' ' {
			return 0
		}
		//去掉多余的后缀空格
		for {
			r = p.read()
			if r != ' ' {
				p.unread()
				break
			}
		}
		if reflect.DeepEqual(symbol, comparisonEq) {
			return EQ
		} else if reflect.DeepEqual(symbol, comparisonNe) {
			return NE
		} else if reflect.DeepEqual(symbol, comparisonGt) {
			return GT
		} else if reflect.DeepEqual(symbol, comparisonGe) {
			return GE
		} else if reflect.DeepEqual(symbol, comparisonLt) {
			return LT
		} else if reflect.DeepEqual(symbol, comparisonLe) {
			return LE
		}
		return 0
	}
	for {
		r := p.read()
		if r == ' ' {
			continue
		}
		switch r {
		case 'e':
			return matchSymbol(comparisonEq)
		case 'n':
			return matchSymbol(comparisonNe)
		case 'g':
			rr := p.read()
			if rr == 't' {
				p.unread()
				return matchSymbol(comparisonGt)
			} else {
				return matchSymbol(comparisonGe)
			}
		case 'l':
			rr := p.read()
			if rr == 't' {
				p.unread()
				return matchSymbol(comparisonLt)
			} else {
				return matchSymbol(comparisonLe)
			}
		default:
			return 0
		}
	}
	return 0
}

func (p *line) scanKey() string {
	//开始单引号 & 结束单引号
	var str []rune
	for {
		r := p.read()
		if r == ' ' {
			p.unread()
			break
		}
		str = append(str, r)
	}
	return string(str)
}

func (p *line) scanNumber() interface{} {
	//开始为0~9.-
	var str []rune
	var i int
	var isfloat bool
	for {
		r := p.read()
		if (r >= '0' && r <= '9') || r == '.' || (i == 0 && r == '-') {
			if r == '.' {
				isfloat = true
			}
			str = append(str, r)
		} else {
			break
		}
		i++
	}
	if isfloat {
		val, _ := strconv.ParseFloat(string(str), 64)
		return val
	}
	val, _ := strconv.ParseInt(string(str), 10, 64)
	//
	return val
}

func (p *line) scanString() string {
	//开始单引号 & 结束单引号
	var str []rune
	var prve rune
	for {
		r := p.read()
		// 39 '
		// 92 \
		if r == 39 && prve != 92 {
			break
		}
		prve = r
		//
		str = append(str, r)
	}

	return string(str)
}

func (p *line) unread() { _ = p.buf.UnreadRune() }
