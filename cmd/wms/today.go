package wms

import (
	"context"
	"fmt"
	"log"

	"github.com/MESMUR/wms/pkg/initialize"
	"github.com/MESMUR/wms/pkg/today"
	"github.com/spf13/cobra"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var todayCmd = &cobra.Command{
	Use:     "today",
	Aliases: []string{"tod"},
	Short:   "Gets today's schedule!",
	Long:    "Gets today's schedule for the given calendar!",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var calendarID string
		if len(args) > 0 {
			calendarID = args[0]
		} else {
			calendarID = "primary"
		}

		fmt.Printf("calID: %v", calendarID)

		ctx := context.Background()

		client := initialize.GetClient()
		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))

		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}

		today.GetEvents(srv, calendarID)
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
