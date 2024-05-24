package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func Echo(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 1024)

	for {
		n, err := conn.Read(buff[0:])
		if err == io.EOF {
			log.Println("End of file")
			break
		}

		if err != nil {
			log.Println("Unexpected error")
			break
		}
		fmt.Printf("Read %d bytes from echo %s\n", n, buff)

		log.Println("Writing data")
		if _, err := conn.Write(buff[0:n]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Listening on port 3000")

	for {
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			os.Exit(1)
		}

		go Echo(conn)
	}
}
