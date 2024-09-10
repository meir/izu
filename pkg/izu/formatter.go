package izu

import "embed"

type State uint8

const (
	StateKeybind State = iota
	StateCommand
	StateBase
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

type Formatter interface{}

//go:embed formatters/*
var Formatters embed.FS
