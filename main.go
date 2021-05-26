package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	defaultBadPath := filepath.Join(filepath.Dir(ex), "bad-words.txt")

	var inPath string
	flag.StringVar(&inPath, "path", ".", "The path of the target file or directory.\n"+
		"May also be passed after all flags as a positional argument.")

	var outPath string
	flag.StringVar(&outPath, "out", ".", "The output directory.")

	var badPath string
	flag.StringVar(&badPath, "bad", defaultBadPath, "The list of words to be stripped.")

	flag.Parse()

	fmt.Println("path:", inPath)
	fmt.Println("out:", outPath)
	fmt.Println("bad:", badPath)

}
