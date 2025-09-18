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
	test, err := vpn.NewRepository().Create(uuid.New())
	//err := searchClient("foo3")
	//lines, err := readAndFilterFile(filePath, "foo2")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(test)
}
