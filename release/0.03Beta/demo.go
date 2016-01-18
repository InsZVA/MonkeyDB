package main
import "./monkey"
import "fmt"

func main() {
	monkeyCli,err := monkey.New("127.0.0.1","1517","monkey")
	if err != nil {
		panic(err)
	}
	r := monkeyCli.Send([]byte("get a"))
	fmt.Println(string(r))
}