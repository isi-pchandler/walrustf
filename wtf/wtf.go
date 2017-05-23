package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type TestSpec struct {
	Name, Launch  string
	Timeout       uint
	Success, Fail []Condition
}

type Condition struct {
	Status, Who, Message string
	Satisfied            bool
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		usage()
	}

	filename := os.Args[1]

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading '%s': %v", filename, err)
	}

	var testSpecs []TestSpec
	err = json.Unmarshal(buf, &testSpecs)
	if err != nil {
		log.Fatalf("error parsing '%s': %v", filename, err)
	}

	bw := color.New(color.FgWhite).Add(color.Underline).Add(color.Bold)
	bw.Printf("running %d tests\n", len(testSpecs))

	for _, t := range testSpecs {
		color.Blue("[%s]", t.Name)
		launch(t)
		wait(t)
	}
}

func launch(t TestSpec) {

}

func wait(t TestSpec) {

	red := color.New(color.FgRed)

	end := time.Now().Add(time.Duration(t.Timeout) * time.Second)
	for time.Now().Before(end) {
		remaining := end.Sub(time.Now()).Seconds()
		fmt.Printf("\r%f", remaining)
		if finished(t) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Print("         \r")

	red.Printf("timeout!\n")

}

func finished(t TestSpec) bool {

	green := color.New(color.FgGreen)

	conn := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := conn.Ping().Result()
	if err != nil {
		log.Fatal("failed to connect to redis: %v", err)
	}

	result := true
	for i, _ := range t.Success {
		c := &t.Success[i]
		if !c.Satisfied {
			testCondition(t.Name, c, conn)
			if c.Satisfied {
				green.Printf("\r%v\n", *c)
			} else {
				result = false
			}
		}
	}
	return result

}

func testCondition(test string, c *Condition, db *redis.Client) {
	//c.Satisfied = true
	match := fmt.Sprintf("%s:%s:*", test, c.Who)
	iter := db.Scan(0, match, 0).Iterator()
	for iter.Next() {
		if strings.HasSuffix(iter.Val(), "~time~") {
			continue
		}
		val, _ := db.Get(iter.Val()).Result()
		ss := strings.Split(val, ":::")
		if len(ss) == 2 && ss[0] == c.Status && ss[1] == c.Message {
			c.Satisfied = true
		}
	}
}

func usage() {
	log.Fatal("usage: wtf testfile.json")
}
