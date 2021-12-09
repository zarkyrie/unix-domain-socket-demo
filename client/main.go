package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const (
	tmpSock   = "/tmp/uds1.sock"
	unixConst = "unix"
)

var semaphore = make(chan bool)

func main() {
	unixSockAddr, _ := net.ResolveUnixAddr(unixConst, tmpSock)
	conn, _ := net.DialUnix(unixConst, nil, unixSockAddr)
	defer func(conn *net.UnixConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	go onMsgReceived(conn)

	b := []byte("client -> time\n")
	_, writeErr := conn.Write(b)
	if writeErr != nil {
		return
	}

	<-semaphore
}

func onMsgReceived(conn *net.UnixConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		fmt.Println(msg)
		if err != nil {
			semaphore <- true
			break
		}
		time.Sleep(time.Second)
		b := []byte(msg)
		_, writeErr := conn.Write(b)
		if writeErr != nil {
			return
		}
	}
}
