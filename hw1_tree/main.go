package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func tab(slFlag []bool) string {
	var str string
	for _, val := range slFlag {
		if val {
			str += "│\t"
		} else {
			str += "\t"
		}
	}
	return str
}

func recDir(out io.Writer, dir *os.File, fl bool, slFlag []bool) error {
	file , err := dir.Readdir(0)
	if err != nil {
		return err
	}
	if !fl {
		testFile := []os.FileInfo{}
		for _, files := range file {
			if files.IsDir() {
				testFile= append(testFile, files)
			}
		}
		file = testFile
	}
	sort.Slice(file, func(i, j int) bool {return file[i].Name() < file[j].Name()})
	for i, files := range file {
		var str string
		var newTabe []bool
		if i != len(file) - 1 {
			str = tab(slFlag) + "├───" + files.Name()
			newTabe = append(slFlag, true)
		} else {
			str = tab(slFlag) + "└───" + files.Name()
			newTabe = append(slFlag, false)
		}

		if !files.IsDir() {
			if files.Size() != 0 {
				str += " (" + strconv.Itoa(int(files.Size())) + "b)"
			} else {
				str += " (empty)"
			}
		}
		fmt.Fprintf(out, "%s\n", str)
		if files.IsDir() {
			str += files.Name()
			dirLevel, err := os.Open(dir.Name() + "/" + files.Name())
			if err != nil {
				return err
			}
			err = recDir(out, dirLevel, fl, newTabe)
		}
	}
	return err
}

func dirTree (out io.Writer, path string, fl bool) error {
	dir , err := os.Open(path)
	if err != nil {
		return err
	}
	return recDir(out, dir, fl, []bool{})
}

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
