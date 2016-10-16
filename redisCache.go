package main

import (
	"gopkg.in/redis.v5"
	"log"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func SaveToRedis(uuid string, payload string) {
	err := redisClient.Set(uuid, payload, 0).Err();

	if err != nil {
		log.Printf("Could not save payload! Error: %v\n", err)
	} else {
		log.Printf("Saved payload with uuid %s\n", uuid)
	}
}
