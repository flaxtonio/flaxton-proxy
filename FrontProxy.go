package Proxy

import (
	"FNTP"
	"fmt"
	"net"
)

//Callback for new TCP connection income, which will allow set server address using string parameter, and make load balance logic
//type FrontConnectionCallback func(*FNTP.Client, *net.TCPConn)

//Proxy That will receive TCP Connections and will send traffic as an FNTP connection to server
type FrontProxy struct {
	TCPIncomeCount     int
	TCPServer          *net.TCPConn
	ServerClientsCount int
	ListenAddress      string
	tcpChanel          chan bool //Chanel for TCP traffic synchronisation ... THIS MAYBE DON'T NEED TO USE
	ErrorHandling      FNTP.ErrorHandler
	NewClient          string // NEED to BE implemented in Load Balance logic
}

func CreateFrontProxy(address string) (frontProxy FrontProxy) {
	frontProxy.ListenAddress = address
	frontProxy.TCPIncomeCount = 0
	frontProxy.ServerClientsCount = 0
	return
}

//Starts TCP Server
func (proxy *FrontProxy) Start() {
	addr, addr_err := net.ResolveTCPAddr("tcp", proxy.ListenAddress)
	if addr_err != nil {
		proxy.ErrorHandling(addr_err)
		return
	}
	socket, err := net.ListenTCP("tcp", addr)
	if err != nil {
		proxy.ErrorHandling(err)
		return
	}
	for {
		conn, packet_err := socket.AcceptTCP()
		if packet_err != nil {
			proxy.ErrorHandling(packet_err)
			continue
		}
		go proxy.HandleNewTCP(conn)
	}
	return
}

func (proxy *FrontProxy) HandleNewTCP(socket *net.TCPConn) {
	proxy.TCPIncomeCount++
	//	var fntp_client FNTP.Client
	var stop_proxy = false
	// With this callback server client will be set
	//	proxy.NewClient(&fntp_client, socket)
	fntp_client := FNTP.NewClient(proxy.NewClient)
	proxy.ServerClientsCount++
	fntp_client.DataReceived = func(data []byte) {
		fmt.Println(string(data))
		socket.Write(data)
	}
	fntp_client.Disconnected = func(err error) {
		stop_proxy = true
		proxy.ServerClientsCount--
	}
	fntp_client.Connect()
	buf_receive := make([]byte, 1024)

	for {
		if stop_proxy {
			socket.Close()
			proxy.TCPIncomeCount--
			fntp_client.Disconnect()
			return
		}
		rlen, err := socket.Read(buf_receive)
		if err != nil {
			socket.Close()
			fntp_client.Disconnect()
			return
		}
		fmt.Println(string(buf_receive[:rlen]))
		go fntp_client.Send(buf_receive[:rlen])
	}
}
