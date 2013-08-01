package main

import (
		"github.com/knieriem/markdown"
		"bufio"
		"os"
		"fmt"
		"path/filepath"
		"strings"
        "io/ioutil"
        "encoding/json"
        "path"
)

type Options struct {
    ContentPath, DestPath, TemplatePath string
}

var config = Options{}

func readConfig () {
    configBytes := readFile("config.json")
    _ = json.Unmarshal(configBytes, &config)
    //fmt.Println(string(configBytes))
}

func readFile (file string) []byte {
    b, err := ioutil.ReadFile(file)
    if err != nil { panic(err) }
    return b
}

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
    newPath := path.Dir(outfile)
    _ = os.MkdirAll(newPath, 0774)

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

func dirWalk (path string, info os.FileInfo, err error) error {
	if info.IsDir() == false && filepath.Ext(path) == ".md"{
		fileroot := strings.Split(info.Name(), ".md")[0]
        outtail := strings.Split(path, config.ContentPath)[1]
        fmt.Println(outtail)
        outdir := filepath.Dir(outtail)
      //  fmt.Println(path)
		htmlfile := config.DestPath + outdir + "/" + fileroot + "/index.html"
	//	fmt.Printf("Found %s\n", fileroot)
        fmt.Println(htmlfile)
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
    readConfig()
    fmt.Printf("DestPath is %s and TemplatePath is %s", config.DestPath, config.TemplatePath)
	root := config.ContentPath
	_ = filepath.Walk(root, dirWalk)
}