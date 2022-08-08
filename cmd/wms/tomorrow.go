package wms

import (
	"context"
	"fmt"
	"log"

	"github.com/MESMUR/wms/pkg/events"
	"github.com/MESMUR/wms/pkg/initialize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var defTomorrow bool

var tomorrowCmd = &cobra.Command{
	Use:     "tomorrow",
	Aliases: []string{"tom"},
	Short:   "Gets tomorrow's schedule!",
	Long:    "Gets tomorrow's schedule for the given calendar!",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var calendarID string
		if len(args) > 0 {
			calendarID = args[0]
		} else {
			calendarID = fmt.Sprint(viper.Get("calendar_name"))
		}

		if defTomorrow {
			viper.Set("calendar_name", calendarID)
			viper.WriteConfig()
		}

		ctx := context.Background()

		client := initialize.GetClient()
		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))

		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}

		events.GetEvents(srv, calendarID, 1)
	},
}

func init() {
	tomorrowCmd.Flags().
		BoolVarP(&defTomorrow, "default", "d", false, "Sets the provided calendar as the default")
	rootCmd.AddCommand(tomorrowCmd)
}
