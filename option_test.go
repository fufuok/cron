package cron

import (
	"log"
	"strings"
	"testing"
	"time"
)

func TestWithCustomTime(t *testing.T) {
	timenow := func() time.Time {
		return time.Now()
	}
	New(WithCustomTime(timenow))
}

func TestWithLocation(t *testing.T) {
	c := New(WithLocation(time.UTC))
	if c.location != time.UTC {
		t.Errorf("expected UTC, got %v", c.location)
	}
}

func TestWithParser(t *testing.T) {
	var parser = NewParser(Dow)
	c := New(WithParser(parser))
	if c.parser != parser {
		t.Error("expected provided parser")
	}
}

func TestWithVerboseLogger(t *testing.T) {
	var buf syncWriter
	var logger = log.New(&buf, "", log.LstdFlags)
	c := New(WithLogger(VerbosePrintfLogger(logger)))
	if c.logger.(printfLogger).logger != logger {
		t.Error("expected provided logger")
	}

	c.AddFunc("@every 1s", func() {})
	c.Start()
	time.Sleep(OneSecond)
	c.Stop()
	out := buf.String()
	if !strings.Contains(out, "schedule,") ||
		!strings.Contains(out, "run,") {
		t.Error("expected to see some actions, got:", out)
	}
}

func TestWithSecondOptional(t *testing.T) {
	c := New(WithSecondOptional())
	id1, err := c.AddFunc("5 * * * *", func() {})
	if err != nil {
		t.Errorf("add func %v", err)
	}
	id2, err := c.AddFunc("* 5 * * * *", func() {})
	if err != nil {
		t.Errorf("add func %v", err)
	}
	c.Start()
	next1 := c.Entry(id1).Next
	next2 := c.Entry(id2).Next
	if next1 != next2 {
		t.Errorf("expect the same execution time: %s, %s", next1, next2)
	}
	c.Stop()
}
