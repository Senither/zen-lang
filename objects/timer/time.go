package timer

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type FakeTime struct {
	callback func()
	lastCall int64
	delay    int64
	interval int64
	isTicker bool
	stopped  bool
	id       string
}

var fakeTime *int64 = nil
var localTimezone *time.Location = time.Local

var timers = map[string]*time.Timer{}
var tickers = map[string]*time.Ticker{}

var fakeTimers = []*FakeTime{}

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

func TimeTravel(duration int64) error {
	if fakeTime == nil {
		return errors.New("cannot time travel outside of testing environments, time must be frozen")
	}

	for range duration {
		*fakeTime += 1

		for _, ft := range fakeTimers {
			if ft.stopped {
				continue
			}

			if !ft.isTicker {
				if (*fakeTime - ft.lastCall) >= ft.delay {
					ft.callback()
					ft.stopped = true
				}
			} else if ft.interval > 0 {
				elapsed := *fakeTime - ft.lastCall
				if elapsed > 0 && (elapsed%ft.interval) == 0 {
					ft.callback()
				}
			}
		}
	}

	return nil
}

func ClearTimers() {
	for _, timer := range timers {
		timer.Stop()
	}

	for _, ticker := range tickers {
		ticker.Stop()
	}

	timers = map[string]*time.Timer{}
	tickers = map[string]*time.Ticker{}
	fakeTimers = []*FakeTime{}
}

func Now() int64 {
	if fakeTime != nil {
		return *fakeTime
	}

	return time.Now().UnixMilli()
}

func Sleep(milliseconds int64) {
	if fakeTime != nil {
		TimeTravel(milliseconds)
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

func StartDelayedTimer(callback func(), delay int64) *time.Timer {
	if fakeTime != nil {
		timer := time.AfterFunc(time.Millisecond*time.Duration(delay), func() {})

		fakeTimers = append(fakeTimers, &FakeTime{
			callback: callback,
			lastCall: *fakeTime,
			delay:    delay,
			isTicker: false,
			id:       fmt.Sprintf("%p", timer),
		})

		return timer
	}

	timer := time.AfterFunc(time.Millisecond*time.Duration(delay), func() {
		callback()
	})

	timers[fmt.Sprintf("%p", timer)] = timer

	return timer
}

func StopDelayedTimer(timer *time.Timer) bool {
	if timer == nil {
		return false
	}

	id := fmt.Sprintf("%p", timer)
	if fakeTime != nil {
		for _, ft := range fakeTimers {
			if ft.id == id && !ft.isTicker {
				ft.stopped = true
				break
			}
		}
	}

	stopped := timer.Stop()
	delete(timers, id)

	return stopped
}

func StartScheduledTimer(callback func(), interval int64) *time.Ticker {
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))

	if fakeTime != nil {
		fakeTimers = append(fakeTimers, &FakeTime{
			callback: callback,
			lastCall: *fakeTime,
			interval: interval,
			isTicker: true,
			id:       fmt.Sprintf("%p", ticker),
		})

		return ticker
	}

	go func() {
		for range ticker.C {
			callback()
		}
	}()

	tickers[fmt.Sprintf("%p", ticker)] = ticker

	return ticker
}

func StopScheduledTimer(ticker *time.Ticker) {
	if ticker == nil {
		return
	}

	id := fmt.Sprintf("%p", ticker)

	if fakeTime != nil {
		for _, ft := range fakeTimers {
			if ft.id == id && ft.isTicker {
				ft.stopped = true
				break
			}
		}
	}

	ticker.Stop()
	delete(tickers, id)
}
