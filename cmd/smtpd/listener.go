package main

import (
	"net"
	"log"
	"errors"
)

func listener(l net.Listener, q chan<- net.Conn) {
	for {
		c, err := l.Accept()
		if nil != err {
			log.Printf("failed on accepting connection on [%v]: %v", l, err)
			return
		}
		q <- c
	}
}

func startListeners(listenOnAddress []string) (connections chan net.Conn, listeners []net.Listener, err error) {
	connections = make(chan net.Conn)
	listeners = make([]net.Listener, 0)
	for _, addr := range listenOnAddress {
		if l, err := net.Listen("tcp", addr); nil != err {
			log.Printf("cannot listen on given address [%v]: %v", addr, err)
		} else {
			go listener(l, connections)
			listeners = append(listeners, l)
		}
	}
	if 0 == len(listeners) {
		return nil, nil, errors.New("ERR: require at least one listener")
	}
	return connections, listeners, nil
}
