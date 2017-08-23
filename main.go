package main

import (
	"github.com/liunxprocess/funcs"
	"fmt"
	"encoding/json"
)



func main() {
	data := funcs.StartLiunxcollect()
	d,_ := json.Marshal(data)
	err :=funcs.Writefile(d)
	if err != nil{
		fmt.Println("写入失败！")
	}else {
		fmt.Println("写入成功")
	}
}
// 启动采集并处理





