package http

import (
	"net/http"
	//"html/template"
	"github.com/zhengjianwen/Processcollection/funcs"
	"encoding/json"
	"fmt"


	"io/ioutil"
	"os"
)

func showpress(w http.ResponseWriter, req *http.Request)  {
	path := "/tmp/hairui/index.html"

	fin, err := os.Open(path)
	defer fin.Close()
	if err != nil {
		fmt.Println("读取错误")
	}
	fd, _ := ioutil.ReadAll(fin)
	w.Write(fd)

}

func getdata(w http.ResponseWriter, req *http.Request)  {
	data := funcs.StartLiunxcollect()
	bytes, _ := json.Marshal(data)
	fmt.Fprint(w, string(bytes))

}