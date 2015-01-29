package main

import (
	//	"FNTP"
	"Proxy"
	"fmt"
	//	"net"
	"os"
)

func main() {
	if len(os.Args) >= 3 {
		switch os.Args[1] {
		case "front":
			{
				proxy := Proxy.CreateFrontProxy(os.Args[2])
				proxy.NewClient = os.Args[3]
				proxy.Start()
			}
		case "server":
			{
				if len(os.Args) != 4 {
					fmt.Println("Need to Enter Proxy Remote Address")
					return
				}
				proxy, err := Proxy.CreateServerProxy(os.Args[2], os.Args[3])
				if err != nil {
					fmt.Print(err.Error())
					return
				}
				proxy.ErrorHandler = func(err error) {
					fmt.Println(err.Error())
				}
				proxy.Start()
				fmt.Println(proxy)
			}
		}
	}
}
