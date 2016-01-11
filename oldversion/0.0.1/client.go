package main
import (
	"net"
	"fmt"
	"./tcp"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:1517")  
    if err != nil {
		panic(err)
	}
    conn, err := net.DialTCP("tcp", nil, tcpAddr)  
    if err != nil {
		panic(err)
	}
	
	for {
		buf1 := ""
		buf2 := ""
		buf3 := ""
		buf := ""
		fmt.Print("monkey>")
		fmt.Scanf("%s",&buf1)
		if buf1 == "set" {
			fmt.Scanf("%s",&buf2)
			fmt.Scanf("%s",&buf3)
			buf = buf1 + " " + buf2 + " " + buf3
		}else if buf1 == "get"{
			fmt.Scanf("%s",&buf2)
			buf = buf1 + " " + buf2
		}else if buf1 == "remove" || buf1 == "delete" {
			fmt.Scanf("%s",&buf2)
			buf = buf1 + " " + buf2
		}else if buf1 == "createdb"{
			fmt.Scanf("%s",&buf2)
			buf = buf1 + " " + buf2
		}else if buf1 == "switchdb"{
			fmt.Scanf("%s",&buf2)
			buf = buf1 + " " + buf2
		}else if buf1 == "dropdb"{
			fmt.Scanf("%s",&buf2)
			buf = buf1 + " " + buf2
		}else if buf1 == "listdb"{
			buf = buf1 + " "
		}else if buf1 == "exit"{
			fmt.Println("Bye!")
			break;
		}else if buf1 == "auth"{
			fmt.Scanf("%s",&buf2)
			buf = buf1 + " " + buf2 + " "
		}
		tcp.Send(conn,[]byte(buf))
		buff := tcp.Receive(conn)
		fmt.Println(string(buff))
	}
}