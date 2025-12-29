package ggrep

import (
	"regexp"
)

type Config struct {
	Pattern              string
	Files                []string
	IgnoreCase           bool
	InvertMatch          bool
	Count                bool
	LineNumber           bool
	FilesWithMatches     bool
	FilesWithoutMatch    bool
	NoFilename           bool
	WithFilename         bool
	OnlyMatching         bool
	Quiet                bool
	Recursive            bool
	MaxCount             int
	BeforeContext        int
	AfterContext         int
	FixedStrings         bool
	WordRegexp           bool
	LineRegexp           bool
	Color                bool
	ByteOffset           bool
	ExtendedRegexp       bool
	PerlRegexp           bool
	Regexp               *regexp.Regexp
	ShowFilenameForMulti bool
}

type Match struct {
	Line       []byte
	LineNum    int
	ByteOffset int64
	Matches    [][]int
}

type FileResult struct {
	Filename string
	Matches  []Match
	Found    bool
	Error    error
}
