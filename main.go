package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Result struct {
	Value   string
	Good    bool
	Channel chan string
}

var badWords []string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkWord(word string) bool {
	for _, s := range badWords {
		if strings.Contains(word, s) {
			return false
		}
	}
	return true
}

func sanitizeList(fpath string, opath string, threads int) {
	fmt.Println(fpath)
	info, err := os.Stat(fpath)
	check(err)
	if info.IsDir() {
		dir, err := ioutil.ReadDir(fpath)
		check(err)
		for _, f := range dir {
			sanitizeList(filepath.Join(fpath, f.Name()), opath, threads)
		}
	} else {
		content, err := ioutil.ReadFile(fpath)
		check(err)
		words := strings.Fields(string(content))

		results := make(chan string)
		queue := make(chan string)

		var workerGroup sync.WaitGroup
		workerGroup.Add(threads)

		var mutex sync.Mutex

		for i := 0; i < threads; i++ {
			go func() {
				defer workerGroup.Done()
				for s := range queue {
					if checkWord(s) {
						results <- s
					}
				}
			}()
		}

		go func() {
			mutex.Lock()
			defer mutex.Unlock()

			dirs := filepath.SplitList(fpath)
			for i := 0; i < len(dirs); i++ {
				dirs[i] = dirs[i] + "-clean"
			}
			out := filepath.Join(opath, filepath.Join(dirs[:len(dirs)-1]...))
			err := os.MkdirAll(out, os.ModePerm)
			check(err)

			f, err := os.Create(filepath.Join(out, dirs[len(dirs)-1]))
			check(err)
			defer f.Close()

			w := bufio.NewWriter(f)
			defer w.Flush()

			for s := range results {
				_, err := w.WriteString(s + "\n")
				check(err)
			}
		}()

		for _, s := range words {
			queue <- s
			fmt.Println(s)
		}
		close(queue)
		workerGroup.Wait()
		close(results)

		mutex.Lock()
		mutex.Unlock()
	}
}

func main() {
	ex, err := os.Executable()
	check(err)
	defaultBadPath := filepath.Join(filepath.Dir(ex), "bad-words.txt")

	var inPath string
	flag.StringVar(&inPath, "path", ".", "The path of the target file or directory.\n"+
		"May also be passed after all flags as a positional argument.")

	var outPath string
	flag.StringVar(&outPath, "out", ".", "The output directory.")

	var badPath string
	flag.StringVar(&badPath, "bad", defaultBadPath, "The list of words to be stripped.")

	var threads int
	flag.IntVar(&threads, "threads", 100, "Concurrent worker count.")

	flag.Parse()

	badWordsContent, err := ioutil.ReadFile(badPath)
	check(err)

	badWords = strings.Fields(string(badWordsContent))

	sanitizeList(inPath, outPath, threads)
}
