# About
Flaxton Proxy server is network trafic proxy server for manageing multiple servers TCP trafic. This proxy server written in Go programming language (<a href="http://golang.org" target="_blank">http://golang.org</a>) using <a href="https://github.com/flaxtonio/fntp"  target="_blank">FNTP</a> protocol implementation.<br/>
<b>Flaxton Proxy server is a part of <a href="http://flaxton.io" target="_blank">flaxton.io</a> cloud server load balancer software.</b>

# How it Works
Almost all internet trafic works using TCP. To handle TCP trafic Flaxton Proxy is receiving TCP and converts it to FNTP for better communication between multiple cloud server in load balancing mode, and after excecution Flaxton Proxy getting response from server using FNTP protocol and after sending it back to client as a TCP trafic.
<img src="http://flaxton.io/img/proxyser.gif" />

# "Hello World"
<b>Front Proxy - This proxy will handle TCP requests from clients and will transfer it to other server using FNTP </b>
```go
package main

import (
	"Proxy"
	"fmt"
	"os"
)

func main() {
  proxy := Proxy.CreateFrontProxy(":80")
  proxy.NewClient = "192.168.1.15:8888"
  proxy.Start()
}
```
<b>Server Proxy - This proxy will handle FNTP requests and will proxy it to Web Server application using TCP</b>
```go
package main

import (
	"Proxy"
	"fmt"
	"os"
)

func main() {
	proxy, err := Proxy.CreateServerProxy(":8888", "127.0.0.1:8080")
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
```
<b>Read <a href="https://github.com/flaxtonio/flaxton-proxy/blob/master/tests/proxyTest.go" target="_blank"><code>tests/proxyTest.go</code></a> file for more detailed example</b>
