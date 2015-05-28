package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "2101", "mainPort of server")
	flag.Parse()

	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	checkError(err)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text+"\n")

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err == nil {
			fmt.Print("Message from server: " + message)
		}

	}
}
