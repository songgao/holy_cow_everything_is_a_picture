package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var acceptedImageExt = []string{".png", ".gif", ".jpg", ".jpeg"}

func isAcceptedFile(path string) bool {
	ext := filepath.Ext(path)
	for _, v := range acceptedImageExt {
		if ext == v {
			return true
		}
	}
	return false
}

type contentManager struct {
	filesDirPath            string
	structureDefinitionPath string

	structureJSON   string
	muStructureJSON *sync.RWMutex

	fileServer http.Handler
}

func newContentManager(filesDirPath, structureDefinitionPath string) (*contentManager, error) {
	var err error
	i := &contentManager{filesDirPath: filesDirPath, structureDefinitionPath: structureDefinitionPath}
	i.fileServer = http.FileServer(http.Dir(filesDirPath))
	i.structureJSON, err = i.parseStructure()
	i.muStructureJSON = new(sync.RWMutex)
	if err != nil {
		return nil, err
	}
	i.startWatching()
	return i, nil
}

func (i *contentManager) parseStructure() (string, error) {
	file, err := os.Open(i.structureDefinitionPath)
	if err != nil {
		return "", err
	}
	parsed, err := parseStructure(file)
	if err != nil {
		return "", err
	}
	j, err := json.Marshal(parsed)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (i *contentManager) startWatching() error {
	structureStat, err := os.Stat(i.structureDefinitionPath)
	if err != nil {
		return err
	}
	structureLastMod := structureStat.ModTime().Unix()
	go func() {
		for {
			time.Sleep(4 * time.Second)
			structureStat, err := os.Stat(i.structureDefinitionPath)
			if err == nil {
				mod := structureStat.ModTime().Unix()
				if structureLastMod != mod {
					j, err := i.parseStructure()
					if err == nil {
						i.muStructureJSON.Lock()
						i.structureJSON = j
						i.muStructureJSON.Unlock()
						structureLastMod = mod
					}
				}
			}
		}
	}()
	return nil
}

func (i *contentManager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, ext := range acceptedImageExt {
		if strings.HasSuffix(req.RequestURI, ext) {
			i.fileServer.ServeHTTP(w, req)
			return
		}
	}
}
