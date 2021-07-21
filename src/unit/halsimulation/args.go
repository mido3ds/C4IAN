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
	fragmentDurSecs int
	live     bool
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
		live:     *live,
	}, nil
}
