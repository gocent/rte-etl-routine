package scheduler

import (
	"errors"
	"log"
	"rte-etl-routine/executor"
	"sync"
	"time"
)

type Scheduler struct {
	frequency time.Duration
	executors []executor.Executor
	state     State
}

type State int8

const (
	RUNNING State = 0
	PAUSED        = 1
	STOPPED       = 2
)

func NewScheduler(frequency time.Duration) *Scheduler {
	return &Scheduler{
		frequency: frequency,
		state:     STOPPED,
	}
}

func (s *Scheduler) Add(e executor.Executor) *Scheduler {
	s.executors = append(s.executors, e)
	return s
}

func (s *Scheduler) Start(immediate bool) error {
	s.state = RUNNING
	if len(s.executors) == 0 {
		return errors.New("no job to schedule")
	}

	if immediate == true {
		s.execute()
	}

	for s.state != STOPPED {
		time.Sleep(s.frequency)
		if s.state == RUNNING {
			s.execute()
		}
	}
	return nil
}

func (s *Scheduler) execute() {
	var wg sync.WaitGroup
	for _, e := range s.executors {
		wg.Add(1)
		go func(e executor.Executor) {
			defer wg.Done()
			err := e.Execute()
			if err != nil {
				log.Fatal(err)
			}
		}(e)
	}
	wg.Wait()
}

func (s *Scheduler) Pause() {
	s.state = PAUSED
}

func (s *Scheduler) Resume() {
	s.state = RUNNING
}

func (s *Scheduler) Stop() {
	s.state = STOPPED
}
