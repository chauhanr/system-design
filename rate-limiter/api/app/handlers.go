package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func (s *server) NotFoundHanlder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
}

var (
	TimeLimit  = 60 // seconds
	BucketSize = 5
)

func (s *server) ForwardingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathVars := mux.Vars(r)
		userId := pathVars["userId"]
		domain := pathVars["domain"]
		log.Printf("Domain: %s, UserId: %s\n", domain, userId)
		uniqueKey := fmt.Sprintf("%s:%s", domain, userId)

		resp, err, status := CheckRateLimit(s.redisClient, uniqueKey)
		if err != nil {
			resp = map[string]interface{}{
				"error": err.Error(),
			}
		}
		respond(w, r, status, resp)
	}
}

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("Response error in copying: %s\n", err)
	}
}

func CheckRateLimit(redisClient *redis.Client, uniqueKey string) (map[string]interface{}, error, int) {
	c, err := redisClient.Get(uniqueKey).Int64()
	if err == redis.Nil {
		err = redisClient.Set(uniqueKey, 1, time.Duration(TimeLimit)*time.Second).Err()
		if err != nil {
			log.Printf("Error setting time limit: %s\n", err)
			return nil, err, http.StatusInternalServerError
		}
		c = 1
	} else if err != nil {
		log.Printf("Error connecting to redis server: %s\n", err)
		return nil, err, http.StatusInternalServerError
	} else {
		if c >= int64(BucketSize) {
			return nil, errors.New("RateLimit Reached"), http.StatusTooManyRequests
		}
		c, err = redisClient.Incr(uniqueKey).Result()
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
	}
	dt := time.Now()
	res := map[string]interface{}{
		"data":  dt.Format("2006-01-02 15:40:05"),
		"count": c,
	}
	return res, nil, http.StatusOK
}
