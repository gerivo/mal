package main

import (
	"errors"
	"regexp"
	"strconv"
)

type reader struct {
	tokens   []string
	position int
	readably bool
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

func read_str(input string, readably bool) (MalTyper, error) {
	return read_form(tokenize(input, readably))
}

func tokenize(input string, readably bool) *reader {
	r := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)

	matches := r.FindAllStringSubmatch(input, -1)

	x := new(reader)
	for _, captures := range matches {
		x.tokens = append(x.tokens, captures[1]) // the tokens are in the second capture group
	}

	x.position = 0
	x.readably = readably
	return x
}

type listContext struct {
	wrappers [2]string
	onlyOne  bool
}

func read_form(r *reader) (MalTyper, error) {
	first := r.Peek()

	// this map directs the token that starts a list to its wrappers
	starter_to_wrapper := map[string]listContext{
		"'":  {[2]string{"(quote ", ")"}, true},
		"~@": {[2]string{"(splice-unquote ", ")"}, true},
		"~":  {[2]string{"(unquote ", ")"}, true},
		"`":  {[2]string{"(quasiquote ", ")"}, true},
		"@":  {[2]string{"(deref ", ")"}, true},
		"^":  {[2]string{"(with-meta ", ")"}, true},

		"(": {[2]string{"(", ")"}, false},
		"{": {[2]string{"{", "}"}, false},
		"[": {[2]string{"[", "]"}, false},
	}

	for starter, listParams := range starter_to_wrapper {
		if first == starter {
			return read_list(r, listParams)
		}
	}

	return read_atom(r)
}

func read_list(r *reader, params listContext) (MalTyper, error) {
	first, ok := r.Next()
	list := new(MalList)
	list.Wraps = params.wrappers

	for first != list.Wraps[1] && ok == nil {
		app, err := read_form(r)
		if err != nil {
			return list, err
		}

		list.Data = append(list.Data, app)
		if params.onlyOne {
			break
		} else {
			first, ok = r.Next()
		}
	}

	return list, ok
}

func read_atom(r *reader) (MalTyper, error) {
	first := r.Peek() // first value always the read atom

	value, is := strconv.ParseInt(first, 10, 0)
	if is == nil {
		return MalInt{int(value)}, nil
	}

	// TODO: MalTrue and MalFalse (they are different)
	if len(first) != 0 && first[0] == '"' {
		if first[len(first)-1] == '"' {
			return MalString{first, r.readably}, nil
		} else {
			return MalString{"", r.readably}, errors.New("unbalanced")
		}
	}

	return MalSymbol{first}, nil
}
