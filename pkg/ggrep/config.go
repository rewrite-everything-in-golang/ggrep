package ggrep

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

func ParseArgs() *Config {
	config := &Config{}

	flag.BoolVar(&config.IgnoreCase, "i", false, "Ignore case distinctions")
	flag.BoolVar(&config.IgnoreCase, "ignore-case", false, "Ignore case distinctions")

	flag.BoolVar(&config.InvertMatch, "v", false, "Select non-matching lines")
	flag.BoolVar(&config.InvertMatch, "invert-match", false, "Select non-matching lines")

	flag.BoolVar(&config.Count, "c", false, "Print only a count of matching lines")
	flag.BoolVar(&config.Count, "count", false, "Print only a count of matching lines")

	flag.BoolVar(&config.LineNumber, "n", false, "Print line numbers")
	flag.BoolVar(&config.LineNumber, "line-number", false, "Print line numbers")

	flag.BoolVar(&config.FilesWithMatches, "l", false, "Print only names of files with matches")
	flag.BoolVar(&config.FilesWithMatches, "files-with-matches", false, "Print only names of files with matches")

	flag.BoolVar(&config.FilesWithoutMatch, "L", false, "Print only names of files without matches")
	flag.BoolVar(&config.FilesWithoutMatch, "files-without-match", false, "Print only names of files without matches")

	flag.BoolVar(&config.NoFilename, "h", false, "Suppress file name prefix")
	flag.BoolVar(&config.NoFilename, "no-filename", false, "Suppress file name prefix")

	flag.BoolVar(&config.WithFilename, "H", false, "Print file name for each match")
	flag.BoolVar(&config.WithFilename, "with-filename", false, "Print file name for each match")

	flag.BoolVar(&config.OnlyMatching, "o", false, "Show only matching parts of lines")
	flag.BoolVar(&config.OnlyMatching, "only-matching", false, "Show only matching parts of lines")

	flag.BoolVar(&config.Quiet, "q", false, "Suppress all output")
	flag.BoolVar(&config.Quiet, "quiet", false, "Suppress all output")
	flag.BoolVar(&config.Quiet, "silent", false, "Suppress all output")

	flag.BoolVar(&config.Recursive, "r", false, "Search directories recursively")
	flag.BoolVar(&config.Recursive, "R", false, "Search directories recursively")
	flag.BoolVar(&config.Recursive, "recursive", false, "Search directories recursively")

	flag.BoolVar(&config.FixedStrings, "F", false, "Interpret pattern as fixed strings")
	flag.BoolVar(&config.FixedStrings, "fixed-strings", false, "Interpret pattern as fixed strings")

	flag.BoolVar(&config.WordRegexp, "w", false, "Match whole words only")
	flag.BoolVar(&config.WordRegexp, "word-regexp", false, "Match whole words only")

	flag.BoolVar(&config.LineRegexp, "x", false, "Match whole lines only")
	flag.BoolVar(&config.LineRegexp, "line-regexp", false, "Match whole lines only")

	flag.BoolVar(&config.ByteOffset, "b", false, "Print byte offset of matches")
	flag.BoolVar(&config.ByteOffset, "byte-offset", false, "Print byte offset of matches")

	flag.BoolVar(&config.ExtendedRegexp, "E", false, "Extended regular expressions")
	flag.BoolVar(&config.ExtendedRegexp, "extended-regexp", false, "Extended regular expressions")

	flag.BoolVar(&config.PerlRegexp, "P", false, "Perl-compatible regular expressions")
	flag.BoolVar(&config.PerlRegexp, "perl-regexp", false, "Perl-compatible regular expressions")

	flag.IntVar(&config.MaxCount, "m", 0, "Stop after NUM matches")
	flag.IntVar(&config.MaxCount, "max-count", 0, "Stop after NUM matches")

	flag.IntVar(&config.AfterContext, "A", 0, "Print NUM lines after match")
	flag.IntVar(&config.AfterContext, "after-context", 0, "Print NUM lines after match")

	flag.IntVar(&config.BeforeContext, "B", 0, "Print NUM lines before match")
	flag.IntVar(&config.BeforeContext, "before-context", 0, "Print NUM lines before match")

	context := flag.Int("C", 0, "Print NUM lines before and after match")
	contextLong := flag.Int("context", 0, "Print NUM lines before and after match")

	color := flag.Bool("color", false, "Use colors in output")
	colorLong := flag.Bool("colour", false, "Use colors in output")
	noColor := flag.Bool("no-color", false, "Disable colors")
	noColorLong := flag.Bool("no-colour", false, "Disable colors")

	pattern := flag.String("e", "", "Pattern to search for")
	patternLong := flag.String("regexp", "", "Pattern to search for")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: ggrep [OPTIONS] PATTERN [FILE...]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *context > 0 {
		config.BeforeContext = *context
		config.AfterContext = *context
	}
	if *contextLong > 0 {
		config.BeforeContext = *contextLong
		config.AfterContext = *contextLong
	}

	if *color || *colorLong {
		config.Color = true
	}
	if *noColor || *noColorLong {
		config.Color = false
	} else if !(*color || *colorLong) {
		config.Color = IsTerminal()
	}

	if *pattern != "" {
		config.Pattern = *pattern
	} else if *patternLong != "" {
		config.Pattern = *patternLong
	} else if flag.NArg() > 0 {
		config.Pattern = flag.Arg(0)
		config.Files = flag.Args()[1:]
	}

	if len(config.Files) == 0 {
		if flag.NArg() > 1 {
			config.Files = flag.Args()[1:]
		} else {
			config.Files = []string{"-"}
		}
	}

	config.ShowFilenameForMulti = len(config.Files) > 1

	return config
}

func ValidateConfig(config *Config) error {
	if config.Pattern == "" {
		return fmt.Errorf("no pattern specified")
	}
	return nil
}

func CompilePattern(config *Config) error {
	pattern := config.Pattern

	if config.FixedStrings {
		pattern = regexp.QuoteMeta(pattern)
	}

	if config.WordRegexp {
		pattern = `\b` + pattern + `\b`
	}

	if config.LineRegexp {
		pattern = `^` + pattern + `$`
	}

	flags := ""
	if config.IgnoreCase {
		flags += "(?i)"
	}

	pattern = flags + pattern

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex: %v", err)
	}

	config.Regexp = re
	return nil
}

func IsTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
