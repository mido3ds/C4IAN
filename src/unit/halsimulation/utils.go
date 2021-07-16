package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
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

func listDir(path string) []string {
	log.Println(path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}

	list := make([]string, 0)
	for _, f := range files {
		path := buildPath(path, f.Name())
		if f.IsDir() {
			list = append(list, listDir(path)...)
		} else {
			list = append(list, path)
		}
	}

	return list
}
