package http

import (
	"net/http"
	"strings"
	"github.com/zhengjianwen/Processcollection/funcs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func showpress(w http.ResponseWriter, req *http.Request)  {
	path := "index.html"

	fin, err := os.Open(path)
	defer fin.Close()
	if err != nil {
		fmt.Println("读取错误")
	}
	fd, _ := ioutil.ReadAll(fin)
	w.Write(fd)

}

func getdata(w http.ResponseWriter, req *http.Request)  {
	//data := funcs.StartLiunxcollect()
	data := funcs.StartWindowscollect()
	bytes, _ := json.Marshal(data)
	fmt.Fprint(w, string(bytes))

}

func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil{

	}
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}