package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type HttpPropery struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a buffer to read from the connection
	reader := bufio.NewReader(conn)

	for {

		httpProperty, err := GetHttpProperty(reader)

		// something went wrong while reading
		if err != nil {
			return
		}

		// Check if the request is for the root path
		if httpProperty.Method == "GET" && httpProperty.Path == "/" {
			// Respond with a simple HTTP response
			response := "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Connection: Keep-Alive\r\n" +
				"\r\n" +
				"Hello, World!"

			if httpProperty.Headers["Connection"] == "keep-alive" {

				conn.Write([]byte(response))
				break

			} else {
				conn.Write([]byte(response))
				break

			}

		} else {
			// Respond with a 404 Not Found for other paths
			response := "HTTP/1.1 404 Not Found\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n" +
				"404 - Not Found"

			conn.Write([]byte(response))
		}

	}

}

func GetHttpProperty(reader *bufio.Reader) (*HttpPropery, error) {

	HttpPropery := new(HttpPropery)

	HttpPropery.Headers = make(map[string]string)

	lineCount := 0

	for {

		requestLine, err := reader.ReadString('\n')

		if err != nil {
			return nil, err
		}

		if lineCount == 0 {

			HttpPropery.Method, HttpPropery.Path, HttpPropery.Version, err = parseRequestLine(requestLine)

			if err != nil {

				return nil, err

			}

			lineCount += 1

			continue

		}

		values := strings.SplitN(requestLine, ":", 2)

		if len(values) != 2 {
			break
		}

		key := strings.TrimSpace(values[0])
		value := strings.TrimSpace(values[1])

		// the above things must have atleast 2 values first one will be key and last one will be value itself.
		HttpPropery.Headers[key] = value

	}

	return HttpPropery, nil

}

func parseRequestLine(line string) (method, path, version string, err error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid request line")
	}
	return parts[0], parts[1], parts[2], nil
}

func main() {
	// Listen on TCP port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server is listening on port 8080...")

	// Accept incoming connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle each connection in a separate goroutine
		go handleConnection(conn)
	}
}
