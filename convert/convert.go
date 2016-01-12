package convert

import "unsafe"

func String2C(str string) (unsafe.Pointer) {
	bytes := []byte(str)
	return Bytes2C(bytes)
}

func Bytes2C(bytes []byte) (unsafe.Pointer) {
	if bytes[len(bytes) - 1] != 0 {
		bytes = append(bytes,0)
	}
	return (unsafe.Pointer(&bytes[0]))
}

func Equal(bytes1 []byte,bytes2 []byte) bool {
	Stringfy(&bytes1)
	Stringfy(&bytes2)
	for i := 0;i < len(bytes1);i++ {
		if bytes1[i] == 0 && bytes2[i] == 0 {
			break
		}
		if bytes1[i] == 0 || bytes2[i] == 0 {
			return false
		}
		if bytes1[i] == bytes2[i] {
			continue
		} else {
			return false
		}
	}
	return true
}

func Stringfy(bytes *[]byte) {
	if (*bytes)[len(*bytes) - 1] != 0 {
		*bytes = append(*bytes,0)
	}
}

func ParseUntil(bytes []byte,b byte,start int) ([]byte,int) {
	i := start
	for ;i < len(bytes);i++ {
		if bytes[i] == b {
			break
		}
	}
	return bytes[start:i],i
}

func StartBy(bytes []byte,header string) bool {
	buff := []byte(header)
	for i := 0;i < len(buff);i++ {
		if bytes[i] != header[i] {
			return false
		}
	}
	return true
}

func UpperHead(str string) string {
	bytes := []byte(str)
	if bytes[0] > 'Z' {
		bytes[0] -= ('a' - 'A')
	}
	return string(bytes)
}