package tcp
import "net"
import "fmt"

func ok(bytes []byte) bool {
	return bytes[0] == 111 && bytes[1] == 107 && bytes[2] == 0;
}

func bytes4uint(bytes []byte) uint32 {
	total := uint32(0);	
	for i := 0;i < 4;i++ {
		total <<= 8;
		total += uint32(bytes[i]);
	}
	return total
}

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


type TCPSession struct {
	Conn *net.TCPConn
	ToSend chan interface{}	//要发送的数据
	Received chan interface{}	//接受到的数据
	Closed bool //是否已经关闭
}

func (s *TCPSession) Init() {
	s.ToSend = make(chan interface{})
	s.Received = make(chan interface{})
	go s.Send()
	go s.Recv()
}

func (s *TCPSession) Send() {
	for {
		if s.Closed {
			return
		}
		buf0 := <- s.ToSend	//取出要发送的数据
		buf := buf0.([]byte)
		
		_,err := s.Conn.Write(buf)	//发送掉	
		//fmt.Println("send,",buf)
		if err != nil {
			s.Closed = true
			return
		}
	}
	
}

func (s *TCPSession) Recv() {
	for {
		if s.Closed {
			return
		}
		buf := make([]byte,1024)
		_,err := s.Conn.Read(buf)
		if err != nil {
			s.Closed = true
			return
		}
		s.Received <- buf
		//fmt.Println("read,",buf)
		}
	
}

func (s *TCPSession) SendMessage(bytes []byte) {
	total := len(bytes) / 1024
	if len(bytes) % 1024 != 0 {
		total++
	}
	header := uint32bytes(uint32(total))	//计算条数
	s.ToSend <- header
	//fmt.Println(header)
	for i := 0;i < total-1;i++ {
		buf := bytes[0:1024]	//发送这一段
		bytes = bytes[1024:]
		s.ToSend <- buf
		continue
	}
	//发送最后一段
	if total == 0 {
		return
	}
	buf := bytes[0:]	//发送这一段
	s.ToSend <- buf
}

func (s *TCPSession) ReadMessage() []byte {
	buf0 := <- s.Received
	buf := buf0.([]byte)
	//fmt.Println(buf)
	total := bytes4uint(buf)
	var buff []byte
	if buf[4] != 0 {	//两份报表被合并
		buff = buf[4:]
		total--
	} else {
		buff = []byte{}		
	}

	for i := uint32(0);i < total;i++ {
		buf0 := <- s.Received
		buf := buf0.([]byte)
		buff = append(buff,buf...)
	}
	return buff
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////Duplicated/////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func Send(conn *net.TCPConn,bytes []byte) {	//Duplicated
	fmt.Println("to Send:",bytes)
	total := len(bytes) / 1024
	if len(bytes) % 1024 != 0 {
		total++
	}
	header := uint32bytes(uint32(total))
	buff := make([]byte,1024)
	n,err := conn.Write(header)	//先发送一共有多少段
	fmt.Println("sending:",header)
	fmt.Println(n)
	if err != nil {
		fmt.Println(err.Error());
		conn.Close()
		return
	}
	
	_,err = conn.Read(buff)		//等待客户端回应
	if err != nil {
		fmt.Println(err.Error());
		conn.Close()
		return
	}
	buff = append(buff,0)
	fmt.Println("receive:",buff)
	if ok(buff) {
		for i := 0;i < total-1;i++ {
			buf := bytes[0:1024]	//发送这一段
			bytes = bytes[1024:]
			_,err := conn.Write(buf)
			fmt.Println("sending:",buf)
			if err != nil {
				fmt.Println(err.Error());
				conn.Close()
				return
			}
			_,err = conn.Read(buff)
			if err != nil {
				fmt.Println(err.Error());
				conn.Close()
				return
			}
			buff = append(buff,0)
			fmt.Println("receive:",buff)
			if ok(buff) {
				continue
			}
		}
		//发送最后一段
		if total == 0 {
			return
		}
		buf := bytes[0:]	//发送这一段
		conn.Write(buf)
		fmt.Println("sending:",buf)
		conn.Read(buff)
		fmt.Println("receiving:",buff)
	}
}

func Receive(conn *net.TCPConn) []byte {	//Duplicated
	buff := []byte{}
	buf := make([]byte,4)	//先传输总计条数，每条1024byte以内
	fmt.Println("to Receive")
	n,err := conn.Read(buf)
	fmt.Println("receiving:",buf)
	if err != nil {
		fmt.Println(err.Error());
		conn.Close()
	}
	total := bytes4uint(buf)
	n,err = conn.Write([]byte("ok"))	//回应客户端，让其继续传输
	fmt.Println("sending:ok")
	fmt.Println(n)
	if total == 0 {
		fmt.Println(err.Error());
		conn.Close()
		return []byte{}
	}
	if err != nil {
		fmt.Println(err.Error());
		conn.Close()
	}
	for i := uint32(0);i < total;i++ {
		buf := make([]byte,1024)
		_,err := conn.Read(buf)			//每读取一段，回应一次
		fmt.Println("receiving:",buf)
		if err != nil {
			fmt.Println(err.Error());
			conn.Close()
		}
		_,err = conn.Write([]byte("ok"))
		fmt.Println("sending:ok")
		if err != nil {
			fmt.Println(err.Error());
			conn.Close()
		}
		buff = append(buff,buf...)
	}
	return buff
}