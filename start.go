package main

import (
	"fmt"
	"vpn/src/infra/vpn"
)

func main() {
	//startServer()
	//startClient("foo6")
	//time.Sleep(3 * time.Second)
	//uuidStr, err := uuid.Parse("7acb0ab6-c7fe-4563-916b-50b3bb8ccee2")
	//if err != nil {
	//	fmt.Println(err)
	//}
	err := vpn.NewRepository().CreateServer()
	//fff, err := vpn.NewRepository().Create(uuid.New())
	//err := searchClient("foo3")
	//lines, err := readAndFilterFile(filePath, "foo2")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(fff)
}
