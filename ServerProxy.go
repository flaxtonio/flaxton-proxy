package Proxy

import (
	"FNTP"
	"fmt"
	"net"
)

type ServerProxy struct {
	FNTPIncomeCount  int
	ServerSocket     FNTP.Server // FNTP protocol server for getting FNTP traffic
	PassProxyAddress *net.TCPAddr
	ProxyPassCount   int       // Count of Proxy passed connections
	serverChanel     chan bool // This will be used in a future just for not forgetting it
	ErrorHandler     FNTP.ErrorHandler
}

func CreateServerProxy(address, passAddress string) (proxy ServerProxy, err error) {
	proxy.ServerSocket = FNTP.NewServer(address)
	proxy.PassProxyAddress, err = net.ResolveTCPAddr("tcp", passAddress)
	if err != nil {
		return
	}
	proxy.FNTPIncomeCount = 0
	proxy.ProxyPassCount = 0
	return
}

func (proxy *ServerProxy) Start() {
	proxy.ServerSocket.OnNewClient = func(socket *FNTP.Socket) {
		proxy.FNTPIncomeCount++
		remote, err := net.DialTCP("tcp", nil, proxy.PassProxyAddress)
		if err != nil {
			proxy.ErrorHandler(err)
			return
		}
		proxy.ProxyPassCount++
		socket.DataReceived = func(data []byte) {
			fmt.Println(string(data))
			remote.Write(data)
		}
		socket.Disconnected = func(err error) {
			proxy.FNTPIncomeCount--
		}
		data_receive := make([]byte, 1024)
		for {
			len, e := remote.Read(data_receive)
			if e != nil {
				proxy.ProxyPassCount--
				remote.Close()
				socket.Close()
				return
			}
			go socket.Send(data_receive[:len])
		}
	}
	proxy.ServerSocket.ErrorHandling = proxy.ErrorHandler
	proxy.ServerSocket.Listen()
}
