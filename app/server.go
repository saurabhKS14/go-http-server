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
		if data == 0 {
			fmt.Println("No Data Read")
			break
		}
		messages := strings.Split(string(buf), "\r\n")

		// for _, message := range messages {
		fmt.Println(messages[0])
		// call := messages[0]
		request := strings.Split(messages[0], " ")

		fmt.Println(request[0])
		path := request[1]
		// http_version := messages[2]
		if path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else if strings.HasPrefix(path, "/echo") {
			passed_string := strings.SplitN(path, "/", 3)[2]
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(passed_string), passed_string)))
		} else if strings.HasPrefix(path, "/user-agent") {
			user_agent := strings.TrimSpace(strings.Split(messages[2], ":")[1])
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(user_agent), user_agent)))
		} else if strings.HasPrefix(path, "/files/") {
			dir := os.Args[2]
			file_name := strings.Split(path, "/")[2]

			content, err := os.ReadFile(dir + file_name)
			if err != nil {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(content), content)))
			}

		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		// }
	}
}
