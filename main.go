package main

import (
	"log"

	"github.com/tianhongw/tinyid/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
