package main

import (
	"errors"
	"regexp"
	"strconv"
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

func (r *reader) Next() (string, error) {
	r.Advance()

	if r.InBounds() {
		return r.Peek(), nil
	} else {
		return "", errors.New("unbalanced")
	}
}

func (r *reader) Advance() reader {
	r.position = r.position + 1
	return *r
}

func read_str(input string) (MalTyper, error) {
	return read_form(tokenize(input))
}

func tokenize(input string) *reader {
	r := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)

	matches := r.FindAllStringSubmatch(input, -1)

	x := new(reader)
	for _, captures := range matches {
		x.tokens = append(x.tokens, captures[1]) // the tokens are in the second capture group
	}

	x.position = 0
	return x
}

func read_form(r *reader) (MalTyper, error) {
	first := r.Peek()

	if first == "(" {
		list, ok := read_list(r)

		return list, ok
	} else {
		return read_atom(r), nil
	}
}

func read_list(r *reader) (MalTyper, error) {
	first, ok := r.Next()
	list := new(MalList)

	for first != ")" && ok == nil {
		app, err := read_form(r)
		if err != nil {
			return list, nil
		}

		list.Data = append(list.Data, app)
		first, ok = r.Next()
	}

	return list, ok
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
