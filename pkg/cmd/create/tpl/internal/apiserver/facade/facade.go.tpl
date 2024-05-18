package facade

import (
	"github.com/bingo-project/component-base/cache"
	"github.com/bingo-project/component-base/crypt"
	"github.com/bwmarrin/snowflake"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"{[.RootPackage]}/internal/apiserver/config"
	"{[.RootPackage]}/pkg/mail"
)

var (
	Config    config.Config
	AES       *crypt.AES
	Redis     *redis.Client
	Cache     *cache.CacheService
	Queue     *asynq.Client
	Worker    *asynq.Server
	Snowflake *snowflake.Node
	Mail      *mail.Mailer
)
