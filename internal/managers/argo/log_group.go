package argo

import (
	"github.com/hanjunlee/argocui/pkg/argo"
)

type logGroup []argo.Log

func (g logGroup) Len() int {
	return len(g)
}

func (g logGroup) Less(i, j int) bool {
	return g[i].Time.Before(g[j].Time)
}

func (g logGroup) Swap(i, j int) {
	t := g[i]
	g[i] = g[j]
	g[j] = t
	return
}
