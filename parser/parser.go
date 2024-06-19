package parser

import (
	"strings"
)

type Parser struct {
	s string
}

func New(s string) *Parser {
	return &Parser{s}
}

func (p Parser) Remaining() string {
	return p.s
}

func (p *Parser) ConsumeUntil(chars ...rune) string {
	minIdx := len(p.s)
	for _, c := range chars {
		idx := strings.IndexRune(p.s, c)
		if idx != -1 && idx < minIdx {
			minIdx = idx
		}
	}
	ret := p.s[:minIdx]
	p.advance(minIdx + 1)
	return ret
}

func (p *Parser) Consume(prefix string) (consumed bool) {
	if strings.HasPrefix(p.s, prefix) {
		p.advance(len(prefix))
		consumed = true
	} else {
		consumed = false
	}
	return
}

func (p *Parser) ConsumeN(n int) string {
	ret := p.s[:n]
	p.advance(n)
	return ret
}

func (p *Parser) advance(n int) {
	p.s = p.s[n:]
}
