package scheduler

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type OChainScheduler struct {
	scheduler gocron.Scheduler
}

func NewScheduler() (OChainScheduler, error) {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		return OChainScheduler{}, err
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func(a string, b int) {
				// do things
			},
			"hello",
			1,
		),
	)

	if err != nil {
		return OChainScheduler{}, err
	}

	// start the scheduler
	s.Start()

	return OChainScheduler{
		scheduler: s,
	}, nil
}
