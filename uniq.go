package main

import (
	"flag"
	"os"
	"regexp"
	"strings"
)

func getFilenamesIfPassed(opts *Options) {
	args := os.Args[1:]
	str := strings.Join(args, " ")
	pattern := `\b\w+\.txt\b`
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	matches := r.FindAllString(str, -1)
	switch len(matches) {
	case 1:
		opts.inputFile = matches[0]
	case 2:
		opts.inputFile, opts.outputFile = matches[0], matches[1]
	}
}

func initFlags(opts *Options) {
	flag.BoolVar(&opts.getCount, "c", false, "Counts the number of occurrences of the string in the input data.")
	flag.BoolVar(&opts.getRepeated, "d", false, "Outputs only those strings that are repeated in the input data.")
	flag.BoolVar(&opts.getNotRepeated, "u", false, "Outputs only those strings that are not repeated in the input data.")
	flag.IntVar(&opts.skippedFieldsCount, "f", 0, "skips the first num_fields fields in the string.")
	flag.IntVar(&opts.skippedCharsCount, "s", 0, "Skips the first num_chars chars in the string.")
	flag.BoolVar(&opts.ignoreCase, "i", false, "Ignores letter case.")
	flag.Parse()
	getFilenamesIfPassed(opts)
}

func main() {
	var opts Options
	initFlags(&opts)
	uniqueStringsUtility(opts)
}
