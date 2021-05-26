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

// A one copy read-only global list of "bad words"
var badWords []string

// Total count of words removed from all files
var badCount uint64

// Total count of words processed
var totalWords uint64

// Panic if an error is not `nil`
// `e` error: The error to check
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Check if the word in question contains any word from `badWords`
// `word` string: The word in question
//
// Returns: bool (false if contains bad word, otherwise true)
func checkWord(word string) bool {
	for _, s := range badWords {
		if strings.Contains(word, s) {
			return false
		}
	}
	return true
}

// Remove "bad words" from a file.
// Recurse when the filepath is a directory, with an actual file being the base case.
// `fpath` string: The filepath of the input file
// `opath` string: The output directory path
// `threads`  int: The maximum number of concurrent goroutines processing the file
func sanitizeList(fpath string, opath string, threads int) {
	// Print the current filepath cause leet
	fmt.Println(fpath)

	// Obtain file information for the current path
	info, err := os.Stat(fpath)
	check(err)

	// Check if file is directory (Base Case Check)
	if info.IsDir() {
		// File is directory, obtain directory contents
		dir, err := ioutil.ReadDir(fpath)
		check(err)

		// Call `sanitizeList` recursively on each listing in current path
		for _, f := range dir {
			sanitizeList(filepath.Join(fpath, f.Name()), opath, threads)
		}
	} else {
		// File is NOT a directory: Base Case Reached

		// Read file into memory
		content, err := ioutil.ReadFile(fpath)
		check(err)

		// Split content of file into array of whitespace separated words
		words := strings.Fields(string(content))

		// Append file word count to global word count
		totalWords += uint64(len(words))

		// Create channels for passing strings and queueing work
		results := make(chan string)
		queue := make(chan string)

		// If `threads` is greater than file word count,
		// Reduce threads to word count to remove excessive resource allocation
		if threads > len(words) {
			threads = len(words)
		}

		// Create Blocking WaitGroup for worker goroutines
		// Add number of threads to WaitGroup
		var waitGroup sync.WaitGroup
		waitGroup.Add(threads)

		// Create a goroutine for each "thread"
		for i := 0; i < threads; i++ {
			go func() {
				// Decrease WorkGroup before function exits
				defer waitGroup.Done()

				// Wait for words from work queue, breaks when `queue` closes
				for s := range queue {
					// Push word to results if good, otherwise increment global bad word counter
					if checkWord(s) {
						results <- s
					} else {
						badCount++
					}
				}
			}()
		}

		// Lock mutex to prevent parent from exiting prematurely
		var mutex sync.Mutex
		mutex.Lock()

		// Goroutine creating new file and processing results from workers
		go func() {
			// Unlock mutex when function is finished
			defer mutex.Unlock()

			// Split filepath into array of directory names
			tempPath := fpath
			if opath != "." {
				tempPath = filepath.Join(opath, fpath)
			}
			dirs := strings.Split(strings.ReplaceAll(tempPath, "\\", "/"), "/")

			// Append -clean to each directory and filename
			for i := 0; i < len(dirs); i++ {
				dirs[i] = dirs[i] + "-clean"
			}

			// Create the new directory structure
			err := os.MkdirAll(filepath.Join(dirs[:len(dirs)-1]...), os.ModePerm)
			check(err)

			// Create and open the new file
			f, err := os.Create(filepath.Join(dirs...))
			check(err)
			defer f.Close()

			// Create buffer for new file
			w := bufio.NewWriter(f)
			defer w.Flush()

			// Wait for words from results channel, and write them to the new file.
			// Breaks when `results` closes
			for s := range results {
				_, err := w.WriteString(s + "\n")
				check(err)
			}
		}()

		// Add all words to work queue and then immediately close queue channel
		for _, s := range words {
			queue <- s
		}
		close(queue)

		// Wait for workers to finish, then close results channel
		waitGroup.Wait()
		close(results)

		// Obtain lock on mutex
		// Prevents function from exiting while results are still being processed and file is still open
		mutex.Lock()
		mutex.Unlock()
	}
}

// Entry point
func main() {
	// Obtain filepath of executable to find path of default bad words list
	ex, err := os.Executable()
	check(err)
	defaultBadPath := filepath.Join(filepath.Dir(ex), "bad-words.txt")

	// Parse command line arguments with `flag` package
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

	// If extra arguments tail flags, use as `inPath`
	if len(flag.Args()) > 0 {
		inPath = strings.Join(flag.Args(), " ")
	}

	// Read bad words into memory
	badWordsContent, err := ioutil.ReadFile(badPath)
	check(err)

	// Split bad words into whitespace separated array (available globally)
	badWords = strings.Fields(string(badWordsContent))

	// Call `sanitizeList`. If the input path is a directory, `sanitizeList` will handle the recursion internally
	sanitizeList(inPath, outPath, threads)

	// After `sanitizeList` is done, print the number of removed/processed words cause leet
	fmt.Printf("%d bad words were removed out of %d words.", badCount, totalWords)
}

// BUY DOGECOIN
