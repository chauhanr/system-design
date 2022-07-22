package main

import (
	"log"
	"net/http"

	"github.com/chauhanr/system-design/rate-limiter/api/app"
	"github.com/go-redis/redis"
)

func main() {
	client := http.Client{}
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rClient.Ping().Result()

	if err != nil {
		log.Fatalf("Error Connecting to DB: %s\n", err)
	}
	log.Printf("Start Rate limiter app: response from DB %s\n", pong)
	app.Startup(&client, rClient)
}
