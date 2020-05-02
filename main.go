package main

import "os"

func main() {
	opmlFileName := os.Args[1]
	NewArchiver().Run(opmlFileName)
}
