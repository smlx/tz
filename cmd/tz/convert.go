package main

import (
	"fmt"
	"time"

	"github.com/smlx/tz/internal/converter"
)

// ConvertCmd represents the convert command.
type ConvertCmd struct {
	Target   string   `kong:"arg,name='TARGET',help='Target location to convert to (e.g. Zürich)'"`
	Source   string   `kong:"arg,optional,name='SOURCE',help='Source location to convert from (e.g. Sydney). @ is shorthand for the local timezone.'"`
	TimeSpec []string `kong:"arg,optional,name='TIME SPECIFICATION',help='Time specification (e.g. 5am friday)'"`
}

// Run executes the convert command.
func (c *ConvertCmd) Run() error {
	res, err := converter.Convert(c.Target, c.Source, c.TimeSpec, time.Now())
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
