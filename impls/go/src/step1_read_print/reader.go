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
	return r.position < len(r.tokens)-1 // last token always an empty string
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

func recursive_read(r *reader, listData []MalTyper) ([]MalTyper, error) {
	app, err := read_form(r)

	if err == nil {
		return append(listData, app), nil
	} else {
		return listData, err
	}
}

type listType interface {
	Read(r *reader) (MalTyper, error)
}

type normal struct {
	Wraps [2]string
}

func (nl normal) Read(r *reader) (MalTyper, error) {
	first, err := r.Next()
	list := new(MalList)
	list.Wraps = nl.Wraps

	for first != list.Wraps[1] && err == nil {
		list.Data, err = recursive_read(r, list.Data)
		if err != nil {
			return list, err
		}

		first, err = r.Next()
	}

	return list, err
}

type special struct {
	Wraps [2]string
}

func (sl special) Read(r *reader) (MalTyper, error) {
	_, err := r.Next()
	list := new(MalList)
	list.Wraps = sl.Wraps

	list.Data, err = recursive_read(r, list.Data)
	if err != nil {
		return list, err
	}

	return list, err
}

type meta struct {
	Wraps [2]string
}

func (ml meta) Read(r *reader) (MalTyper, error) {
	_, err := r.Next()
	list := new(MalList)
	list.Wraps = ml.Wraps

	last_app, err := read_form(r)
	if err != nil {
		return list, err // only should throw error if there isnt a value after the meta call or
	}

	_, err = r.Next()
	for err == nil {
		list.Data, err = recursive_read(r, list.Data)
		if err != nil {
			return list, err // the recursive call fails
		}

		_, err = r.Next()
	}

	list.Data = append(list.Data, last_app)
	return list, nil
}

func read_form(r *reader) (MalTyper, error) {
	first := r.Peek()

	// this map directs the token that starts a list to its wrappers
	starter_to_wrapper := map[string]listType{
		"'":  special{[2]string{"(quote ", ")"}},
		"~@": special{[2]string{"(splice-unquote ", ")"}},
		"~":  special{[2]string{"(unquote ", ")"}},
		"`":  special{[2]string{"(quasiquote ", ")"}},
		"@":  special{[2]string{"(deref ", ")"}},

		"^": meta{[2]string{"(with-meta ", ")"}},

		"(": normal{[2]string{"(", ")"}},
		"{": normal{[2]string{"{", "}"}},
		"[": normal{[2]string{"[", "]"}},
	}

	for starter, list := range starter_to_wrapper {
		if first == starter {
			return list.Read(r)
		}
	}

	return read_atom(r)
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
