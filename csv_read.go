package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

var test_phone = make([]string, 0, 10)

func main() {
	//准备读取文件
	fileName := "Master.csv"
	fs, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()

	r := csv.NewReader(fs)
	//针对大文件，一行一行的读取文件
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println(row)
		//test_phone = append(test_phone,row[0])
	}
	//fmt.Println(test_phone)


	//针对小文件，也可以一次性读取所有的文件
	//注意，r要重新赋值，因为readall是读取剩下的
	fs1, _ := os.Open(fileName)
	r1 := csv.NewReader(fs1)
	content, err := r1.ReadAll()
	if err != nil {
		log.Fatalf("can not readall, err is %+v", err)
	}
	for _,row := range content {

		fmt.Println(row)
	}

	//创建一个新文件
	newFileName := "newfile.csv"
	//这样打开，每次都会清空文件内容
	//nfs, err := os.Create(newFileName)

	//这样可以追加写
	nfs, err := os.OpenFile(newFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("can not create file, err is %+v", err)
	}
	defer nfs.Close()
	nfs.Seek(0, io.SeekEnd)

	w := csv.NewWriter(nfs)
	//设置属性
	w.Comma = ','
	w.UseCRLF = true
	row := []string{"1", "2", "3", "4", "5,6"}
	err = w.Write(row)
	if err != nil {
		log.Fatalf("can not write, err is %+v", err)
	}
	//这里必须刷新，才能将数据写入文件。
	w.Flush()

	//一次写入多行
	var newContent [][]string
	newContent = append(newContent, []string{"1", "2", "3", "4", "5", "6"})
	newContent = append(newContent, []string{"11", "12", "13", "14", "15", "16"})
	newContent = append(newContent, []string{"21", "22", "23", "24", "25", "26"})
	w.WriteAll(newContent)

}