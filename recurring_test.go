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

func ExampleRecurring() {
	// set a recurring tick at 01:02:03 UTC for the next 3 days
	r := New(1, 2, 3, 0)
	defer r.Stop()
	for i := 0; i < 3; i++ {
		<-r.C
		// time will now be 01:02:03 UTC
	}
}

func TestRecurring(t *testing.T) {
	now := time.Now().UTC()
	r := New(now.Hour(), now.Minute(), now.Second(), now.Nanosecond()+5e6)
	defer r.Stop()
	// we don't want to wait 24 hours for the test to complete, so we cheat and
	// change the deadline. This won't affect the first event, but will change
	// the second one.
	r.nsec += 5e6
	<-r.C
	elapsed := time.Since(now)
	if elapsed.Nanoseconds()/1e6 != 5 {
		t.Fatalf("Bad duration, expected 5ms, got %d", elapsed.Nanoseconds()/1e6)
	}
	now = time.Now().UTC()
	<-r.C
	elapsed = time.Since(now)
	msec := elapsed.Nanoseconds() / 1e6
	// there is a little slop, so look for a range
	if msec < 4 || msec > 6 {
		t.Errorf("Bad duration, expected 5ms, got %d", msec)
	}
}
