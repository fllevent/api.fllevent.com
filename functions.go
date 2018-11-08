package main

import (
	"flag"
	"fmt"
	"os"
)

func printHelp() {
	flag.PrintDefaults()
	os.Exit(0)
}

func printVersion() {
	fmt.Println("version number: " + versionNumber)
	os.Exit(0)
}
