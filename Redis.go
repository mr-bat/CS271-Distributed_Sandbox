package main

import "gopkg.in/redis.v5"

func storeData(data string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()
	err := client.Set(getAddress(), data, 0).Err()
	if err != nil {
		panic(err)
	}
}

func getData() string {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	val, _ := client.Get(getAddress()).Result()
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