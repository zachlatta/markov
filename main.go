package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Prefix []string

func (p Prefix) String() string {
	return strings.Join(p, " ")
}

func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

type Suffix struct {
	Word  string
	Count int
}

func (s *Suffix) String() string {
	return s.Word
}

type Chain struct {
	chain     map[string][]*Suffix
	prefixLen int
}

func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]*Suffix), prefixLen}
}

func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}

		key, suffixExists := p.String(), false
		for _, suffix := range c.chain[key] {
			if suffix.String() == s {
				suffix.Count++
				suffixExists = true
			}
		}

		if !suffixExists {
			c.chain[key] = append(c.chain[key], &Suffix{s, 1})
		}

		p.Shift(s)
	}
}

func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		suffixes := c.chain[p.String()]
		if len(suffixes) == 0 {
			break
		}

		next := suffixes[rand.Intn(len(suffixes))]
		words = append(words, next.Word)
		p.Shift(next.Word)
	}

	return strings.Join(words, " ")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	c := NewChain(2)
	c.Build(os.Stdin)
	fmt.Println(c.Generate(100))
}
