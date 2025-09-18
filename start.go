package main

import (
	"fmt"
	"vpn/src/infra/vpn"
)

func main() {
	//startServer()
	//startClient("foo6")
	//time.Sleep(3 * time.Second)
	err := vpn.NewRepository().CreateServer()
	//err := searchClient("foo3")
	//lines, err := readAndFilterFile(filePath, "foo2")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(lines)
}
