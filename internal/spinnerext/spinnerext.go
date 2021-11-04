package spinnerext

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

var (
	s *spinner.Spinner
)

func Start() {
	s = spinner.New(spinner.CharSets[13], 100*time.Millisecond, spinner.WithWriter(os.Stdout))
	s.Start()
}

func Text(message string) {
	if s == nil {
		return
	}

	s.Suffix = " " + message

	time.Sleep(200 * time.Millisecond)
}

func Pause() {
	if s == nil {
		return
	}

	s.Suffix = ""
	s.Stop()
}

func Unpause() {
	if s == nil {
		return
	}

	s.Restart()
}

func Stop(message string) {
	if s == nil {
		return
	}

	s.Suffix = ""
	s.FinalMSG = message + "\n"
	s.Stop()

	time.Sleep(200 * time.Millisecond)
}
