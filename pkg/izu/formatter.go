package izu

import "embed"

type State uint8

const (
	StateBase State = iota
	StateMultiple
	StateSingle
	StateSinglePart
	StateString
)

type Part interface {
	Info() (State, []Part)
	Parse([]byte) (int, error)
	String() string
}

type Formatter interface {
	Format(State, []Part) []string
}

//go:embed formatters/*
var Formatters embed.FS
