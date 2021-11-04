package spinnerext

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

type Spinner struct {
	s *spinner.Spinner
}

func (s *Spinner) Start() {
	s.s.Start()
}

func (s *Spinner) Update(message string) {
	s.s.Suffix = " " + message

	time.Sleep(200 * time.Millisecond)
}

func (s *Spinner) Stop(message string) {
	s.s.Suffix = ""
	s.s.FinalMSG = message + "\n"
	s.s.Stop()

	time.Sleep(200 * time.Millisecond)
}

func New() *Spinner {
	s := spinner.New(spinner.CharSets[13], 100*time.Millisecond, spinner.WithWriter(os.Stdout))

	return &Spinner{s}
}
