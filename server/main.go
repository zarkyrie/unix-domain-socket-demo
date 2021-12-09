package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	tmpSock   = "/tmp/uds1.sock"
	unixConst = "unix"
)

func main() {
	unixSockAddr, _ := net.ResolveUnixAddr(unixConst, tmpSock)
	unixListener, _ := net.ListenUnix(unixConst, unixSockAddr)
	defer func(unixListener *net.UnixListener) {
		err := unixListener.Close()
		if err != nil {
			log.Fatal("unix domain socket close error")
		}
	}(unixListener)

	for {
		unixConn, err := unixListener.AcceptUnix()
		if err != nil {
			continue
		}
		go unixHandler(unixConn)
	}
}

func unixHandler(conn *net.UnixConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnected :" + ipStr)
		err := conn.Close()
		if err != nil {
			return
		}
	}()
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		fmt.Println(msg)
		res := "server->" + time.Now().String() + "\n"
		b := []byte(res)

		_, writeErr := conn.Write(b)
		if writeErr != nil {
			return
		}
	}
}
