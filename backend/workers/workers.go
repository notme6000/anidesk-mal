package workers

import (
	"log"
	"time"
)

type Worker struct {
	name     string
	interval time.Duration
	task     func() error
	stopCh   chan struct{}
}

func New(name string, interval time.Duration, task func() error) *Worker {
	return &Worker{
		name:     name,
		interval: interval,
		task:     task,
		stopCh:   make(chan struct{}),
	}
}

func (w *Worker) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		if err := w.task(); err != nil {
			log.Printf("[worker:%s] error: %v", w.name, err)
		}

		for {
			select {
			case <-ticker.C:
				if err := w.task(); err != nil {
					log.Printf("[worker:%s] error: %v", w.name, err)
				}
			case <-w.stopCh:
				log.Printf("[worker:%s] stopped", w.name)
				return
			}
		}
	}()
	log.Printf("[worker:%s] started (interval: %s)", w.name, w.interval)
}

func (w *Worker) Stop() {
	close(w.stopCh)
}
