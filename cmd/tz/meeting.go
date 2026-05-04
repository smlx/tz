package main

import (
	"fmt"
	"time"

	"github.com/smlx/tz/internal/meeting"
)

// MeetingCmd represents the meeting command.
type MeetingCmd struct {
	Locations []string `kong:"arg,name='LOCATIONS',help='Locations to plan meeting for'"`
}

// Run executes the meeting command.
func (c *MeetingCmd) Run() error {
	res, err := meeting.Plan(c.Locations, time.Now())
	if err != nil {
		return err
	}
	fmt.Print(res)
	return nil
}
