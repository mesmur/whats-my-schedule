package wms

import (
	"context"
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
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		client := initialize.GetClient()
		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))

		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}

		today.GetEvents(srv)
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
