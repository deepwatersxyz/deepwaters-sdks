package util

import (
	"time"
)

type Periodic struct {
	Timeout        time.Duration
	Interval       time.Duration
	ReportInterval time.Duration
	ExitOnSuccess  bool
	start          time.Time
}

type periodicHelper struct {
	stop        chan bool
	report      chan bool
	timeoutQuit chan bool
	reportQuit  chan bool
}

func (p *Periodic) prepare() periodicHelper {
	h := periodicHelper{
		stop:        make(chan bool, 1),
		report:      make(chan bool, 1),
		timeoutQuit: make(chan bool, 1),
		reportQuit:  make(chan bool, 1),
	}

	if p.Timeout > 0 {
		go p.runTimeout(h)
	}

	if p.ReportInterval > 0 {
		go p.runReport(h)
	}

	return h
}

func (p *Periodic) GetStartTime() time.Time {
	return p.start
}

func (p *Periodic) Run(f func() error, reportF func(attempts int, err error)) error {
	p.start = time.Now()
	attempts := 1
	var err error
	if err = f(); err == nil && p.ExitOnSuccess {
		return nil
	} else if err != nil && reportF != nil {
		reportF(attempts, err)
	}

	h := p.prepare()
	defer func() {
		h.timeoutQuit <- true
		h.reportQuit <- true
	}()

	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()
	for running := true; running; {
		select {
		case <-ticker.C:
			attempts += 1
			if err = f(); err == nil && p.ExitOnSuccess {
				running = false
			} else if err != nil && p.ReportInterval <= 0 && reportF != nil {
				reportF(attempts, err)
			}
		case <-h.report:
			if reportF != nil {
				reportF(attempts, err)
			}
		case <-h.stop:
			running = false
			break
		}
	}

	return err
}

func (p *Periodic) runTimeout(h periodicHelper) {
	select {
	case <-time.After(p.Timeout):
		h.stop <- true
		break
	case <-h.timeoutQuit:
		break
	}
}

func (p *Periodic) runReport(h periodicHelper) {
	ticker := time.NewTicker(p.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.report <- true
		case <-h.reportQuit:
			return
		}
	}
}
