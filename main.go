package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"os"
	"strings"
	"time"
)

func readEnv(name string) (string, error) {
	value := os.Getenv(name)
	if len(value) == 0 {
		return "", fmt.Errorf("missing %s env", name)
	}
	return value, nil
}

func main() {
	addresses, err := readEnv("ADDRESSES")
	if err != nil {
		fmt.Println(err)
		return
	}
	password, err := readEnv("PASSWORD")
	if err != nil {
		fmt.Println(err)
		return
	}
	masterName, err := readEnv("MASTER")
	if err != nil {
		fmt.Println(err)
		return
	}

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: strings.Split(addresses, " "),
		Password:      password,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(fmt.Errorf("ping: %w", err))
		return
	}

	_, err = rdb.Set(context.Background(), "k1", "v1", time.Minute).Result()
	if err != nil {
		fmt.Println(fmt.Errorf("set: %w", err))
		return
	}
	<-time.After(time.Second * 30)
	value, err := rdb.Get(context.Background(), "k1").Result()
	if err != nil {
		fmt.Println(fmt.Errorf("get: %w", err))
		return
	}
	if value != "v1" {
		fmt.Println("Value mismatch")
		return
	}
	fmt.Println("Everything seems ok")
}
