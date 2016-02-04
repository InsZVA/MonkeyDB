# MonkeyDB - 高性能的内存缓存数据库
## 如何编译

> $ export CC=clang-3.5	//仅支持使用clang-3.5进行编译

> $ go build server.go

> $ go build client.go

> $ go build brenchmark.go

## 如何使用

> $ ./server -h

> Usage of ./server:
  -p string
    	侦听端口号 (default "1517")

> $ ./client -h

> Usage of ./client:
  -r string
    	远程服务器地址 (default "127.0.0.1:1517")

## 命令规范

createdb `dbname`
创建一个新的数据库

switchdb `dbname`
切换到数据库

dropdb `dbname`
删除数据库

listdb
列出所有数据库

auth `password`
获得授权使用数据库（密码字段为config数据库的passwd字段）

set `key` `value`
设置一个键值对（最大长度4M）

get `key`
得到键的值

remove `key`
删除键

命令采用TCP明文未加密传输

## BrenchMark

在256KB L2Cache CPU*8 + DDR3 1600MHz RAM上的测试结果

> $ ./brenchmark

> 100000 set req within  85874753 ns

> 100000 get req within  71206799 ns


每秒处理12w以上请求

## 适用场景

1. 由于未加入日志功能，请勿作为单一数据库使用，建议作为缓存
2. 储存引擎未做遍历优化，不适用于列表型应用

## 声明

本仓库自带monkeyS储存引擎为Ubuntu-256K版本，Windows和Mac和其他缓存大小版本请自行编译

## 协议

本产品采用`GNU General Public License`协议