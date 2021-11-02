package cmd

type Flags struct {
	filename    string
	name        string
	version     string
	platform    string
	channel     string
	constraints []string
	signingKey  string
}

var (
	flags = &Flags{}
)
