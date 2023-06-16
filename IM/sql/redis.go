package sql

import (
	"IM/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"time"
)

type Redisdao struct {
}

func Redstr(key string, value interface{}) string {
	//data := make([]byte, unsafe.Sizeof(value))
	data, _ := json.Marshal(value)
	set := utils.Red.Set(utils.Ctx, key, string(data), 1000*time.Second)
	return set.Val()
}
func Redget(key string) string {
	get := utils.Red.Get(utils.Ctx, key)
	if get.Err() != nil {
		return get.Err().Error()
	}
	RedExpire(key, 1)
	return get.Val()
}
func RedListGetAll(key string) []string {
	get := utils.Red.LRange(utils.Ctx, key, 0, -1)
	if get.Err() != nil {
		return []string{get.Err().Error()}
	}
	RedExpire(key, 1)
	return get.Val()
}

func RedListGroup(key string, value []int64) bool {
	for i := 0; i < len(value); i++ {
		if utils.Red.LPush(utils.Ctx, key, value[i]).Err() != nil {
			return false
		}
	}
	return true
}

func RedExpire(key string, level int) { //-1 永不过期  0  立马过期  1 增加1000秒   2 增加7200秒

	switch level {
	case -1:
		utils.Red.Expire(utils.Ctx, key, -1)
	case 1:
		utils.Red.Expire(utils.Ctx, key, 1000*time.Second)
	case 0:
		utils.Red.Expire(utils.Ctx, key, 1*time.Second)
	case 2:
		utils.Red.Expire(utils.Ctx, key, 7200*time.Second)
	}

}

func RedZRevRange(key string, start, stop int64) []string {
	revRange := utils.Red.ZRevRange(utils.Ctx, key, start, stop)
	if revRange.Err() != nil {
		fmt.Println(revRange.Err())
		return nil
		//return []string{revRange.Err().Error()}
	}
	return revRange.Val()
}

func RedZAdd(key string, score float64, value interface{}) bool {
	add := utils.Red.ZAdd(utils.Ctx, key, &redis.Z{score, value})
	if add.Err() != nil {
		fmt.Println(add.Err().Error())
		return false
	}
	RedExpire(key, 2)
	return true
}
