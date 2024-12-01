package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"
	"text/template"

	_ "embed"
)

type event struct {
	Name  string
	Value string
	Src   []state
	Dst   state
}

type state struct {
	Name  string
	Value string
}

type data struct {
	PkgName        string
	States         []state
	Events         []event
	PossibleEvents map[state][]event
	GenDynamic     bool
	GenType        bool
}

//go:embed fsm.go.tmpl
var s string

func possibleEvents(states []state, events []event) map[state][]event {
	possible := map[state][]event{}
	for _, s := range states {
		possible[s] = []event{}
	}

	for _, e := range events {
		for _, s := range e.Src {
			possible[s] = append(possible[s], e)
		}
	}

	return possible
}

func Gen(w io.Writer, d data) error {
	funcs := template.FuncMap{
		"ToTitle": strings.Title,
	}
	tm, err := template.New("fsm").Funcs(funcs).Parse(s)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	var buff bytes.Buffer
	err = tm.Execute(&buff, d)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	code, err := io.ReadAll(&buff)
	if err != nil {
		return err
	}

	code, err = format.Source(code)
	if err != nil {
		return fmt.Errorf("formatting code: %w", err)
	}

	_, err = w.Write(code)
	if err != nil {
		return err
	}

	return nil
}
