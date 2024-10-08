package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct { }

func (fs FileServer) start() {
	tcp, err := net.Listen("tcp",":3000")
	if err != nil {
		 log.Fatal(err)
	}
	defer tcp.Close()
	for {
		c, err := tcp.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.readLoop(c)
	}
}

func (fs *FileServer)readLoop (c net.Conn){
	packet := new(bytes.Buffer)
	for {
		var size int64
		err := binary.Read(c, binary.LittleEndian, &size)
		if err != nil {
			fmt.Errorf("Not able to read from binary")
		}
		n, err := io.CopyN(packet, c, int64(size))
		if err != nil {
		 log.Fatal(err)
		}
		fmt.Println(packet.Bytes())
		fmt.Println("recieved %d bytes over the network",n)
	}
}

func sendFile(size int) error{
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}
	
	conn, err := net.Dial("tcp",":3000")
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.CopyN(conn, bytes.NewReader(file),int64(size))
	if err != nil {
		return err 
	}
	fmt.Println("written %d bytes over the network", n)
	return nil
}

func main() {
	go func(){
		time.Sleep(3 * time.Second)
		sendFile(20000)
	}()
	server := &FileServer{}
	server.start()
}