package bootstrap

import (
	"github.com/bingo-project/component-base/log"
	"github.com/bingo-project/component-base/redis"

	"{[.RootPackage]}/internal/apiserver/cache"
	"{[.RootPackage]}/internal/apiserver/facade"
)

func InitCache() {
	r, err := redis.NewClient(facade.Config.Redis.Host, facade.Config.Redis.Password, facade.Config.Redis.Database)
	if err != nil {
		log.Errorw("init cache failed", "err", err)

		return
	}

	cache.NewCache(r)
}
