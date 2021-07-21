package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
)

func waitSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func fileExists(path string) bool {
	_, err := os.Open(path)
	return err == nil
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
		if iface != "lo" && iface != "" {
			return iface
		}
	}

	log.Panic("no default interface found")
	return "unreachable"
}
