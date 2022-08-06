package today

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/fatih/color"
	"google.golang.org/api/calendar/v3"
)

var green = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
var blue = color.New(color.FgHiBlue, color.Bold).SprintfFunc()
var red = color.New(color.FgHiRed, color.Bold).SprintfFunc()
var magenta = color.New(color.FgHiMagenta, color.Bold, color.Underline).SprintfFunc()

var linkRegex, _ = regexp.Compile(`https://(.*?\.)?zoom.us/[a-z]/[^\"\s]*`)
var docRegex, _ = regexp.Compile(
	`(https://docs\.google\.com/document/d/[^/\s]*)|(https://(.*?\.)?atlassian.net/wiki/spaces/(.*?/)pages/.+?$)`,
)

// GetEvents oops
func GetEvents(srv *calendar.Service, calendarID string) {
	tod := time.Now()
	tom := tod.AddDate(0, 0, 1)
	tMin := time.Date(tod.Year(), tod.Month(), tod.Day(), 0, 0, 0, 0, tod.Location()).
		Format(time.RFC3339)
	tMax := time.Date(tom.Year(), tom.Month(), tom.Day(), 0, 0, 0, 0, tom.Location()).
		Format(time.RFC3339)
	events, err := srv.Events.List("maahir@xendit.co").
		ShowDeleted(false).
		SingleEvents(true).
		OrderBy("startTime").
		TimeMin(tMin).
		TimeMax(tMax).
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve user events: %v", err)
	}
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		fullDayEvents, otherEvents := sortEvents(events)
		printFullEvents(fullDayEvents)
		printEvents(otherEvents)
	}
}

func sortEvents(events *calendar.Events) ([]*calendar.Event, []*calendar.Event) {
	fullDayEvents := []*calendar.Event{}
	otherEvents := []*calendar.Event{}

	for _, event := range events.Items {
		// If the start and end datetime values are not set they are 'day long events'
		// But considering the scope of our setting they are at most 1 day long
		if event.Start.DateTime == "" && event.End.DateTime == "" {
			fullDayEvents = append(fullDayEvents, event)
		} else {
			otherEvents = append(otherEvents, event)
		}
	}
	return fullDayEvents, otherEvents
}

func printFullEvents(events []*calendar.Event) {
	var output string
	for _, event := range events {
		output = fmt.Sprintf(
			"%v: %v \n",
			magenta(event.Start.Date),
			event.Summary,
		)
		fmt.Printf(output)
	}
	fmt.Println()
}

func printEvents(events []*calendar.Event) {
	var output string
	var tStart, tEnd string
	format := "15:04"
	for _, event := range events {
		timeStart, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		timeEnd, _ := time.Parse(time.RFC3339, event.End.DateTime)

		if t := time.Now(); t.After(timeStart) && t.Before(timeEnd) {
			tStart = blue(timeStart.Format(format))
			tEnd = blue(timeEnd.Format(format))
		} else if t.After(timeStart) {
			tStart = red(timeStart.Format(format))
			tEnd = red(timeEnd.Format(format))
		} else {
			tStart = green(timeStart.Format(format))
			tEnd = green(timeEnd.Format(format))
		}

		output = fmt.Sprintf(
			"%v - %v: %v",
			tStart,
			tEnd,
			event.Summary,
		)
		if link := getEventLink(event); link != "" {
			output += fmt.Sprintf("\n\t- Link: %v ", link)
		}
		if docLink := getDocumentLink(event); docLink != "" {
			output += fmt.Sprintf("\n\t- Docs: %v ", docLink)
		}
		output += fmt.Sprint("\n")
		fmt.Printf(output)
	}
}

func getEventLink(event *calendar.Event) string {
	link := event.HangoutLink
	if link == "" {
		link = linkRegex.FindString(event.Description)
	}
	return link
}

func getDocumentLink(event *calendar.Event) string {
	link := docRegex.FindString(event.Description)
	return link
}
