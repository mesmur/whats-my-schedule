package wms

import (
	"context"
	"log"

	"github.com/MESMUR/wms/pkg/initialize"
	"github.com/MESMUR/wms/pkg/list"
	"github.com/spf13/cobra"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Gets a list of Calendars by ID!",
	Long:  "Gets a list of Calendars by ID that can be used by the 'schedule' commands to choose from which Calendar you would like your events displayed!",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		client := initialize.GetClient()
		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))

		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}

		list.GetCalendars(srv)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
