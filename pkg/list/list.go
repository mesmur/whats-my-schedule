package list

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"google.golang.org/api/calendar/v3"
)

var green = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
var blue = color.New(color.FgHiBlue, color.Bold).SprintfFunc()
var red = color.New(color.FgHiRed, color.Bold).SprintfFunc()

// GetCalendars prints a list of calendars
func GetCalendars(srv *calendar.Service) {
	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve user calendars: %v", err)
	}
	for _, calendar := range calendars.Items {
		fmt.Printf("%v \n", colorCalendar(calendar))
	}
}

func colorCalendar(calendar *calendar.CalendarListEntry) string {
	if calendar.Primary == true {
		return red(calendar.Id)
	}
	return blue(calendar.Id)
}
