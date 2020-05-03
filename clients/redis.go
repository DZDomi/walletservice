package clients

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"time"
)

var client *redis.Client
var service string

func InitRedis(servicename string) {
	service = servicename
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func generateKey(prefix string, key string) string {
	return fmt.Sprintf("service:%s:%s:%s", service, prefix, key)
}

func SetObject(prefix string, key string, object interface{}, duration time.Duration) error {
	keyToSet := generateKey(prefix, key)
	result, err := json.Marshal(object)
	if err != nil {
		log.Println("Error while marshaling: " + keyToSet + ", " + err.Error())
		return err
	}
	if _, err = client.Set(keyToSet, result, duration).Result(); err != nil {
		log.Println("Error while trying to set redis key: " + keyToSet + ", " + err.Error())
		return err
	}
	return nil
}

func GetObject(prefix string, key string, object interface{}) error {
	keyToGet := generateKey(prefix, key)
	val, err := client.Get(keyToGet).Result()
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		log.Println("Error while trying to get redis key: " + keyToGet + ", " + err.Error())
		return err
	}
	if err = json.Unmarshal([]byte(val), object); err != nil {
		log.Println("Error while unmarshalling: " + err.Error())
		return err
	}
	return nil
}

func DeleteObject(prefix string, key string) error {
	keyToDelete := generateKey(prefix, key)
	if _, err := client.Del(keyToDelete).Result(); err != nil {
		log.Println("Error while trying to delete redis key: " + keyToDelete + ", " + err.Error())
		return err
	}
	return nil
}
