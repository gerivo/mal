package main

import (
	"strconv"
	"strings"
)

type MalTyper interface {
	GetValue() any
	String() string
}

type MalList struct {
	Data  []MalTyper
	Wraps [2]string
}

func (ml MalList) GetValue() any {
	return ml.Data
}

func (ml MalList) String() string {
	str := make([]string, 0)
	for _, value := range ml.Data {
		str = append(str, value.String())
	}

	return ml.Wraps[0] + strings.Join(str, " ") + ml.Wraps[1]
}

type MalInt struct {
	value int
}

func (ml MalInt) GetValue() any {
	return ml.value
}

func (ml MalInt) String() string {
	return strconv.Itoa(ml.value)
}

type MalSymbol struct {
	symbol string
}

func (ml MalSymbol) GetValue() any {
	return ml.symbol
}

func (ml MalSymbol) String() string {
	return ml.symbol
}

type MalString struct {
	symbol   string
	readably bool
}

func (ml MalString) GetValue() any {
	return ml.symbol
}

func (ml MalString) String() string {
	return ml.symbol
}
