package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

var (
	fContentPath             string
	fStructureDefinitionPath string
	fLAddr                   string
)

func init() {
	flag.StringVar(&fContentPath, "content", "", "path to the folder that has images (supported formats: .jpg, .jpeg, .png, .gif)")
	flag.StringVar(&fStructureDefinitionPath, "structure", "", "path to the json file that defines structure of images")
	flag.StringVar(&fLAddr, "laddr", "localhost:7428", "http listening address")
}

func checkFlags() (flagsOK bool) {
	flag.Parse()
	if fContentPath == "" {
		return false
	}
	if fStructureDefinitionPath == "" {
		return false
	}

	return true
}

func main() {
	if !checkFlags() {
		flag.PrintDefaults()
		return
	}
	cm, err := newContentManager(fContentPath, fStructureDefinitionPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	mux, err := initMux(cm)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = http.ListenAndServe(fLAddr, mux)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func initMux(cm *contentManager) (http.Handler, error) {
	mux := http.NewServeMux()
	assetsPath, err := getAssetsPath()
	if err != nil {
		return nil, err
	}
	mux.Handle("/", http.FileServer(http.Dir(assetsPath)))
	mux.Handle("/content/", http.StripPrefix("/content", cm))
	mux.HandleFunc("/content.json", func(w http.ResponseWriter, req *http.Request) {
		cm.muStructureJSON.RLock()
		io.WriteString(w, cm.structureJSON)
		cm.muStructureJSON.RUnlock()
	})
	return mux, nil
}
