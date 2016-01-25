package main

import (
	"os"
	"io"
	"fmt"
	"bufio"
)

func checkFileIsExist(filename string) (bool) {
	var exist = true;
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false;
	}
	return exist;
}

func Write(bytes []byte) {
	filename := "./log/redo.txt"
	var f *os.File
	if checkFileIsExist(filename) {  //如果文件存在
		f, err := os.OpenFile(filename, os.O_APPEND | os.O_RDWR, 0777)  //打开文件
		if err != nil {
			return err
		}
	}else {
		f, err := os.Create(filename)  //创建文件
		if err != nil {
			return err
		}
	}
	_,err := io.Write(f, append(bytes,30)) //写入文件 30作为记录分隔符
	if err != nil {
		return err
	}
	return nil
}

func Recover(r func(record []type)) {
	filename := "./log/redo.txt"
	var f *os.File
	if checkFileIsExist(filename) {  //如果文件存在
		f, err := os.Open(filename)  //打开文件
		if err != nil {
			return err
		}
		fmt.Println("开始恢复数据库")
		i := 0
		for r := bufio.NewReader(f);err == nil;buff,err := r.ReadBytes(30) {
			r(buff)
			i++
		}
		fmt.Println("数据库恢复已完成，共",i,"条")
		
	}else {
		return
	}
}