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
	var strn string

	//Create channels for sincronizing with go rutines
	c := make(chan bool, 1)
	ch := make(chan bool, 1)

	//Parse Port number, for communicating with IXAdaemon_server
	flag.StringVar(&port, "port", "2101", "mainPort of server")
	flag.Parse()

	//Connect with IXAdaemon_server
	//mainService := ":" + port
	//tcpAddr1, err := net.ResolveTCPAddr("tcp4", mainService)
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Can't connect to server")
		os.Exit(1)
	}

	//Create Writers and Readers for communication
	bufin := bufio.NewWriter(conn)
	bufout := bufio.NewReader(conn)
	sc := bufio.NewScanner(conn)

	interceptSignals(conn)

	//Launch routines
	go write_server(conn, *bufin, c, &strn)
	go read_server(*sc, *bufout, ch)
	//write_server ended
	correctWriting := <-c
	if correctWriting == false {
		//fmt.Println("Failure Writing Server")
	}
	//read_server ended
	correctReading := <-ch
	if correctReading == false {
		//fmt.Println("Failure Reading Server")
	}

}

func write_server(conn net.Conn, bufin bufio.Writer, c chan bool, strn *string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		str := scanner.Text()
		*strn = *strn + str + "\n"
		_, err := bufin.WriteString(string(str))
		if err != nil {
			c <- true
			return
		}
		_, err = bufin.WriteString("\n")
		if err != nil {
			c <- true
			return
		}
		err = bufin.Flush()
		if err != nil {
			c <- true
			return
		}
	}
	bufin.WriteString("[IXAdaemon]EOF\n")
	bufin.Flush()
	c <- false
	return
}

func read_server(scanner bufio.Scanner, bufout bufio.Reader, c chan bool) {
	//canner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		str := scanner.Text()
		if str == "[IXAdaemon]EOD" {
			c <- true
		} else {
			//bufio.NewWriter(os.Stdout).WriteString(str)
			fmt.Println(str)
		}
	}
	c <- true
	return
}
