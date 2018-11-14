package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lestrrat-go/cron"
)

type CronWrapper map[string]doExecute

//cronDispatchJob schedule an executor every time unit
func cronDispatchJob(cw CronWrapper) {
	tab := cron.New()
	var mu sync.Mutex
	for cornPattern, executor := range cw {
		tab.Schedule(cornPattern, cron.JobFunc(func(context.Context) {
			mu.Lock()
			defer mu.Unlock()
			executor() // in case i pass map of 2 executors fistone is never executed howevere second is executed 2 times 
		}))
	}

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go tab.Run(ctx)
	<-sigCh
}
