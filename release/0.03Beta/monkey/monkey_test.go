package monkey

import "testing"
import "fmt"

func TestNew(t *testing.T) {
	monkey,err := New("127.0.0.1","1517","monkey")
	if err != nil {
		t.Error(err)
	}
	r := monkey.Send([]byte("set a 123"))
	fmt.Println(string(r))
	r = monkey.Send([]byte("get a"))
	fmt.Println(string(r))
}