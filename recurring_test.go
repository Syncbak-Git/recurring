package recurring

import (
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	now, _ := time.Parse(time.RFC3339Nano, "2015-09-23T19:26:59.87654Z")
	d := deadline(now, now.Hour(), now.Minute(), now.Second()+1, now.Nanosecond())
	if int(d.Seconds()) != 1 {
		t.Errorf("Bad duration, expected 1 second, got %f", d.Seconds())
	}
	d = deadline(now, now.Hour(), now.Minute(), now.Second()-1, now.Nanosecond())
	if int(d.Seconds()) != 24*60*60-1 {
		t.Errorf("Bad duration, expected 24 hours - 1 second, got %f", d.Seconds())
	}
	d = deadline(now, now.Hour(), now.Minute()+1, now.Second(), now.Nanosecond())
	if int(d.Minutes()) != 1 {
		t.Errorf("Bad duration, expected 1 minute, got %f", d.Minutes())
	}
	d = deadline(now, now.Hour(), now.Minute()-1, now.Second(), now.Nanosecond())
	if int(d.Minutes()) != 1440-1 {
		t.Errorf("Bad duration, expected 24 hours - 1 minute, got %f", d.Minutes())
	}
	d = deadline(now, now.Hour()+1, now.Minute(), now.Second(), now.Nanosecond())
	if int(d.Hours()) != 1 {
		t.Errorf("Bad duration, expected 1 hour, got %f", d.Hours())
	}
	d = deadline(now, now.Hour()-1, now.Minute(), now.Second(), now.Nanosecond())
	if int(d.Hours()) != 23 {
		t.Errorf("Bad duration, expected 23 hours, got %f", d.Hours())
	}
	d = deadline(now, now.Hour(), now.Minute(), now.Second(), now.Nanosecond())
	if int(d.Hours()) != 24 {
		t.Errorf("Bad duration, expected 24 hours, got %f", d.Hours())
	}
}
