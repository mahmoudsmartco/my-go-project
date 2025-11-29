package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // لو عندك باسورد اكتبها هنا
		DB:       0,
	})
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("❌ فشل الاتصال بـ Redis: %v", err))
	}
	fmt.Println("✅ تم الاتصال بـ Redis بنجاح")
}
