package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func interceptSignals(conn net.Conn) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGINT:
			fmt.Println("boom")
			conn.Close()
			os.Exit(1)
		case syscall.SIGTERM:
			panic(2)
		}
	}()
}
