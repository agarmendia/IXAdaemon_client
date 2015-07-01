package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	//"time"
)

func main() {

	var port string
	c := make(chan struct{})
	ch := make(chan struct{})
	flag.StringVar(&port, "port", "2101", "mainPort of server")
	flag.Parse()

	mainService := ":" + port
	tcpAddr1, err := net.ResolveTCPAddr("tcp", mainService)
	conn, err := net.DialTCP("tcp", nil, tcpAddr1)
	if err != nil {
		fmt.Println("Can't connect to server")
		os.Exit(1)
	}
	interceptSignals(conn)
	bufin := bufio.NewWriter(conn)
	bufout := bufio.NewReader(conn)

	go write_server(*conn, *bufin, c)
	go read_server(*conn, *bufout, ch)
	<-c
	//write_server ended
	<-ch
	//read_server ended

}

func write_server(conn net.TCPConn, bufin bufio.Writer, c chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		str := scanner.Text()
		_, err := bufin.WriteString(string(str))
		if err != nil {
			fmt.Print("1")
			fmt.Println(err)
		}
		_, err = bufin.WriteString("\n")
		if err != nil {
			fmt.Print("2")
			fmt.Println(err)
		}
		err = bufin.Flush()
		if err != nil {
			fmt.Print("3")
			fmt.Println(err)
		}
	}
	bufin.WriteString("[IXAdaemon]EOF\n")
	bufin.Flush()
	c <- struct{}{}
	return
}

func read_server(conn net.TCPConn, bufout bufio.Reader, c chan struct{}) {
	for {
		//time.Sleep(time.Second)
		m, err := bufout.ReadString('\n')
		if err != nil {
			fmt.Print("4")
			fmt.Println(err)
			break
		}
		if m == "[IXAdaemon]EOD\n" {
			break
		}
		fmt.Println(m)
	}
	c <- struct{}{}
	return
}
