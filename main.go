package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dyleme/fsm/generator"
)

var (
	isHelp      = flag.Bool("help", false, "Display help message")
	typeName    = flag.String("typeName", "", "Type name")
	parsedFile  = flag.String("parsedFile", "", "File name")
	packageName = flag.String("package", "", "Package name")
	docType     = flag.String("docType", "mermaid", "Doc type")
	genFileName = flag.String("genFile", "", "Gen file name")
	genType     = flag.Bool("genType", false, "Gen type")
	genDynamic  = flag.Bool("genDynamic", false, "Gen dynamic")
)

func validateFlags() error {
	if *isHelp {
		return fmt.Errorf("Usage: myapp [command] [arguments]")
	}

	if *typeName == "" {
		return fmt.Errorf("typeName is required")
	}

	if *parsedFile == "" {
		return fmt.Errorf("parsedFile is required")
	}

	if *genFileName == "" {
		return fmt.Errorf("genFile name is required")
	}

	return nil
}

func main() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	flag.Parse()
	if err := validateFlags(); err != nil {
		log.Println(err)

		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)

		return
	}
	parts := strings.Split(*parsedFile, "/")
	filename := parts[len(parts)-1]
	pathToFile := wd + "/" + *parsedFile
	file, err := os.Open(pathToFile)
	if err != nil {
		log.Println(err)

		return
	}

	commentsParser, err := generator.NewCommentsParser(file, filename, *typeName, *docType)
	if err != nil {
		log.Println(err)

		return
	}

	rParser := &generator.RealParser{
		CommentParser: commentsParser,
	}

	data, err := rParser.Parse()
	if err != nil {
		log.Println(err)

		return
	}

	data = generator.InjectFlags(data, generator.InjectedFlags{
		PkgName:    *packageName,
		GenType:    *genType,
		GenDynamic: *genDynamic,
	})

	f, err := os.Create(*genFileName)
	if err != nil {
		log.Println(err)

		return
	}

	err = generator.Gen(f, data)
	if err != nil {
		log.Println(err)

		return
	}
}
