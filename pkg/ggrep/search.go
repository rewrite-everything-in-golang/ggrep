package ggrep

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[1;31m"
	colorGreen  = "\033[32m"
	colorPurple = "\033[35m"
)

func SearchPath(path string, config *Config) (bool, error) {
	if path == "-" {
		return SearchReader(os.Stdin, "(standard input)", config)
	}

	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		if config.Recursive {
			return SearchDirectory(path, config)
		}
		return false, fmt.Errorf("%s is a directory", path)
	}

	return SearchFile(path, config)
}

func SearchDirectory(dir string, config *Config) (bool, error) {
	found := false

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		result, err := SearchFile(path, config)
		if err == nil && result {
			found = true
		}

		return nil
	})

	return found, err
}

func SearchFile(filename string, config *Config) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	return SearchReader(file, filename, config)
}

func SearchReader(reader io.Reader, filename string, config *Config) (bool, error) {
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var lines [][]byte
	var byteOffset int64

	for scanner.Scan() {
		line := make([]byte, len(scanner.Bytes()))
		copy(line, scanner.Bytes())
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	found := false
	matchCount := 0
	totalMatches := 0

	for idx, line := range lines {
		if config.MaxCount > 0 && matchCount >= config.MaxCount {
			break
		}

		isMatch := config.Regexp.Match(line)
		shouldPrint := isMatch != config.InvertMatch

		if shouldPrint {
			found = true
			matchCount++
			totalMatches++

			if config.Quiet {
				return true, nil
			}

			if config.FilesWithMatches {
				fmt.Println(filename)
				return true, nil
			}

			if !config.Count {
				PrintMatch(line, idx+1, byteOffset, filename, config, lines, idx)
			}
		}

		byteOffset += int64(len(line)) + 1
	}

	if config.Count && !config.Quiet {
		if config.WithFilename || (!config.NoFilename && config.ShowFilenameForMulti) {
			fmt.Printf("%s:", filename)
		}
		fmt.Println(totalMatches)
	}

	if config.FilesWithoutMatch && !found {
		fmt.Println(filename)
	}

	return found, nil
}

func PrintMatch(line []byte, lineNum int, offset int64, filename string, config *Config, allLines [][]byte, idx int) {
	if filename != "(standard input)" && !config.NoFilename {
		if config.WithFilename || config.ShowFilenameForMulti {
			if config.Color {
				fmt.Printf("%s%s%s:", colorPurple, filename, colorReset)
			} else {
				fmt.Printf("%s:", filename)
			}
		}
	}

	if config.LineNumber {
		if config.Color {
			fmt.Printf("%s%d%s:", colorGreen, lineNum, colorReset)
		} else {
			fmt.Printf("%d:", lineNum)
		}
	}

	if config.ByteOffset {
		fmt.Printf("%d:", offset)
	}

	if config.OnlyMatching && !config.InvertMatch {
		matches := config.Regexp.FindAll(line, -1)
		for _, match := range matches {
			if config.Color {
				fmt.Printf("%s%s%s\n", colorRed, match, colorReset)
			} else {
				fmt.Printf("%s\n", match)
			}
		}
	} else {
		if config.Color && !config.InvertMatch {
			PrintColored(line, config)
		} else {
			fmt.Printf("%s\n", line)
		}
	}
}

func PrintColored(line []byte, config *Config) {
	matches := config.Regexp.FindAllIndex(line, -1)
	if len(matches) == 0 {
		fmt.Printf("%s\n", line)
		return
	}

	last := 0
	for _, match := range matches {
		fmt.Print(string(line[last:match[0]]))
		fmt.Printf("%s%s%s", colorRed, line[match[0]:match[1]], colorReset)
		last = match[1]
	}
	fmt.Printf("%s\n", line[last:])
}

func SearchFilesParallel(files []string, config *Config) bool {
	var wg sync.WaitGroup
	results := make(chan bool, len(files))
	semaphore := make(chan struct{}, 10)

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			found, err := SearchPath(f, config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", f, err)
				results <- false
				return
			}
			results <- found
		}(file)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	found := false
	for result := range results {
		found = found || result
	}

	return found
}
