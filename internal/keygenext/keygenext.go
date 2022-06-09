// keygenext is an internal package that allows us to create a separate
// shadow instance of keygen that can be used for the end-user, since
// the main package is used for CLI auto-upgrades.
package keygenext

var (
	Account   string
	Product   string
	Token     string
	PublicKey string
	UserAgent string
)
