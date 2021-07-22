package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"strings"
)

func uniform(min, max float64) float64 {
	return min + (rand.Float64() * (max - min))
}

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

func getDefaultInterface() string {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		const line = 1 // line containing the gateway addr. (first line: 0)
		// jump to line containing the agteway address
		for i := 0; i < line; i++ {
			scanner.Scan()
		}

		// get field containing gateway address
		tokens := strings.Split(scanner.Text(), "\t")
		iface := tokens[0]
		return iface
	}

	log.Panic("no default interface found")
	return "unreachable"
}
