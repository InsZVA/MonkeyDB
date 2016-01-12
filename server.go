package main
// #cgo LDFLAGS: -L ./lib -lmonkeyS
// #include "./lib/core.h"
// #include <stdlib.h>
import "C"
import (
	"unsafe"
	"fmt"
	"net"
	"strings"
	"./tcp"
	"./convert"
	"reflect"
)
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////// command类型 用于解析处理各种数据库命令 //////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
var (
	acceptedCmd = []string{"set","get","delete","remove","createdb","switchdb","dropdb","push","pop","destroy","listdb"}
)

type command []byte

func (cmd command) Set(db **C.Database) []byte {
	response := []byte{}
	key,next := convert.ParseUntil(cmd,' ',4)
	value,_ := convert.ParseUntil(cmd,0,next+1)
	r := C.Set(&(*db).tIndex,(*C.char)(convert.Bytes2C(key)),(convert.Bytes2C(value)))
	for i := 0;;i++ {
		response = append(response,byte(r.msg[i]))
		if response[i] == 0 { break; }
	}
	return response
}

func (cmd command) Get(db **C.Database) []byte {
	response := []byte{}
	key,_ := convert.ParseUntil(cmd,0,4)
	r := C.Get(&(*db).tIndex,(*C.char)(convert.Bytes2C(key)))
	if int(r.code) == 0 {
		for i := 0;;i++ {
			response = append(response,byte(*(*C.char)(unsafe.Pointer((uintptr(r.pData)+uintptr(i))))))
			if response[i] == 0 { break; }
		}
	}else {
	}
	return response
}

func (cmd command) Delete(db **C.Database) []byte {
	return cmd.Remove(db)
}

func (cmd command) Remove(db **C.Database) []byte {
	response := []byte{}
	key,_ := convert.ParseUntil(cmd,0,7)
	r := C.Delete(&(*db).tIndex,(*C.char)(convert.Bytes2C(key)))
	for i := 0;;i++ {
		response = append(response,byte(r.msg[i]))
		if response[i] == 0 { break; }
	}
	return response
}

func (cmd command) Createdb(db **C.Database) []byte {
	response := []byte{}
	key,_ := convert.ParseUntil(cmd,0,9)
	d := C.CreateDB((*C.char)(convert.Bytes2C(key)))
	if d != nil {
		*db = d
		response = []byte("Already exist,switched\n")
	}else {
		response = []byte("Created\n")
	}
	return response
}

func (cmd command) Switchdb(db **C.Database) []byte {
	response := []byte{}
	key,_ := convert.ParseUntil(cmd,0,9)
	d := C.SwitchDB((*C.char)(convert.Bytes2C(key)))
	if d != nil {
		*db = d
		response = []byte("ok\n")
	}else {
		response = []byte("fail\n")
	}
	return response
}

func (cmd command) Dropdb(db **C.Database) []byte {
	response := []byte{}
	key,_ := convert.ParseUntil(cmd,0,7)
	*db = C.DropDB((*C.char)(convert.Bytes2C(key)))
	return response
}

func (cmd command) Listdb(db **C.Database) []byte {
	response := []byte{}
	r := C.ListDB()
	fmt.Println(r)
	for i := 0;i < 1024;i++ {
		b := byte(*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(r))+uintptr(i))))
		response = append(response,b)
		if(b == 0){ break; }
	}
	C.free(unsafe.Pointer(r))
	return response
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func initDB() {	//初始化数据库
	str0 := "monkey"
	C.CreateDB((*C.char)(convert.String2C(str0)))	//创建基础数据库
	str0 = "config"						//配置数据库
	C.CreateDB((*C.char)(convert.String2C(str0)))
	key := "passwd"		//初始密码
	data := "monkey"
	db := C.SwitchDB((*C.char)(convert.String2C(str0)))
	C.Set(&(db.tIndex),(*C.char)(convert.String2C(key)),(convert.String2C(data)))
}

func listen() {
	fmt.Println("Server Started!")
	servicePort := ":1517"
	tcpAddr,err := net.ResolveTCPAddr("tcp4",servicePort)
	if err != nil {
		panic(err)
	}
	l,err := net.ListenTCP("tcp",tcpAddr)	//侦听TCP
	if err != nil {
		panic(err)
	}
	for{
		conn,err := l.AcceptTCP()
		conn.SetKeepAlive(true)
		conn.SetNoDelay(true)
		if err != nil {
			panic(err)
		}
		s := tcp.TCPSession{Conn:conn}
		s.Init()
		go Handler(&s)
	}
}

func main() {
	initDB()
	listen()
}

func auth(s *tcp.TCPSession) bool {
	buff := s.ReadMessage()	
	params := strings.Split(string(buff)," ")
	str0 := "config"
	db := C.SwitchDB((*C.char)(convert.String2C(str0)))
	if params[0] != "auth" {
		s.SendMessage([]byte("Please auth first!"))
		return false
	}
	r := C.Get(&(db.tIndex),(*C.char)(convert.String2C("passwd")))
	if int(r.code) == 0 {
		passwd := []byte{}
		for i := 0;;i++ {
			passwd = append(passwd,byte(*(*C.char)(unsafe.Pointer((uintptr(r.pData)+uintptr(i))))))
			if passwd[i] == 0 {
				break
			}
		}
		if convert.Equal(passwd,[]byte(params[1])) {
			s.SendMessage([]byte("Auth success"))
			return true
		} else {
			s.SendMessage([]byte("Auth fail"))
			return false
		}
	}else {
		s.SendMessage([]byte("Auth success"))
		return true
	}
}

func Handler(s *tcp.TCPSession) {
	for !auth(s){
	}
	str := "monkey"							
	db := C.SwitchDB((*C.char)(convert.String2C(str)))//环境变量-当前数据库
	for {
		if s.Closed {
			return
		}				
		buff := s.ReadMessage()
		// if err != nil {
		// 	conn.Close()
		// 	break
		// }
		if len(buff) == 0 {
			return
		}
		//commands := bytes.Split(buff,[]byte{0})
		//for _,cmd := range commands {
			TranslateMessage2(s,&db,buff)
		//}						//解析消息
	}
	
}

func TranslateMessage(s *tcp.TCPSession,db **C.Database,message []byte) {
	command := string(message)
	params := strings.Split(command," ")
	//fmt.Println(params)
	response := []byte{}
	if params[0] == "set" {
		r := C.Set(&(*db).tIndex,(*C.char)(convert.String2C(params[1])),(convert.String2C(params[2])))
		for i := 0;;i++ {
			response = append(response,byte(r.msg[i]))
			if response[i] == 0 { break; }
		}
	}else if params[0] == "get" {
		r := C.Get(&(*db).tIndex,(*C.char)(convert.String2C(params[1])))
		// for i := 0;;i++ {
		// 	response = append(response,byte(r.msg[i]))
		// 	if response[i] == 0 { break; }
		// }
		if int(r.code) == 0 {
			for i := 0;;i++ {
				response = append(response,byte(*(*C.char)(unsafe.Pointer((uintptr(r.pData)+uintptr(i))))))
				if response[i] == 0 { break; }
			}
		}else {
			// for i := 0;;i++ {
			// response = append(response,byte(r.msg[i]))
			// if response[i] == 0 { break; }
			// }
		}
		
	}else if params[0] == "delete" || params[0] == "remove" {
		r := C.Delete(&(*db).tIndex,(*C.char)(convert.String2C(params[1])))
		for i := 0;;i++ {
			response = append(response,byte(r.msg[i]))
			if response[i] == 0 { break; }
		}
		
	}else if params[0] == "createdb" {
		d := C.CreateDB((*C.char)(convert.String2C(params[1])))
		if d != nil {
			*db = d
			response = []byte("Already exist,switched\n")
		}else {
			response = []byte("Created\n")
		}
	}else if params[0] == "switchdb" {
		d := C.SwitchDB((*C.char)(convert.String2C(params[1])))
		if d != nil {
			*db = d
			response = []byte("ok\n")
		}else {
			response = []byte("fail\n")
		}
	}else if params[0] == "dropdb" {
		*db = C.DropDB((*C.char)(convert.String2C(params[1])))
	}else if strings.EqualFold("listdb",params[0]) {
		r := C.ListDB()
		for i := 0;i < 1024;i++ {
			b := byte(*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(r))+uintptr(i))))
			response = append(response,b)
			if(b == 0){ break; }
		}
		C.free(unsafe.Pointer(r))
	}else {
		//fmt.Println("unkown command:",params[0])
	}
	s.SendMessage(response)
}

func TranslateMessage2(s *tcp.TCPSession,db **C.Database,message []byte) {
	com := command(message)
	response := []byte{}
	for _,cmd := range acceptedCmd {
		if convert.StartBy(message,cmd) {
			result := reflect.ValueOf(com).MethodByName(convert.UpperHead(cmd)).Call([]reflect.Value{reflect.ValueOf(db)})
			response = result[0].Interface().([]byte)
		}
	}
	s.SendMessage(response)
}