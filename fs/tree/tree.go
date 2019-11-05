package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

const (
	itemPrefix     = "├───"
	nestedPrefix   = "\t"
	outerPrefix    = "│"
	lastItemPrefix = "└───"
	emptySuffix    = " (empty)"
)

var gIndent = ""

func dirTree(out io.Writer, root string, printFiles bool) error {
	rootDir, err := os.Open(root)
	defer rootDir.Close()
	if err != nil {
		return fmt.Errorf("Error encountered while opening %s", root)
	}

	files, err := rootDir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("Error encountered while retrieving the contents of %s", root)
	}
	files = filterContents(files, printFiles)
	dirLen := len(files)

	indent := outerPrefix + nestedPrefix
	prefix := itemPrefix
	suffix := ""
	for i, f := range files {
		if i == dirLen-1 {
			prefix = lastItemPrefix
			indent = nestedPrefix
		}
		if f.IsDir() {
			fmt.Fprintln(out, strings.Join([]string{gIndent, prefix, f.Name()}, ""))
			gIndent = gIndent + indent
			dirTree(out,
				strings.Join([]string{root, f.Name()}, string(os.PathSeparator)),
				printFiles)
			gIndent = gIndent[:len(gIndent)-len(indent)]
		} else {
			if printFiles {
				size := f.Size()
				if size == 0 {
					suffix = emptySuffix
				} else {
					suffix = strings.Join([]string{" (", strconv.FormatInt(size, 10), "b)"}, "")
				}
			}
			fmt.Fprintln(out, strings.Join([]string{
				gIndent,
				prefix,
				f.Name(),
				suffix,
			}, ""))
		}
	}
    
	return nil
}

func filterContents(list []os.FileInfo, printFiles bool) []os.FileInfo {
	var formatted []os.FileInfo
	if printFiles {
		formatted = list
	} else {
		formatted = list[:0]
		for _, file := range list {
			if file.IsDir() {
				formatted = append(formatted, file)
			}
		}
	}

	sort.Slice(formatted, func(i, j int) bool {
		return formatted[i].Name() < formatted[j].Name()
	})

	return formatted
}
