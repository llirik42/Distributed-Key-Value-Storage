package main

import (
	"distributed-algorithms/config"
	"log"
)

func main() {
	cfg, err := config.NewConfiguration()

	if err != nil {
		panic(err)
	} else {
		log.Println(cfg)
	}
}
