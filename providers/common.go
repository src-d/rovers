package providers

import (
	"fmt"
	"strings"
)

const (
	isForkKey = "is_fork"
	aliasesKey = "aliases"
)

type ContextBuilder map[string]string

func (c ContextBuilder) Fork(isFork bool) ContextBuilder {
	c[isForkKey] = fmt.Sprintf("%s", isFork)

	return c
}

func (c ContextBuilder) Aliases(aliases []string) ContextBuilder {
	a := strings.Join(aliases, ", ")
	c[aliasesKey] = a

	return c
}
