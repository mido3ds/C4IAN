package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type Args struct {
	halSocketPath   string
	videoPath       string
	audiosDirPath   string
	ffmpegPath      string
	locationSocket  string
	iface           string
	fragmentDurSecs int
	live            bool
	sensors         bool
}

func (a *Args) String() string {
	return fmt.Sprintf("&Args{halSocketPath: %v, videosDirPath: %v, audiosDirPath: %v, fragmentDurSecs: %v}", a.halSocketPath, a.videoPath, a.audiosDirPath, a.fragmentDurSecs)
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("halsimulation", "Simulator for unit client HAL")

	halSocketPath := parser.String("", "hal-socket-path", &argparse.Options{Help: "Path to unix socket file to communicate over with HAL.", Default: "/tmp/unit.hal.sock"})
	video := parser.String("v", "video", &argparse.Options{Help: "Path to video to send, maybe ommited if no video to send.", Default: ""})
	ffmpegPath := parser.String("", "ffmpeg-path", &argparse.Options{Help: "Path to ffmpeg executable.", Default: "ffmpeg"})
	audiosDirPath := parser.String("", "audios-dir", &argparse.Options{Help: "Path to directory of audios to send. If not provided won't send audio messages.", Default: ""})
	fragmentDurSecs := parser.Int("", "fragment-dur-secs", &argparse.Options{Help: "Duration in secs of one video fragment.", Default: 2})
	live := parser.Flag("", "live", &argparse.Options{Help: "Indicates this is streaming from camera not from file, false by default."})
	noSensors := parser.Flag("", "no-sensors", &argparse.Options{Help: "Turn off sensors data generation, false by default."})
	locationSocket := parser.String("", "location-socket", &argparse.Options{Help: "Path to unixgrams to receive location updates on, works the same as router, if not provided will create fake locations.", Default: ""})

	iface := parser.String("", "iface", &argparse.Options{Help: "Name of this interface. Default is to list the ifaces with /proc/net/route.", Default: ""})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		halSocketPath:   *halSocketPath,
		videoPath:       *video,
		audiosDirPath:   *audiosDirPath,
		fragmentDurSecs: *fragmentDurSecs,
		ffmpegPath:      *ffmpegPath,
		live:            *live,
		sensors:         !(*noSensors),
		locationSocket:  *locationSocket,
		iface:           *iface,
	}, nil
}
