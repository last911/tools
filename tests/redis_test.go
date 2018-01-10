package tests

import (
	"github.com/last911/tools/cache"
	"testing"
)

func TestRedis(t *testing.T) {
	redis := cache.NewRedis(&cache.Options{Addr: "127.0.0.1:6379"})
	err := redis.Set("name", "scnjl", 0).Err()
	if err != nil {
		t.Fatal(err)
	}
	name, err := redis.Get("name").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name:", name)
}
