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
	Data []MalTyper
}

func (ml MalList) GetValue() any {
	return ml.Data
}

func (ml MalList) String() string {
	str := make([]string, 0)
	for _, value := range ml.Data {
		str = append(str, value.String())
	}

	return "(" + strings.Join(str, " ") + ")"
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
