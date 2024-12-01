package generator

import (
	"fmt"
	"os"
)

type Parser interface {
	Parse() (data, error)
}

type FSMGenerator struct {
	parser Parser
}

func NewFSMGenerator(parser Parser) *FSMGenerator {
	return &FSMGenerator{parser: parser}
}

func (g *FSMGenerator) Do() error {
	data, err := g.parser.Parse()
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	file, err := os.Create("fsm_gen.go")
	if err != nil {
		return err
	}

	err = Gen(file, data)
	if err != nil {
		return fmt.Errorf("gen: %w", err)
	}

	return nil
}

type InjectedFlags struct {
	PkgName    string
	GenType    bool
	GenDynamic bool
}

func InjectFlags(d data, flags InjectedFlags) data {
	if flags.PkgName != "" {
		d.PkgName = flags.PkgName
	}

	d.GenType = flags.GenType
	d.GenDynamic = flags.GenDynamic

	return d
}
