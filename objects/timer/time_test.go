package timer

import (
	"fmt"
	"testing"
	"time"
)

func resetTimePackage() {
	Unfreeze()
	ClearTimers()
	ResetTimezone()
}

func TestFreezeUnfreeze(t *testing.T) {
	defer resetTimePackage()

	if fakeTime != nil {
		t.Errorf("Expected fakeTime to be nil at start of test")
	}

	Freeze(1234567890)

	if fakeTime == nil || *fakeTime != 1234567890 {
		t.Errorf("Expected fakeTime to be 1234567890 after Freeze")
	}

	Unfreeze()

	if fakeTime != nil {
		t.Errorf("Expected fakeTime to be nil after Unfreeze")
	}
}

func TestTimeTravel(t *testing.T) {
	defer resetTimePackage()

	err := TimeTravel(10)
	if err == nil {
		t.Errorf("Expected error when calling TimeTravel without freezing time")
	}

	Freeze(1000)

	err = TimeTravel(10)
	if err != nil {
		t.Errorf("Did not expect error when calling TimeTravel after freezing time")
	}

	if *fakeTime != 1010 {
		t.Errorf("Expected fakeTime to be 1010 after time travel, got %d", *fakeTime)
	}
}

func TestTimeTravelWithTimers(t *testing.T) {
	defer resetTimePackage()

	Freeze(0)

	timerFired := false
	StartDelayedTimer(func() {
		timerFired = true
	}, 5)

	tickerCount := 0
	StartScheduledTimer(func() {
		tickerCount++
	}, 3)

	err := TimeTravel(10)
	if err != nil {
		t.Errorf("Did not expect error when calling TimeTravel after freezing time")
	}

	if !timerFired {
		t.Errorf("Expected timer to have fired after time travel")
	}

	if tickerCount != 3 {
		t.Errorf("Expected ticker to have fired 3 times after time travel, got %d", tickerCount)
	}
}

func TestClearTimers(t *testing.T) {
	defer resetTimePackage()

	timers["testTimer"] = StartDelayedTimer(func() {
		panic("This should not be called")
	}, 10)
	tickers["testTicker"] = StartScheduledTimer(func() {
		panic("This should not be called")
	}, 10)
	fakeTimers = append(fakeTimers, &FakeTime{
		callback: func() {},
		delay:    10,
		lastCall: 0,
	})

	ClearTimers()

	if len(timers) != 0 {
		t.Errorf("Expected timers map to be empty after ClearTimers")
	}

	if len(tickers) != 0 {
		t.Errorf("Expected tickers map to be empty after ClearTimers")
	}

	if len(fakeTimers) != 0 {
		t.Errorf("Expected fakeTimers slice to be empty after ClearTimers")
	}

	// Ensure no panics occur when original timers would have fired
	time.Sleep(15 * time.Millisecond)
}

func TestNow(t *testing.T) {
	defer resetTimePackage()

	realNow := time.Now().UnixMilli()
	now := Now()
	if now < realNow || now > realNow+5 {
		t.Errorf("Expected Now() to return current time in milliseconds, got %d", now)
	}

	Freeze(1234567890)
	fakeNow := Now()

	if fakeNow != 1234567890 {
		t.Errorf("Expected Now() to return frozen time 1234567890, got %d", fakeNow)
	}
}

func TestSleep(t *testing.T) {
	defer resetTimePackage()

	before := Now()
	Sleep(10)
	after := Now()

	diff := after - before
	if diff < 10 || diff > 12 {
		t.Errorf("Expected at least 10 milliseconds to have passed during Sleep, got %d", diff)
	}

	Freeze(1000)
	Sleep(10)

	if *fakeTime != 1010 {
		t.Errorf("Expected fakeTime to be 1010 after Sleep, got %d", *fakeTime)
	}
}

func TestSetTimezone(t *testing.T) {
	defer resetTimePackage()

	err := SetTimezone("Europe/Copenhagen")
	if err != nil {
		t.Errorf("Did not expect error when setting valid timezone, got %v", err)
	}

	loc, _ := time.LoadLocation("Europe/Copenhagen")
	now := time.Now().In(loc)
	zenNow := time.Now().In(time.Local)
	if now.Hour() != zenNow.Hour() {
		t.Errorf("Expected local time to match set timezone hour, got %d and %d", zenNow.Hour(), now.Hour())
	}

	err = SetTimezone("Invalid/Timezone")
	if err == nil {
		t.Errorf("Expected error when setting invalid timezone")
	}
}

func TestResetTimezone(t *testing.T) {
	if time.Local == nil {
		t.Errorf("Expected time.Local to not be nil")
	}

	time.Local = nil

	ResetTimezone()
	if time.Local != localTimezone {
		t.Errorf("Expected time.Local to be reset to localTimezone")
	}
}

func TestStartDelayedTimerWithoutFreezing(t *testing.T) {
	defer resetTimePackage()

	timesFired := 0
	timer := StartDelayedTimer(func() {
		timesFired++
	}, 5)

	if _, exists := timers[fmt.Sprintf("%p", timer)]; !exists {
		t.Errorf("Expected timer to be added to timers map")
	}

	time.Sleep(15 * time.Millisecond)

	if timesFired != 1 {
		t.Errorf("Expected timer to have fired once, got %d", timesFired)
	}

	StopDelayedTimer(timer)

	if _, exists := timers[fmt.Sprintf("%p", timer)]; exists {
		t.Errorf("Expected timer to be removed from timers map after stopping")
	}
}

func TestStartDelayedTimerWithFreezing(t *testing.T) {
	defer resetTimePackage()
	Freeze(0)

	timesFired := 0
	timer := StartDelayedTimer(func() {
		timesFired++
	}, 10)

	if _, exists := timers[fmt.Sprintf("%p", timer)]; exists {
		t.Errorf("Expected timer to not be in timers map when time is frozen")
	}

	TimeTravel(100)

	stopped := StopDelayedTimer(timer)
	if !stopped {
		t.Errorf("Expected StopDelayedTimer to return true for valid timer")
	}

	err := TimeTravel(20)
	if err != nil {
		t.Errorf("Did not expect error when calling TimeTravel after freezing time")
	}

	if timesFired != 1 {
		t.Errorf("Expected timer to have fired once, got %d", timesFired)
	}
}

func TestStartScheduledTimerWithoutFreezing(t *testing.T) {
	defer resetTimePackage()

	timesFired := 0
	ticker := StartScheduledTimer(func() {
		timesFired++
	}, 5)

	if _, exists := tickers[fmt.Sprintf("%p", ticker)]; !exists {
		t.Errorf("Expected ticker to be added to tickers map")
	}

	time.Sleep(20 * time.Millisecond)

	if timesFired < 2 || timesFired > 4 {
		t.Errorf("Expected ticker to have fired between 2 and 4 times, got %d", timesFired)
	}

	StopScheduledTimer(ticker)

	if _, exists := tickers[fmt.Sprintf("%p", ticker)]; exists {
		t.Errorf("Expected ticker to be removed from tickers map after stopping")
	}
}

func TestStartScheduledTimerWithFreezing(t *testing.T) {
	defer resetTimePackage()
	Freeze(0)

	timesFired := 0
	ticker := StartScheduledTimer(func() {
		timesFired++
	}, 10)

	if _, exists := tickers[fmt.Sprintf("%p", ticker)]; exists {
		t.Errorf("Expected ticker to not be in tickers map when time is frozen")
	}

	err := TimeTravel(50)
	if err != nil {
		t.Errorf("Did not expect error when calling TimeTravel after freezing time")
	}

	stopped := StopScheduledTimer(ticker)
	if !stopped {
		t.Errorf("Expected StopScheduledTimer to return true for valid ticker")
	}

	if timesFired != 5 {
		t.Errorf("Expected ticker to have fired 5 times, got %d", timesFired)
	}
}
