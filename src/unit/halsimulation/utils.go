package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path"
)

func normal(mean, stdDev float64) float64 {
	return rand.NormFloat64()*stdDev + mean
}

func waitSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func buildPath(a, b string) string {
	return fmt.Sprintf("%s/%s", a, b)
}

func listDir(pathToDir string) []string {
	log.Println(pathToDir)
	files, err := ioutil.ReadDir(pathToDir)
	if err != nil {
		log.Panic(err)
	}

	list := make([]string, 0)
	for _, f := range files {
		p := path.Join(pathToDir, f.Name())
		if f.IsDir() {
			list = append(list, listDir(p)...)
		} else {
			list = append(list, p)
		}
	}

	return list
}
