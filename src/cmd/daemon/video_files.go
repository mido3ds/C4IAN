package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mido3ds/C4IAN/src/models"
)

const metadataFileName = "index.m3u8"

type VideoFilesManager struct {
	path string
}

func NewVideoFilesManager(path string) *VideoFilesManager {
	// Delete the directory recursively if it exists
	err := os.RemoveAll(path)
	if err != nil {
		log.Panic(err)
	}

	// Create a new directory
	err = os.Mkdir(path, 0755)
	if err != nil {
		log.Panic(err)
	}

	return &VideoFilesManager{path: path}
}

func (v *VideoFilesManager) AddFragment(frag *models.VideoFragment) {
	// Create directory for this streams if it does not exist
	dir := filepath.Join(v.path, frag.Src, strconv.Itoa(frag.ID))
	exists := pathExists(dir)
	if !exists {
		os.MkdirAll(dir, 0755)
	}

	// Write metadata
	path := filepath.Join(dir, metadataFileName)
	writeFile(path, frag.Metadata)

	// Write fragment
	path = filepath.Join(dir, frag.FileName)
	writeFile(path, frag.Body)
}

func writeFile(path string, data []byte) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	n, err := file.Write(data)
	if n != len(data) {
		log.Panic("Could not write the whole file")
	}
	if err != nil {
		log.Panic(err)
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Panic(err)
	return false
}
