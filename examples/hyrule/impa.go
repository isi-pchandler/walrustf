package main

import (
	walrus "github.com/rcgoodfellow/walrustf/go"
	"log"
)

const (
	IP   = "192.168.147.3"
	PORT = 4747
)

func main() {
	wtf, err := walrus.NewClient("192.168.147.100", "hyrule", "impa")
	if err != nil {
		log.Fatal(err)
	}

	err = wtf.Ok("spicy")
	if err != nil {
		log.Fatal(err)
	}
	err = wtf.Warning("taco")
	if err != nil {
		log.Fatal(err)
	}
	err = wtf.Error("crunch")
	if err != nil {
		log.Fatal(err)
	}
}
