package main

import (
		"github.com/knieriem/markdown"
		"bufio"
		"os"
)

func main() {
	p := markdown.NewParser(&markdown.Extensions{Smart: true})

   	fi, err := os.Open("test.md")
    if err != nil { panic(err) }
    // close fi on exit and check for its returned error
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()

    r := bufio.NewReader(fi)

   // open output file
    fo, err := os.Create("test.html")
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

	//	w := bufio.NewWriter(os.Stdout)
	//	p.Markdown(os.Stdin, markdown.ToHTML(w))
	//	w.Flush()
	}