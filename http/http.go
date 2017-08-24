package http

import (
	"flag"
	"net/http"

	"fmt"
)

func Start() {
	host := flag.String("host", "0.0.0.0", "listen host")
	port := flag.String("port", "8080", "listen port")

	http.HandleFunc("/", showpress)
	http.HandleFunc("/getdata", getdata)
	fmt.Println("系统启动：", *host, *port)
	err := http.ListenAndServe(*host+":"+*port, nil)

	if err != nil {
		fmt.Println("启动失败！", err)
	}
}
