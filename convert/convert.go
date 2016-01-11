package convert

import "unsafe"

func String2C(str string) (unsafe.Pointer) {
	bytes := []byte(str)
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