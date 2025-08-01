package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6378")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Failed to close listener")
		}
	}(l)

	for {
		conn, err := l.Accept()

		if err != nil {

			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func() {
			err := HandleClient(conn)
			if err != nil {
				fmt.Println("Error handling client: ", err.Error())
			}
		}()

	}

}

func HandleClient(conn net.Conn) error {

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection: ", err.Error())
		}
	}(conn)

	for {
		resp := NewResp(conn)
		val, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(val)

		err = resp.Write(val)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
}

// Test the code using this:
// printf '*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n' | nc localhost 6378
