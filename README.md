# wordlist-sanitizer
Remove Offensive and Profane Words from Wordlists

# About
`wordlist-sanitizer` will create a copy of a file appended with`-clean` that has had a specified list of bad words removed.
If the input is a directory, `wordlist-sanitizer` will recursively create a clone of the directory (directory names also appended with `-clean`) with all files inside sanitized.

# Installation
Ensure that Golang is installed, and the GOPATH variable is in your PATH

```bash
git clone https://github.com/gaberust/wordlist-sanitizer
cd wordlist-sanitizer
```

Windows:
```ps
.\install.bat
```

*nix:
```bash
chmod +x install.sh
./install.sh
```

# Usage
```bash
$ wordlist-sanitizer -h
Usage of wordlist-sanitizer:
  -bad string
        The list of words to be stripped. (default "[EXE_PATH]\\bad-words.txt")
  -out string
        The output directory. (default ".")
  -path string
        The path of the target file or directory.
        May also be passed after all flags as a positional argument. (default ".")
  -threads int
        Concurrent worker count. (default 100)
```

Example:
```bash
$ wordlist-sanitizer -threads 100000 SecLists\Usernames\xato-net-10-million-usernames.txt
SecLists\Usernames\xato-net-10-million-usernames.txt
1101033 bad words were removed out of 8295455 words.
```
