package {{.PackageName}}

import (
	"context"

	"github.com/bingo-project/component-base/log"
	"github.com/go-redsync/redsync/v4"

	"{{.RootPackage}}/{{.StorePath}}"
)

type {{.StructName}}Watcher struct {
	ctx   context.Context
	mutex *redsync.Mutex
}

// Spec is parsed using the time zone of clean Cron instance as the default.
func (w *{{.StructName}}Watcher) Spec() string {
	return "@every 1m"
}

// Init initializes the watcher for later execution.
func (w *{{.StructName}}Watcher) Init(ctx context.Context, rs *redsync.Mutex, config interface{}) error {
	*w = {{.StructName}}Watcher{
		ctx:   ctx,
		mutex: rs,
	}

	return nil
}

// Run runs the watcher job.
func (w *{{.StructName}}Watcher) Run() {
	if err := w.mutex.Lock(); err != nil {
		log.C(w.ctx).Infow("{{.StructName}}Watcher already run.")

		return
	}
	defer func() {
		if _, err := w.mutex.Unlock(); err != nil {
			log.C(w.ctx).Errorw("could not release {{.StructName}}Watcher lock. err: %v", err)

			return
		}
	}()

    // Do your job here.
}
