package main

import (
	"errors"
	"os"

	"github.com/akamensky/argparse"
)

type Args struct {
	halSocketPath string
	videosDirPath string
	audiosDirPath string
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("halsimulation", "Simulator for unit client HAL")

	halSocketPath := parser.String("", "hal-socket-path", &argparse.Options{Help: "Path to unix socket file to communicate over with HAL.", Default: "/tmp/unit.hal.sock"})
	videosDirPath := parser.String("", "videos-dir-path", &argparse.Options{Help: "Path to directory of videos to send.", Default: "/tmp/hal.videos"})
	audiosDirPath := parser.String("", "audios-dir-path", &argparse.Options{Help: "Path to directory of audios to send.", Default: "/tmp/hal.audios"})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		halSocketPath: *halSocketPath,
		videosDirPath: *videosDirPath,
		audiosDirPath: *audiosDirPath,
	}, nil
}
