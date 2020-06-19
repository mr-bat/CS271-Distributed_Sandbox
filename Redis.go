package main

func storeData(key string, data string) {
	panic("minimal prototyping")

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	//defer client.Close()
	//err := client.Set(strconv.Itoa(getId()) + ":" + key, data, 0).Err()
	//if err != nil {
	//	panic(err)
	//}
}

func getData(key string) string {
	panic("minimal prototyping")

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	//defer client.Close()
	//
	//val, _ := client.Get(strconv.Itoa(getId()) + ":" + key).Result()
	//return val
}

func clearData() {
	panic("minimal prototyping")

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	//defer client.Close()
	//
	//client.Del(getAddress())
}