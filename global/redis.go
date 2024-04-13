package global

import (
	"fmt"

	"github.com/gin-contrib/sessions/redis"
)

var Store redis.Store

func InitRedis() {

	RetryWithExponentialBackoff(UseRedis, "Redis Connection", 5)

}

func UseRedis() error {
	var err error
	Store, err = redis.NewStore(Config.Redis.RedisPoolConns,
		"tcp", fmt.Sprintf("%s:%s", Config.Redis.RedisAddr, Config.Redis.RedisPort),
		Config.Redis.RedisPassword,
		[]byte(Config.Server.Secret))
	if err != nil {
		return err
	}
	return nil
}
