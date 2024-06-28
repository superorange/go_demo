package proxy

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
)

// 定义缓冲池
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 128*1024) // 32KB 缓冲区
	},
}

// forwardClientConnection 处理客户端到服务器的数据转发
func forwardClientConnection(clientConn net.Conn) {
	defer func() {
		if err := clientConn.Close(); err != nil {
			log.Println("Error closing client connection:", err)
		}
	}()

	serverConn, err := net.Dial("tcp", "127.0.0.1:5201")
	if err != nil {
		log.Println("Error connecting to server:", err)
		return
	}
	defer func() {
		if err := serverConn.Close(); err != nil {
			log.Println("Error closing server connection:", err)
		}
	}()

	done := make(chan struct{})
	go func() {
		transfer(serverConn, clientConn)
		done <- struct{}{}
	}()
	go func() {
		transfer(clientConn, serverConn)
		done <- struct{}{}
	}()
	<-done
	<-done
}

// transfer 使用缓冲区高效转发数据
func transfer(dst io.Writer, src io.Reader) {
	//buffer := bufferPool.Get().([]byte)
	//defer bufferPool.Put(buffer)
	buffer := make([]byte, 1024*64)

	_, err := io.CopyBuffer(dst, src, buffer)
	if err != nil {
		log.Println("Error during transfer:", err)
	}
}

// Start 初始化代理服务器以监听指定端口
func Start() {

	os.OpenFile("ope", os.O_RDONLY, 0)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("Error listening on port 8082:", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Println("Error closing listener:", err)
		}
	}()
	log.Println("Listening tcp on 127.0.0.1:8082...")
	for {
		clientConn, err := listener.Accept()
		log.Printf("clientConn: %v\n", clientConn)
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go forwardClientConnection(clientConn)
	}
}
