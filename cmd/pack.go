package main

import (
	"flag"
	"fmt"
	"github.com/Oppodelldog/wengine/packer"
	"os"
)

func main() {
	var (
		filename = flag.String("f", "", "-f out.dat")
		files    []string
	)
	flag.Parse()

	if *filename == "" {
		fmt.Println("parameter -f is required")
		os.Exit(1)
	}

	if len(os.Args[3:]) == 0 {
		fmt.Println("no files provided")
		os.Exit(2)
	}

	files = os.Args[3:]

	pf, err := packer.New(files)
	panicOnError(err)

	err = pf.Write(*filename)
	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
