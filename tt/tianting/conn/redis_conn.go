package conn

import (
	"context"
	"fmt"
	"tianting/comm"
	"time"

	"github.com/go-redis/redis/v8"
)

var log comm.ConsoleLogger

func InitRedisConn(key, valu string) interface{} {
	opt, err := redis.ParseURL("redis://:@localhost:6379/0")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 直接执行命令获取错误
	//err := rdb.Do(ctx, "set", "key", 10, "EX", 3600).Err()
	rdb.Do(ctx, "set", key, valu, "EX", 5).Err()
	fmt.Println(err)

	// 执行命令获取结果
	val, err := rdb.Do(ctx, "get", key).Result()
	if err != nil {
		fmt.Printf("get redis result error:%v\n", err)
	}
	log.Debug("get redis result is :%v", val)
	return val

}
