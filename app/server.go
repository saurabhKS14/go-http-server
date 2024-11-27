package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)

	for {
		data, err := conn.Read(buf)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Reached  end of line")
				break
			} else {
				fmt.Println(err.Error())
			}
		}
		messages := strings.Split(string(data), "\r\n")

		// for _, message := range messages {
		fmt.Println(messages[0])
		// call := messages[0]
		not_headers := strings.Split(messages[0], " ")
		fmt.Println(not_headers[0])
		path := not_headers[1]
		// http_version := messages[2]
		switch path {
		case "/":
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		default:
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		// }
	}
}
