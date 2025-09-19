package main

import (
	"fmt"
	"github.com/google/uuid"
	"vpn/src/infra/vpn"
)

func main() {
	//startServer()
	//startClient("foo6")
	//time.Sleep(3 * time.Second)
	uuidStr, err := uuid.Parse("6efc0c5c-4b42-4119-8979-463df786a5f6")
	if err != nil {
		fmt.Println(err)
	}
	err = vpn.NewRepository().Delete(uuidStr)
	//fff, err := vpn.NewRepository().Create(uuid.New())
	//err := searchClient("foo3")
	//lines, err := readAndFilterFile(filePath, "foo2")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(fff)
}
