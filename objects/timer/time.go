package timer

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var fakeTime *int64 = nil
var localTimezone *time.Location = time.Local

// Maps formatting tokens to Go time layout equivalents, right now
// it supports a limited set of tokens from PHP's date function
// See: https://www.php.net/manual/en/datetime.format.php
var layoutReplacements = map[string]string{
	// Days
	"%d": "02", "%D": "Mon",
	// Months
	"%m": "01", "%M": "Jan", "%F": "January",
	// Years
	"%y": "06", "%Y": "2006",
	// Hours, Minutes, Seconds
	"%h": "15", "%H": "03",
	"%i": "04", "%s": "05",
	// AM/PM
	"%a": "pm", "%A": "PM",
}

func init() {
	err := SetTimezone(os.Getenv("TZ"))
	if err != nil {
		panic("invalid timezone set in TZ environment variable: " + err.Error())
	}
}

func Freeze(time int64) {
	fakeTime = &time
}

func Unfreeze() {
	fakeTime = nil
}

func Now() int64 {
	if fakeTime != nil {
		return *fakeTime
	}

	return time.Now().UnixMilli()
}

func Sleep(milliseconds int64) {
	if fakeTime != nil {
		*fakeTime += milliseconds
		return
	}

	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func Parse(dateString, layout string) (int64, error) {
	originalLayout := layout

	for k, v := range layoutReplacements {
		layout = strings.ReplaceAll(layout, k, v)
	}

	t, err := time.Parse(layout, dateString)
	if err != nil {
		message := err.Error()
		if strings.Contains(message, "cannot parse") {
			return 0, fmt.Errorf("cannot parse '%s' as '%s'", dateString, originalLayout)
		}

		if strings.Contains(message, "out of range") {
			return 0, fmt.Errorf("date '%s' is out of range for format '%s'", dateString, originalLayout)
		}

		if strings.Contains(message, "missing") {
			return 0, fmt.Errorf("date '%s' is missing components for format '%s'", dateString, originalLayout)
		}

		if strings.Contains(message, "extra text") {
			return 0, fmt.Errorf("date '%s' has extra text for format '%s'", dateString, originalLayout)
		}

		return 0, errors.New("unknown error: " + message)
	}

	return t.UnixMilli(), nil
}

func Format(timestamp int64, layout string) string {
	time := time.UnixMilli(timestamp)

	for k, v := range layoutReplacements {
		layout = strings.ReplaceAll(layout, k, v)
	}

	return time.Format(layout)
}

func ResetTimezone() {
	time.Local = localTimezone
}

func SetTimezone(timezone string) error {
	if timezone == "" || strings.ToLower(timezone) == "local" {
		ResetTimezone()
		return nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("unknown time zone `%s`", timezone)
	}

	time.Local = loc
	return nil
}
