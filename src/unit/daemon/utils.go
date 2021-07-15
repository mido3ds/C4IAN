package main

import (
	"os"
	"os/signal"
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
