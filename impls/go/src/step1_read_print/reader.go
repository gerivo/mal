package main

import (
	"regexp"
	"strconv"
)

type reader struct {
	tokens   []string
	position int
}

func (r reader) in_bounds() bool {
	return r.position < len(r.tokens)
}

func (r reader) peek() string {
	return r.tokens[r.position]
}

func (r reader) next() string {
	before := r.position
	r.position = r.position + 1

	return r.tokens[before]
}

func read_str(input string) []any {
	return read_form(tokenize(input))
}

func tokenize(input string) reader {
	r, _ := regexp.Compile("[\\s,]*(~@|[\\[\\]{}()'`~^@]|\"(?:\\.|[^\\\"])*\"?|;.*|[^\\s\\[\\]{}('\"`,;)]*)")

	return reader{r.FindAllString(input, -1), 0}
}

func read_form(r reader) []any {
	first := r.peek()
	var ret []any

	if first[0] == '(' {
		ret = append(ret, []any{read_list(r)})
	} else {
		ret = append(ret, read_atom(r))
	}

	return ret
}

func read_list(r reader) []any {
	first := r.next()
	var ret []any

	for first[0] != ')' && r.in_bounds() {
		ret = append(ret, read_form(r))
		first = r.next()
	}

	return ret
}

func read_atom(r reader) []any {
	first := r.peek()
	var ret []any

	value_int, is := strconv.ParseInt(first, 10, 0)
	symbol := true

	if is != nil {
		ret = append(ret, value_int)
		symbol = false
	}

	value_bool, is := strconv.ParseBool(first)

	if is != nil {
		ret = append(ret, value_bool)
		symbol = false
	}

	if symbol {
		ret = append(ret, first)
	}

	return ret
}
