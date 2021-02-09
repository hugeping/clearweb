package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	output            = flag.String("o", "", "Output path. Otherwise uses stdout")
	input             = flag.String("i", "", "Input path. Otherwise uses stdin")
	citationStart     = flag.Int("c", 1, "Start citations from this index")
	citationMarkers   = flag.Bool("m", false, "Use footnote style citation markers")
	numberedLinks     = flag.Bool("n", false, "Number the links")
	emitImagesAsLinks = flag.Bool("e", false, "Emit links to included images")
	linkEmitFrequency = flag.Int("l", 2, "Emit gathered links through the document after this number of paragraphs")
)

func check(e error) {
	if e != nil {
		fmt.Printf("%v\n", e)
		os.Exit(1)
	}
}

func saveFile(contents []byte, path string) {
	d1 := contents
	err := ioutil.WriteFile(path, d1, 0644)
	check(err)
}

func readStdin() string {
	// based on https://flaviocopes.com/go-shell-pipes/
	reader := bufio.NewReader(os.Stdin) //default size is 4096 apparently
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	return string(output)
}

func getInput() (string, error) {
	var inputHtml string

	_, err := os.Stdin.Stat()
	check(err)

	if *input != "" {
		//get the input file from the command line
		dat, err := ioutil.ReadFile(*input)
		check(err)
		inputHtml = string(dat)
	} else {
		// we have a pipe input
		inputHtml = readStdin()
	}
	return inputHtml, nil
}

func main() {
	var inputHtml string

	flag.Parse()

	//get the input from commandline or stdin
	inputHtml, err := getInput()
	check(err)

	//convert html to gmi
	options := NewOptions()
	options.CitationStart = *citationStart
	options.LinkEmitFrequency = *linkEmitFrequency
	options.CitationMarkers = *citationMarkers
	options.NumberedLinks = *numberedLinks
	options.EmitImagesAsLinks = *emitImagesAsLinks

	ctx := NewTraverseContext(*options)

	text, err := FromString(inputHtml, *ctx)

	check(err)

	//process the output
	if *output == "" {
		fmt.Print(text + "\n") //terminate with a new line
	} else {
		//save to the specified output
		gmiBytes := []byte(text + "\n")
		saveFile(gmiBytes, *output)
	}
}
