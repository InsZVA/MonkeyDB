# MonkeyDB 0.03 Beta 使用说明
## 运行MonkeyDB服务器

> $./server&<br/>
> Server Started!

## 运行brenchmark测试性能

> $./brenchmark<br/>
> 10000 set req within  103177348 ns<br/>
> 10000 get req within  82086509 ns

## 使用Golang库

```
package main
import "./monkey"
import "fmt"

func main() {
	monkeyCli,err := monkey.New("127.0.0.1","1517","monkey")	//服务器地址，端口号（默认1517），密码
	if err != nil {
		panic(err)
	}
	r := monkeyCli.Send([]byte("get a"))
	fmt.Println(string(r))
}
```

## 命令及语法

> set `键` `值`

> get `键`

> remove `键`

> createdb `新数据库名`

> switchdb `数据库名`

> dropdb `数据库名`

> listdb

## 特殊数据库及键值

`monkey`数据库为默认数据库
`config`数据库中`passwd`字段为身份验证密码