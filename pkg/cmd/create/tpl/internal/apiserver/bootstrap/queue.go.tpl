package bootstrap

import (
	"github.com/redis/go-redis/v9"

	"{[.RootPackage]}/internal/apiserver/facade"
	"{[.RootPackage]}/pkg/queue"
)

func InitQueue() {
	rds := &redis.Options{
		Addr:     facade.Config.Redis.Host,
		Username: facade.Config.Redis.Username,
		Password: facade.Config.Redis.Password,
		DB:       facade.Config.Redis.Database,
	}

	facade.Queue, facade.Worker = queue.NewQueue(rds)
}
