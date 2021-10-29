package main

import (
	"github.com/keygen-sh/cli/cmd"
)

// func main() {
// 	distCmd := flag.NewFlagSet("dist", flag.ExitOnError)

// 	account := distCmd.String("account", "", "Your Keygen account ID")
// 	product := distCmd.String("product", "", "Your Keygen product ID")
// 	token := distCmd.String("token", "", "Your Keygen product token")

// 	entitlements := distCmd.String("entitlements", "", "Attach entitlements to release")
// 	channel := distCmd.String("channel", "", "The release channel")
// 	platform := distCmd.String("platform", "", "The release platform")

// 	if len(os.Args) < 2 {
// 		// TODO(ezekg) Add help
// 		fmt.Println("<help>")

// 		os.Exit(1)
// 	}

// 	switch os.Args[1] {
// 	case "genkey":
// 		// TODO(ezekg)
// 	case "dist":
// 		distCmd.Parse(os.Args[2:])

// 		fmt.Printf("running dist with args <%v>\n", *entitlements)

// 		keygenext.Account = *account
// 		keygenext.Product = *product
// 		keygenext.Token = *token

// 		r := &keygenext.Release{
// 			Name:      "Version 1.0",
// 			Version:   "1.0.0",
// 			Filename:  "1.0.0.dmg",
// 			Filetype:  "dmg",
// 			Platform:  *platform,
// 			Channel:   *channel,
// 			ProductID: *product,
// 			Constraints: keygenext.Constraints{
// 				{EntitlementID: *entitlements},
// 			},
// 		}
// 		fmt.Println(r.Upsert())
// 		fmt.Println(r)
// 	}
// }

func main() {
	cmd.Execute()
}
