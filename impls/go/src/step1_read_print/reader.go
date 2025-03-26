package main

import (
	"regexp"
	"strconv"
	"strings"
)

type reader struct {
	tokens   []string
	position int
}

func (r reader) InBounds() bool {
	return r.position < len(r.tokens)
}

func (r reader) Peek() string {
	return r.tokens[r.position]
}

func (r *reader) Advance() reader {
	r.position = r.position + 1
	return *r
}

func read_str(input string) MalTyper {
	return read_form(tokenize(input))
}

func tokenize(input string) *reader {
	r := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)

	tokens := r.FindAllString(input, -1)
	for i := range tokens {
		tokens[i] = strings.TrimSpace(tokens[i])
	}

	x := new(reader)
	x.tokens = tokens
	x.position = 0
	return x
}

func read_form(r *reader) MalTyper {
	first := r.Peek()

	if first == "(" {
		list := read_list(r)
		return list
	} else {
		return read_atom(r)
	}
}

func read_list(r *reader) MalTyper {
	r.Advance()
	first := r.Peek()
	list := new(MalList)

	for first != ")" && r.InBounds() {
		list.Data = append(list.Data, read_form(r))

		r.Advance()
		first = r.Peek()
	}

	return list
}

func read_atom(r *reader) MalTyper {
	first := r.Peek() // first value always the read atom

	value, is := strconv.ParseInt(first, 10, 0)
	if is == nil {
		return MalInt{int(value)}
	}

	// TODO: MalTrue and MalFalse (they are different)
	return MalSymbol{first}
}
