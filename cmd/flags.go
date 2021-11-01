package cmd

type Flags struct {
	name        string
	version     string
	platform    string
	channel     string
	signature   string
	checksum    string
	constraints []string
}

var (
	flags = &Flags{}
)
