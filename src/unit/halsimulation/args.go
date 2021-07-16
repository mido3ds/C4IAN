package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type Args struct {
	halSocketPath     string
	videosDirPath     string
	audiosDirPath     string
	fragmentSizeBytes int
}

func (a *Args) String() string {
	return fmt.Sprintf("&Args{halSocketPath: %v, videosDirPath: %v, audiosDirPath: %v, fragmentSizeBytes: %v}", a.halSocketPath, a.videosDirPath, a.audiosDirPath, a.fragmentSizeBytes)
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("halsimulation", "Simulator for unit client HAL")

	halSocketPath := parser.String("", "hal-socket-path", &argparse.Options{Help: "Path to unix socket file to communicate over with HAL.", Default: "/tmp/unit.hal.sock"})
	videosDirPath := parser.String("", "videos-dir-path", &argparse.Options{Help: "Path to directory of videos to send.", Default: "/tmp/hal.videos"})
	audiosDirPath := parser.String("", "audios-dir-path", &argparse.Options{Help: "Path to directory of audios to send.", Default: "/tmp/hal.audios"})
	fragmentSizeBytes := parser.Int("", "fragment-size-bytes", &argparse.Options{Help: "Size in bytes of one video fragment.", Default: 512})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		halSocketPath:     *halSocketPath,
		videosDirPath:     *videosDirPath,
		audiosDirPath:     *audiosDirPath,
		fragmentSizeBytes: *fragmentSizeBytes,
	}, nil
}
