package scheduler

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/ochain-gg/ochain-network/config"
	"github.com/ochain-gg/ochain-network/engine/database"
)

type OChainScheduler struct {
	Scheduler gocron.Scheduler
	db        *database.OChainDatabase
	cfg       config.OChainConfig
}

func NewScheduler(cfg config.OChainConfig, db *database.OChainDatabase) (OChainScheduler, error) {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		return OChainScheduler{}, err
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(
			15*time.Second,
		),
		gocron.NewTask(
			func(db *database.OChainDatabase) {
				CheckAndHandlePortalUpdate(cfg, db)
			},
			db,
		),
	)

	if err != nil {
		return OChainScheduler{}, err
	}

	return OChainScheduler{
		Scheduler: s,
		db:        db,
		cfg:       cfg,
	}, nil
}
