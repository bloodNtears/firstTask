package uniqueStrings

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	getCount, getRepeated, getNotRepeated, ignoreCase bool
	skippedFieldsCount, skippedCharsCount             int
	inputFile, outputFile                             string
}

func uniqueStringsUtility(opts Options) {
	var lines []string
	inputLines(&lines, opts)
	uniqueLines := getUniqueStrings(lines, opts)
	if opts.getCount {
		countMap := countStringsReps(lines)
		formattedUniqueStrings := writeCountsToString(uniqueLines, countMap)
		outputLines(formattedUniqueStrings, opts)
	}
	if opts.getRepeated {
		countMap := countStringsReps(lines)
		repeatedUniqueLines := getRepeatedStrings(uniqueLines, countMap)
		outputLines(repeatedUniqueLines, opts)
	}
	if opts.getNotRepeated {
		countMap := countStringsReps(lines)
		notRepeatedUniqueLines := getNotRepeatedStrings(uniqueLines, countMap)
		outputLines(notRepeatedUniqueLines, opts)
	}
	if !(opts.getCount || opts.getRepeated || opts.getNotRepeated) {
		outputLines(uniqueLines, opts)
	}
}

func removeFields(line string, n int) string {
	fields := strings.Fields(line)
	if n > len(fields) {
		n = len(fields)
	} else if n == len(fields) {
		n--
	}
	fields = fields[n:]
	return strings.Join(fields, " ")
}

func removeChars(line string, n int) string {
	if len(line) <= n {
		return line
	}
	return line[n:]
}

func getUniqueStrings(lines []string, opts Options) []string {
	uniqueMap := make(map[string]bool)
	uniqueLines := make([]string, 0)
	for _, str := range lines {
		originalLine := str
		if opts.ignoreCase {
			str = strings.ToLower(str)
		}
		if opts.skippedFieldsCount != 0 {
			str = removeFields(str, opts.skippedFieldsCount)
		}
		if opts.skippedCharsCount != 0 {
			str = removeChars(str, opts.skippedCharsCount)
		}
		if !uniqueMap[str] {
			uniqueMap[str] = true
			uniqueLines = append(uniqueLines, originalLine)
		}
	}
	return uniqueLines
}

func doesInputFileExist(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else if os.IsNotExist(err) {
		fmt.Println("Input file", filename, "does not exist. Using stdin to pass data.")
		return false
	} else {
		panic(err)
	}
}

func inputLines(lines *[]string, opts Options) {
	reader := io.Reader(os.Stdin)
	if opts.inputFile != "" {
		if doesInputFileExist(opts.inputFile) {
			file, err := os.Open(opts.inputFile)
			if err != nil {
				panic(err)
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					panic(err)
				}
			}(file)
			reader = file
		}
	}
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		readLines := strings.Split(string(buf[:n]), "\n")

		for i, line := range readLines {
			if line != "" || line == "" && i != len(readLines)-1 {
				*lines = append(*lines, line)
			}
		}

		if err == io.EOF {
			break
		}
	}
}

func createOutputFileIfNotExisted(filename string) {
	if _, err := os.Stat(filename); err == nil {
		return
	} else if os.IsNotExist(err) {
		var file, createErr = os.Create(filename)
		if createErr != nil {
			panic(createErr)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)
		fmt.Printf("Output file %s successfully created.\n", filename)
	}
}

func outputLines(uniqueLines []string, opts Options) {
	var writer io.Writer
	if opts.outputFile != "" {
		createOutputFileIfNotExisted(opts.outputFile)
		file, err := os.Create(opts.outputFile)
		if err != nil {
			panic(err)
		}
		writer = file
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)
	} else {
		writer = os.Stdout
	}
	for _, str := range uniqueLines {
		_, err := fmt.Fprintln(writer, str)
		if err != nil {
			panic(err)
		}
	}
}

func countStringsReps(lines []string) map[string]int {
	countMap := make(map[string]int)
	for _, str := range lines {
		countMap[str]++
	}
	return countMap
}

func writeCountsToString(uniqueStrings []string, countMap map[string]int) []string {
	var formattedStrings []string
	for i, line := range uniqueStrings {
		formattedStrings = append(formattedStrings, fmt.Sprintf("%d %s", countMap[line], uniqueStrings[i]))
	}
	return formattedStrings
}

func getRepeatedStrings(uniqueStrings []string, countMap map[string]int) []string {
	repeatedStringsSlice := make([]string, 0)
	for _, line := range uniqueStrings {
		if countMap[line] > 1 {
			repeatedStringsSlice = append(repeatedStringsSlice, line)
		}
	}
	return repeatedStringsSlice
}

func getNotRepeatedStrings(uniqueStrings []string, countMap map[string]int) []string {
	notRepeatedStringsSlice := make([]string, 0)
	for _, line := range uniqueStrings {
		if countMap[line] < 2 {
			notRepeatedStringsSlice = append(notRepeatedStringsSlice, line)
		}
	}
	return notRepeatedStringsSlice
}
