package main

import (
	"flag"
	"fmt"
	"gopher/xkcd"
	"os"
)

var flagP = flag.Bool("p", false, "Parses xkcd.com to comics.json")
var flagR = flag.Bool("r", false, "Reads comics.json")

func main() {
	flag.Parse()

	if *flagP {
		xkcd.SaveToFile()
	} else if *flagR {
		os.Open("comics.json")
	} else {
		fmt.Println("Use -h for help")
	}
}
