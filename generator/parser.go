package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
)

type MockParser struct{}

func (mc *MockParser) Parse() data {
	stillState := state{Name: "Still", Value: "still"}
	movingState := state{Name: "Moving", Value: "moving"}
	crashState := state{Name: "Crash", Value: "crash"}
	states := []state{stillState, movingState, crashState}

	events := []event{
		{
			Name:  "moveEvent",
			Value: "Move",
			Src:   []state{stillState, movingState},
			Dst:   movingState,
		},
		{
			Name:  "stopEvent",
			Value: "Stop",
			Src:   []state{movingState},
			Dst:   stillState,
		},
		{
			Name:  "toCrashEvent",
			Value: "ToCrash",
			Src:   []state{movingState},
			Dst:   crashState,
		},
	}

	d := data{
		PkgName:        "generated",
		States:         states,
		Events:         events,
		PossibleEvents: possibleEvents(states, events),
		GenDynamic:     false,
	}

	return d
}

type RealParser struct {
	CommentParser *CommentsParser
}

func (r *RealParser) Parse() (data, error) {
	fd, err := r.CommentParser.Parse()
	if err != nil {
		return data{}, fmt.Errorf("parsing: %w", err)
	}

	return BetterNaming(fd)
}

type docType string

const (
	mermaid  docType = "mermaid"
	graphviz docType = "graphviz"
)

type CommentsParser struct {
	stateTypeName string
	fileName      string
	lineParser    lineParser
	src           io.Reader
}

func NewCommentsParser(src io.Reader, fileName string, typeName string, t string) (*CommentsParser, error) {
	var cp lineParser
	switch docType(t) {
	case mermaid:
		cp = &mermaidParser{}
	case graphviz:
		cp = &graphvizParser{}
	default:
		return nil, fmt.Errorf("unknown doc type: %s", t)
	}

	parser := &CommentsParser{
		stateTypeName: typeName,
		lineParser:    cp,
		fileName:      fileName,
		src:           src,
	}

	return parser, nil
}

type fileData struct {
	packageName string
	td          TransitionData
}

func (r *CommentsParser) Parse() (fileData, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", r.src, parser.ParseComments)
	if err != nil {
		return fileData{}, err
	}

	var td TransitionData

	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {
			case token.TYPE:
				var name string
				for _, spec := range decl.Specs {
					tspec := spec.(*ast.TypeSpec)
					name = tspec.Name.Name
				}
				if name != r.stateTypeName {
					continue
				}

				comments := commentText(decl.Doc)

				td, err = r.extractTransitionData(comments)
				if err != nil {
					return fileData{}, err
				}
				break
			}
		}
	}

	return fileData{
		packageName: f.Name.Name,
		td:          td,
	}, nil
}

func commentText(cg *ast.CommentGroup) []string {
	if cg == nil {
		return nil
	}

	lines := make([]string, 0, len(cg.List))
	for _, line := range cg.List {
		lines = append(lines, line.Text)
	}

	return lines
}

func (r *CommentsParser) extractTransitionData(lines []string) (TransitionData, error) {
	var (
		states = make(map[string]struct{})
		events []lineData
	)
	for _, line := range lines {
		l, err := r.lineParser.ParseLine(line)
		if err != nil {
			return TransitionData{}, err
		}

		if _, ok := states[l.Src]; !ok {
			states[l.Src] = struct{}{}
		}

		if _, ok := states[l.Dst]; !ok {
			states[l.Dst] = struct{}{}
		}

		events = append(events, l)
	}

	var statesSlice = make([]string, 0, len(states))
	for name := range states {
		statesSlice = append(statesSlice, name)
	}

	td := TransitionData{
		Events: events,
	}

	return td, nil
}

type TransitionData struct {
	Events []lineData
}

type lineData struct {
	Src   string
	Dst   string
	Label string
}

type lineParser interface {
	ParseLine(line string) (lineData, error)
}

type mermaidParser struct{}

func (mp *mermaidParser) ParseLine(line string) (lineData, error) {
	line = strings.TrimPrefix(line, "//")
	var label string
	parts := strings.Split(line, ":")

	if len(parts) > 2 {
		return lineData{}, fmt.Errorf("invalid mermaid line: %s", line)
	}

	if len(parts) > 1 {
		label = strings.TrimSpace(parts[1])
	}

	states := strings.Split(parts[0], "-->")
	if len(states) != 2 {
		return lineData{}, fmt.Errorf("invalid mermaid line: %s", line)
	}

	return lineData{
		Src:   strings.TrimSpace(states[0]),
		Dst:   strings.TrimSpace(states[1]),
		Label: label,
	}, nil
}

type graphvizParser struct{}

func (gp *graphvizParser) ParseLine(line string) (lineData, error) {
	return lineData{}, nil
}
