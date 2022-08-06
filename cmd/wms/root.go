package wms

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "wms",
	Short: "wms - what's my schedule retrieves and shows you your calendar from google calendar",
	Long:  `What's My Schedule is a lightweight tool to retrieve and showcase your schedule from Google Calendar in a pretty manner!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No option specified")
		}
		return nil
	},
}

// Execute runs the CLI!
func Execute() {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	// default cmd if no cmd is given
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != flag.ErrHelp {
		args := append([]string{todayCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Gets the users home directory
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Sets the config path and path of the .wms files
	configPath := home + "/.wms/config.json"
	path := home + "/.wms"

	// Creates the path if it does not exist
	err = os.MkdirAll(path, 0755)
	cobra.CheckErr(err)

	// Set the Viper defaults
	viper.SetDefault("calendar_name", "primary")

	// Save the Config if it does not exist
	viper.SafeWriteConfigAs(configPath)

	// Set the config file and read it in
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	cobra.CheckErr(err)
}
