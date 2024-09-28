package watch

import "time"

// trigger is a debouncer for events.
//
// Example usage:
//
//	trigger := newTrigger(time.Second)
//	for {
//	  select {
//	  case evt := <-events: // Listen on your inbox.
//	    trigger.Ping()
//	  case <-trigger.C: // This gets triggered only after inbox has been quiet for a second.
//	    trigger.Ack() // This is important and avoids race conditions.
//	    // do something
//	  }
//	}
type trigger struct {
	C        <-chan time.Time
	timer    *time.Timer
	running  bool
	interval time.Duration
}

func newTrigger(interval time.Duration) *trigger {
	timer := time.NewTimer(time.Second)
	if !timer.Stop() {
		<-timer.C
	}
	return &trigger{
		C:        timer.C,
		timer:    timer,
		running:  false,
		interval: interval,
	}
}

func (t *trigger) Ping() {
	if t.running {
		if !t.timer.Stop() {
			<-t.timer.C
		}
	}
	t.running = true
	t.timer.Reset(t.interval)
}

func (t *trigger) Ack() {
	t.running = false
}
