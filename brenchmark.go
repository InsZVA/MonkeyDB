package main
import (
	"net"
	"fmt"
	"./tcp"
	"time"
	"strconv"
	"runtime"
)

func uint32bytes(n uint32) []byte {
	header := make([]byte,4)
	i := 0
	for n > 0 {
		header[3-i] = byte(n % 256)
		n /= 256
		i++
	}
	return header
}

var	read chan bool


func set() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:1517")  
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)  
    if err != nil {
		panic(err)
	}
	s := tcp.TCPSession{Conn:conn}
	s.Init()
	for i := 1;i < 10000;i++ {
		k := strconv.Itoa(i)
		data := []byte("set " + k + " " + k+k+k+k+k+k+k+k)
		s.SendMessage(data)
		//fmt.Println(string(s.ReadMessage()))
	}
	read <- true
}

func get() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:1517")  
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)  
    if err != nil {
		panic(err)
	}
	s := tcp.TCPSession{Conn:conn}
	s.Init()
	for i := 1;i < 10000;i++ {
		k := strconv.Itoa(i)
		data := []byte("get " + k)
		s.SendMessage(data)
		//fmt.Println(string(s.ReadMessage()))
	}
	read <- true
}

func main() {
	read = make(chan bool,runtime.NumCPU())
	// tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:1517")  
    // if err != nil {
	// 	panic(err)
	// }
    // conn, err := net.DialTCP("tcp", nil, tcpAddr)  
    // if err != nil {
	// 	panic(err)
	// }
	// s := tcp.TCPSession{Conn:conn}
	// s.Init()
	// s.SendMessage([]byte("auth monkey "));
	//fmt.Println(string(s.ReadMessage()))
	start := time.Now().UnixNano()
	// for i := 1;i < 10000;i++ {
	// 	k := strconv.Itoa(i)
	// 	data := []byte("set " + k + " " + k+k+k+k+k+k+k+k)
	// 	s.SendMessage(data)
	// 	//fmt.Println(string(s.ReadMessage()))
	// }
	for i := 0;i < runtime.NumCPU();i++ {
		go set()
	}
	for i := 0;i < runtime.NumCPU();i++ {
		<- read
	}
	fmt.Println("100000 set req within ",time.Now().UnixNano() - start,"ns")
	start = time.Now().UnixNano()
	// for i := 1;i < 10000;i++ {
	// 	k := strconv.Itoa(i)
	// 	data := []byte("get " + k)
	// 	s.SendMessage(data)
	// 	//fmt.Println(string(s.ReadMessage()))
	// }
	for i := 0;i < runtime.NumCPU();i++ {
		go get()
	}
	for i := 0;i < runtime.NumCPU();i++ {
		<- read
	}
	fmt.Println("100000 get req within ",time.Now().UnixNano() - start,"ns")
}