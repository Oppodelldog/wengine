package main

import (
	"flag"
	"fmt"
	"github.com/Oppodelldog/wengine/packer"
	"io/ioutil"
	"os"
)

func main() {
	var (
		filename = flag.String("f", "", "-f out.dat")
		pattern  = flag.String("p", "%003d", "-p %003d.txt")
	)
	flag.Parse()

	if *filename == "" {
		fmt.Println("parameter -f is required")
		os.Exit(1)
	}

	packedFile, err := packer.Read(*filename)
	exitOnError(err)

	for f := range packedFile.Files() {
		exitOnError(ioutil.WriteFile(fmt.Sprintf(*pattern, f.Index), f.Content, 0644))
	}
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
