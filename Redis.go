package main

import (
	"gopkg.in/redis.v5"
	"strconv"
)

func storeData(key string, data string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()
	err := client.Set(strconv.Itoa(getId()) + ":" + key, data, 0).Err()
	if err != nil {
		panic(err)
	}
}

func getData(key string) string {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	val, _ := client.Get(strconv.Itoa(getId()) + ":" + key).Result()
	return val
}

func clearData() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	client.Del(getAddress())
}