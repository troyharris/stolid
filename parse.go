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
    "bytes"
    "text/template"
    "sort"
)

type Options struct {
    SiteName, WebRoot, ContentPath, DestPath, TemplatePath string
}

type Article struct {
    Body string
}

type Header struct {
    SiteName, Title, Menu string
}

type PageFile struct {
    MarkdownPath, AbsFilePath, RelFilePath, Title string
}

type PageContent struct {
    Content string
}

var config = Options{}
var siteMap = make(map[string][]PageFile)

func readConfig () {
    configBytes := readFile("config.json")
    _ = json.Unmarshal(configBytes, &config)
}

func readFile (file string) []byte {
    b, err := ioutil.ReadFile(file)
    if err != nil { panic(err) }
    return b
}

func parseFile (infile string) string {
    p := markdown.NewParser(&markdown.Extensions{Smart: true})
    fi, err := os.Open(infile)
    if err != nil { panic(err) }
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()
    r := bufio.NewReader(fi)
    bw := make([]byte, 0)
    w := bytes.NewBuffer(bw)
    p.Markdown(r, markdown.ToHTML(w))
    return string(w.Bytes())
}

func parseArticle (infile string) Article {
    content := parseFile(infile)
    return Article{Body: content}
}

func parseHeader (title string) Header {
    return Header{SiteName: config.SiteName, Title: title, Menu: buildMenu()}
}

func parseTemplate (file string, data interface{}) (out []byte, error error) {
    var buf bytes.Buffer
    t, err := template.ParseFiles(file)
    if err != nil {
        return nil, err
    }
    err = t.Execute(&buf, data)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func createArticlePage (infile string, title string) []byte {
    a := parseArticle(infile)
    h := parseHeader(title)
    compiledHead, _ := parseTemplate(config.TemplatePath + "head.html", h)
    compiledArticle, _ := parseTemplate(config.TemplatePath + "article.html", a)
    content := PageContent{Content: string(compiledArticle)}
    compiledHTML, _ := parseTemplate(config.TemplatePath + "index.html", content)
    return append(compiledHead, compiledHTML...)
}

func createCategoryPage (files []PageFile, name string) []byte {
    title := cleanTitle(name)
    h := parseHeader(title)
    compiledHead, _ := parseTemplate(config.TemplatePath + "head.html", h)
    articles := make([]byte, 0)
    for _, p := range files {
        a := parseArticle(p.MarkdownPath)
        b, _ := parseTemplate(config.TemplatePath + "article.html", a)
        articles = append(articles, b...)
    }
    content := PageContent{Content: string(articles)}
    compiledHTML, _ := parseTemplate(config.TemplatePath + "index.html", content)
    return append(compiledHead, compiledHTML...)
}

func createIndexPage (files []PageFile) []byte {
    return createCategoryPage(files, config.SiteName)
}

func dirWalk (path string, info os.FileInfo, err error) error {
    if info.IsDir() == false && filepath.Ext(path) == ".md"{
        fileroot := strings.Split(info.Name(), ".md")[0]
        outtail := strings.Split(path, config.ContentPath)[1]
        //fmt.Println(outtail)
        outdir := filepath.Dir(outtail)
        //fmt.Println(outdir)
        relfile := outdir + "/" + fileroot + "/index.html"
        htmlfile := config.DestPath + relfile
        rootwebfile := "/" + relfile
        //fmt.Println(htmlfile)
        title := cleanTitle(fileroot)
        fullHTML := createArticlePage(path, title)
        writeHTML(htmlfile, fullHTML)
        page := PageFile{MarkdownPath: path, AbsFilePath: htmlfile, RelFilePath: rootwebfile, Title: title}
        a, ok := siteMap[outdir]
        if (ok) {
            siteMap[outdir] = append(a, page)
        }  else {
            siteMap[outdir] = []PageFile{page}
        }
    }
    return nil
}

func cleanTitle (filename string) string {
    return strings.Replace(filename, "_", " ", -1)
}

func compileFiles (root string) {
    _ = filepath.Walk(root, dirWalk)
    defer func () {
        buildIndexes()
    }()
}

func buildIndexes () {
    index := make([]PageFile, 0)
    for cat, page := range siteMap {
        index = append(index, page...)
        catPage := createCategoryPage(page, cat)
        catPath := config.DestPath + cat + "/index.html"
        writeHTML(catPath, catPage)
    }
    writeHTML(config.DestPath + "/index.html", createIndexPage(index))
    siteMap = make(map[string][]PageFile)
}

func writeHTML (filePath string, content []byte) {
    newPath := path.Dir(filePath)
    _ = os.MkdirAll(newPath, 0774)
    ioutil.WriteFile(filePath, content, 0774)
}

func getCategories (rootPath string) []string {
    d, _ := ioutil.ReadDir(rootPath)
    a := make([]string, 0)
    for _, v := range d {
        if v.IsDir() {
            a = append(a, v.Name())
        }
    }
    sort.Strings(a)
    return a
}

func buildMenu () string {
    cats := getCategories(config.ContentPath)
    menu := "<ul>"
    for _, v := range cats {
        menu += "<li><a href='" + config.WebRoot + v + "'>" + v + "</a></li>"
    }
    menu += "</ul>"
    return menu
}

func buildSite () {
    //fmt.Printf("DestPath is %s and TemplatePath is %s", config.DestPath, config.TemplatePath)
    root := config.ContentPath
    //compileFiles(root)
    compileFiles(root)
    fmt.Println("Done")
}