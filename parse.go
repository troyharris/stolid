package main

import (
		"github.com/knieriem/markdown"
		"bufio"
		"os"
		"fmt"
		"path/filepath"
		"strings"
)

func parseFile (infile string, outfile string) {
		p := markdown.NewParser(&markdown.Extensions{Smart: true})

   	fi, err := os.Open(infile)
    if err != nil { panic(err) }
    // close fi on exit and check for its returned error
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()

    r := bufio.NewReader(fi)

   // open output file
    fo, err := os.Create(outfile)
    if err != nil { panic(err) }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }() 

    // make a write buffer
    w := bufio.NewWriter(fo)

    p.Markdown(r, markdown.ToHTML(w))

    if err = w.Flush(); err != nil { panic(err) }
}

func dirWalk (path string, info os.FileInfo, err error) error {
	if info.IsDir() == false && filepath.Ext(path) == ".md"{
		fileroot := strings.Split(info.Name(), ".md")[0]
		htmlfile := filepath.Dir(path) + "/" + fileroot + ".html"
		fmt.Printf("Found %s\n", fileroot)
		parseFile(path, htmlfile)
	}
	return nil
}

func main() {
	parseFile("test.md", "test.html");
	root := "."
	_ = filepath.Walk(root, dirWalk)
}