package main

import (
    "fmt"
    "net/http"
    "os"
    "io/ioutil"
    "strings"
    "path/filepath"
    "time"
)

var rootPath string

func loadPage (file string) (string, error) {
	filePath := strings.TrimRight(rootPath, "/") + file
	var fileBytes []byte
	var err error
	info, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		fileBytes, err = ioutil.ReadFile(filePath + "index.html")
	} else {
		fileBytes, err = ioutil.ReadFile(filePath)
	}
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}

func handler (w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimRight(rootPath, "/") + r.URL.Path
	var err error
	info, err := os.Stat(filePath)
	if info.IsDir() {
		http.ServeFile(w, r, filePath + "index.html")
	} else {
		http.ServeFile(w, r, filePath)
	}
	if err != nil {
		fmt.Fprintf(w, "There was an error: %s", err)
		return
	}
}


func staticHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Println("YO")
	http.ServeFile(w, r, r.URL.Path)
}

func walkHandlers (path string, info os.FileInfo, err error) error {
	root := "/" + strings.Split(path, rootPath)[1]
	if info.IsDir() {
		http.HandleFunc(root, handler)
	}
	return nil
}

func updateExists() (bool, error) {
    _, err := os.Stat("update")
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func checkUpdate () error {
	u, err := updateExists()
	if (u) {
		os.Remove("update")
		buildSite()
	}
	return err
}

func startUpdateLoop () {
	for {
		_ = checkUpdate()
		time.Sleep(5 * time.Second)
	}
}

func main() {
	readConfig()
	rootPath = config.DestPath
	go startUpdateLoop()
	_ = filepath.Walk(rootPath, walkHandlers)
 //   http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}