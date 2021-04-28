package main

import (
	"fmt"
	"os"
)

func main() {
	workDir, _ := os.Getwd()
	fmt.Println(workDir)
}
