package client

import (
	"context"

	utils "github.com/angelthump/cache-main/utils"
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

func Initalize() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: utils.Config.Redis.Hostname,
		DB:   0,
	})
}
