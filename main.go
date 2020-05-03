package main

import "os"

func main() {
	opmlFileName := os.Args[1]
	err := NewArchiver().Run(opmlFileName)
	if err != nil {
		panic(err)
	}
}
