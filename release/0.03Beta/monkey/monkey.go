package monkey

import "./tcp"
import "net"
import "./convert"
import "errors"

type MonkeyCli struct {
	tcpSession *tcp.TCPSession
}

func New(ipAddr string,port string,passwd string) (*MonkeyCli,error){
	tcpAddr,err := net.ResolveTCPAddr("tcp4",ipAddr + ":" + port)
	if err != nil {
		return nil,err
	}
	conn,err := net.DialTCP("tcp",nil,tcpAddr)
	if err != nil {
		return nil,err
	}
	s := &tcp.TCPSession{Conn:conn}
	monkeyCli := &MonkeyCli{tcpSession:s}
	s.Init()
	s.SendMessage([]byte("auth " + passwd))
	buff := s.ReadMessage()
	if !convert.Equal(buff,[]byte("Auth success")) {
		return nil,errors.New(string(buff))
	}
	return monkeyCli,nil
}

func (this *MonkeyCli) Send(cmd []byte) []byte {
	convert.Stringfy(&cmd)
	this.tcpSession.SendMessage(cmd)
	return this.tcpSession.ReadMessage()
}