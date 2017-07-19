package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var conn *redis.Client
var tds *TestDefs
var collector = flag.String("collector", "localhost", "walrus instance to use")

func main() {

	log.SetFlags(0)
	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
	}

	dbConnect()
	loadTest(flag.Args()[1])

	switch flag.Args()[0] {
	case "launch":
		doLaunch()
	case "watch":
		doWatch()
	default:
		usage()
	}

}

func doLaunch() {

	log.Printf(bold("running %d tests"), len(tds.Tests))

	go launchAll()

	doWatch()

}

func launchAll() {

	for _, t := range tds.Tests {
		err := launch(t)
		if err != nil {
			log.Printf(red("%v\n"), err)
			continue
		}
	}

	writeResults(tds)

}

func doWatch() {

	for _, t := range tds.Tests {
		color.Blue("[%s]", t.Name)
		wait(t)
	}

}

func launch(t TestSpec) error {

	//extract and execute the launch command for the test
	//fmt.Printf("running launch script...")
	cmd_line := strings.Fields(t.Launch)
	cmd := exec.Command(cmd_line[0], cmd_line[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s", out)
		return fmt.Errorf("launch command failed: %v", err)
	}
	return nil
}

func wait(t TestSpec) {

	begin := time.Now()
	end := time.Now().Add(time.Duration(t.Timeout) * time.Second)
	for time.Now().Before(end) {
		elapsed := time.Now().Sub(begin).Seconds()
		remaining := end.Sub(time.Now()).Seconds()
		fmt.Printf("\r%f                  ", remaining)
		if finished(t, begin, elapsed) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Print("         \r")

	log.Printf(red("timeout!\n"))

}

func finished(t TestSpec, start time.Time, elapsed float64) bool {

	_, err := conn.Ping().Result()
	if err != nil {
		log.Fatal("failed to connect to redis: %v", err)
	}

	//TODO check failure condition
	result := true
	for i, _ := range t.Success {
		c := &t.Success[i]
		if !c.Satisfied {
			testCondition(t.Name, c, conn, start)
			if c.Satisfied {
				logTestPassed(c, elapsed)
			} else {
				result = false
			}
		}
	}
	return result

}

func logTestPassed(c *Condition, elapsed float64) {
	log.Printf("\r%s: %s @%ss",
		cyan(c.Who),
		green(c.Message),
		bold(fmt.Sprintf("%f", elapsed)),
	)
}

func testCondition(test string, c *Condition, db *redis.Client, start time.Time) {

	var ts time.Time

	match := fmt.Sprintf("%s:%s:*", test, c.Who)
	iter := db.Scan(0, match, 0).Iterator()

	for iter.Next() {
		v := iter.Val()

		parts := strings.Split(v, ":")
		if len(parts) != 4 {
			log.Printf("bad key format %s", v)
			continue
		}

		sec, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Printf("bad key format - unable to parse seconds %s", v)
			continue
		}
		usec, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Printf("bad key format - unable to parse microseconds %s", v)
			continue
		}
		ts = time.Unix(int64(sec), 0)
		ts = ts.Add(time.Duration(usec) * time.Microsecond)

		val, _ := db.Get(v).Result()
		ss := strings.Split(val, ":::")
		if len(ss) == 2 && ss[0] == c.Status && ss[1] == c.Message {
			if ts.After(start) {
				c.Satisfied = true
			}
		}
	}

}

func loadTest(filename string) {

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading '%s': %v", filename, err)
	}

	tds = &TestDefs{}
	err = json.Unmarshal(buf, tds)
	if err != nil {
		log.Fatalf("error parsing '%s': %v", filename, err)
	}

}

func writeResults(tds *TestDefs) {

	buf, err := json.MarshalIndent(tds, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(tds.Name+"_results.json", buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func dbConnect() {

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

func usage() {
	s := red("usage:\n")
	s += fmt.Sprintf("  %s [-collector=host] ( %s | %s ) tests.json",
		blue("wtf"),
		green("launch"),
		green("watch"),
	)
	log.Fatal(s)
}

type TestSpec struct {
	Name, Launch  string
	Timeout       uint
	Success, Fail []Condition
}

type TestDefs struct {
	Name  string
	Tests []TestSpec
}

type Condition struct {
	Status, Who, Message string
	Satisfied            bool
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
