package global

import (
	"fmt"
	"log"

	"github.com/gin-contrib/sessions/redis"
)

var Store redis.Store

func InitRedis()  {
	addr := Config.Redis.RedisAddr
	port := Config.Redis.RedisPort
	passowrd := Config.Redis.RedisPassword
	selectstr := Config.Server.Secret
	poolConns := Config.Redis.RedisPoolConns
	var err error
	Store, err = redis.NewStore(poolConns, "tcp", fmt.Sprintf("%s:%s",addr,port), passowrd, []byte(selectstr))
	if err != nil {
		log.Fatalln("redis error:",err)
	}
}