package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	if listener, err := net.Listen("tcp", "localhost:8088"); err != nil {
		fmt.Println(err)
		return
	} else {

		defer listener.Close()
		fmt.Println("listen start(localhost:8085)")
		for {
			if conn, err := listener.Accept(); err != nil {
				fmt.Println("connection started")
			} else {
				time.Sleep(time.Second)
				conn.Write([]byte("connection started\r\n"))
				// conn.Close()
			}
		}
	}
}
