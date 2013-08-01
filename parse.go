package main

import (
		"github.com/knieriem/markdown"
		"bufio"
		"os"
		"fmt"
		"path/filepath"
		"strings"
        "io/ioutil"
)

func parseFile (name string, infile string, outfile string) {
    top, bottom := buildPage(name)
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

    w.WriteString(top)
    p.Markdown(r, markdown.ToHTML(w))
    w.WriteString(bottom)

    if err = w.Flush(); err != nil { panic(err) }
}
/*

func markdownToHTML (infile string) string {
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
    return "hi"
}
*/

func dirWalk (path string, info os.FileInfo, err error) error {
	if info.IsDir() == false && filepath.Ext(path) == ".md"{
		fileroot := strings.Split(info.Name(), ".md")[0]
		htmlfile := filepath.Dir(path) + "/" + fileroot + ".html"
		fmt.Printf("Found %s\n", fileroot)
		parseFile("index", path, htmlfile)
	}
	return nil
}

func templateSection (sectionName string) string {
    h, err := ioutil.ReadFile("templates/" + sectionName + ".html")
    if err != nil { panic(err) }
    return string(h)
}

func buildPage (name string) (string, string) {
    head := templateSection("head")
    top := templateSection(name + "-top")
    bottom := templateSection(name + "-bottom")
    return head + top, bottom
}


func main() {
	root := "content"
	_ = filepath.Walk(root, dirWalk)
}