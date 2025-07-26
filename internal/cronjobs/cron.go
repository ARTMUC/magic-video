package cronjobs

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

type Func struct {
	Cron string
	F    func()
}

type Stopper interface {
	Stop() context.Context
}

func Start(funcs []Func) (Stopper, error) {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))
	for _, f := range funcs {
		_, err := c.AddFunc(f.Cron, f.F)
		if err != nil {
			return nil, fmt.Errorf("failed scheduling cron function: %w", err)
		}
	}

	c.Start()
	return c, nil
}
