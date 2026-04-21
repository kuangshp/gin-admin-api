package main

import (
	"fmt"
	"gin-admin-api/test1/test2"
)

func main() {
	repository1 := test2.NewAccountRepository()
	fmt.Println(repository1.GetByName("你死"))
}
