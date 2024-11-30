package generator

import "os"

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
		return err
	}

	file, err := os.Create("fsm_gen.go")
	if err != nil {
		return err
	}

	return Gen(file, data)
}
