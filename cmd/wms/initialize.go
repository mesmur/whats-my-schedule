package wms

import (
	"fmt"

	"github.com/MESMUR/wms/pkg/initialize"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "initialize",
	Aliases: []string{"init"},
	Short:   "Initializes What's My Schedule",
	Long:    "Initialize the What's My Schedule connection go Google Calendar using the given credentials",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		clientID := args[0]
		clientSecret := args[1]

		if err := initialize.CheckToken(); err == nil {
			fmt.Printf("Credentials have already been initialized!\n")
			return
		}

		config := initialize.CreateOauth2Config(clientID, clientSecret)
		token := initialize.GetTokenFromWeb(config)
		initialize.SaveToken(token)
		initialize.CreateConfig()

		fmt.Println("Successfully authenticated and stored oAuth Token Credentials")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
