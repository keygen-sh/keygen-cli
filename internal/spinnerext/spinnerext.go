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
	s.Suffix = " " + message

	time.Sleep(200 * time.Millisecond)
}

func Pause() {
	s.Suffix = ""
	s.Stop()
}

func Unpause() {
	s.Restart()
}

func Stop(message string) {
	s.Suffix = ""
	s.FinalMSG = message + "\n"
	s.Stop()

	time.Sleep(200 * time.Millisecond)
}
