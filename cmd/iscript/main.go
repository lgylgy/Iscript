package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lgylgy/iscript/s4"
)

var (
	encode = flag.Bool("e", false, "true -> encode / false -> decode")
)

func usage() {
	fmt.Fprintf(os.Stderr, "no description...\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	config, err := s4.LoadConfiguration(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	message, err := s4.Run(*encode, config)
	if err != nil {
		log.Fatal(err)
	}

	if message != "" {
		log.Println(message)
	}
}
