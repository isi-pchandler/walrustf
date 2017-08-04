package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

var conn *redis.Client
var collector = flag.String("collector", "localhost", "walrus instance to use")

func main() {

	log.SetFlags(0)
	flag.Parse()
	
	dbConnect()

	pubsub := conn.PSubscribe("__key*__:*")
	defer pubsub.Close()

	for {
		data, err := pubsub.Receive()
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		msg, ok := data.(*redis.Message)
		if ok {
			if strings.HasSuffix(msg.Channel, "set") {
				if strings.HasSuffix(msg.Payload, "~seq~") {
					continue
				}
				showKey(msg.Payload)
			}
		}
	}

}

func showKey(key string) {

	value, err := conn.Get(key).Result()
	if err != nil {
		log.Printf("get_error: key=%s %v", key, err)
		return
	}

	key_parts := strings.Split(key, ":")
	if len(key_parts) != 4 {
		log.Printf("bad key format %s", key)
		return
	}

	sec, err := strconv.Atoi(key_parts[2])
	if err != nil {
		log.Printf("bad key format - unable to parse seconds %s", key)
		return
	}
	nsec, err := strconv.Atoi(key_parts[3])
	if err != nil {
		log.Printf("bad key format - unable to parse nanoseconds %s", key)
		return
	}

	t := time.Unix(int64(sec), int64(nsec))

	//log.Printf("%s = %s", key, value)

	parts := strings.Split(value, ":::")
	if len(parts) != 2 {
		log.Printf("bad value format %s", value)
		return
	}

	msg := ""
	switch parts[0] {
	case "ok":
		msg = green(parts[1])
	case "warning":
		msg = yellow(parts[1])
	default:
		msg = red(parts[1])
	}

	log.Printf("%s:%s %s %s",
		cyan(key_parts[0]), blue(key_parts[1]),
		msg,
		t.Format("15:04:05.000"),
	)

}

func dbConnect() {

	log.Printf("connecting to %s", *collector)

	conn = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", *collector),
	})

	if conn == nil {
		log.Fatal("could not connect to redis")
	}

	_, err := conn.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

}

var blueb = color.New(color.FgBlue, color.Bold).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var cyanb = color.New(color.FgCyan, color.Bold).SprintFunc()
var greenb = color.New(color.FgGreen, color.Bold).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var redb = color.New(color.FgRed, color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var bold = color.New(color.Bold).SprintFunc()
