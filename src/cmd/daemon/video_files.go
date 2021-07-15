package main

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/mido3ds/C4IAN/src/models"
)

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

func (v *VideoFilesManager) CreateVideoFile(src string, id int) string {
	// Create directory for unit streams if it does not exist
	dirPath := path.Join(v.path, src)
	exists, err := pathExists(dirPath)
	if err != nil {
		log.Panic(err)
	}
	if !exists {
		os.Mkdir(dirPath, 0755)
	}

	// Create file for this stream
	filePath := path.Join(dirPath, strconv.Itoa(id))
	exists, err = pathExists(filePath)
	if err != nil {
		log.Panic(err)
	}
	if !exists {
		os.Create(filePath)
	} else {
		log.Panic("Stream file already exists")
	}
	return filePath
}

func (v *VideoFilesManager) AppendVideoFragment(frag *models.VideoFragment) {
	// Check that the stream file exists
	filePath := path.Join(v.path, frag.Src, strconv.Itoa(frag.ID))
	exists, err := pathExists(filePath)
	if err != nil {
		log.Panic(err)
	}
	if !exists {
		log.Panic("Stream file does not exist")
	}

	// Open file to append fragment
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// Append fragment
	_, err = file.Write(frag.Body)
	if err != nil {
		log.Panic(err)
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
