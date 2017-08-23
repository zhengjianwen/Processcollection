package funcs

import (
	"os/exec"
	"fmt"
	"bufio"
	"io"
	"strings"
	"github.com/liunxprocess/models"
)


func StartWindowscollect() (data []models.ProcessLinux ) {
	//windows_process := liunx_dk()
	//windows_process := liunx_process()
	//for k,v := range processdata{
	//	_, ok := portdata[k]
	//	if ok{
	//		tmp := portdata[k]
	//		v.Proto = tmp.Proto
	//		v.Recvq = tmp.Recvq
	//		v.Sendq = tmp.Sendq
	//		v.LocalAddr = tmp.LocalAddr
	//		v.ForeignAddr = tmp.ForeignAddr
	//		v.State = tmp.State
	//		v.Program_name = tmp.Program_name
	//	}
	//	data = append(data,v)
	//}
	return
}

func windows_process() (data []models.ProcessWindows) {
	cmd := exec.Command("tasklist")
	//显示运行的命令
	//fmt.Println("运行的命令", cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("命令错误", err)
		return nil
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	var status = 0
	for {
		line, err2 := reader.ReadString('\n')
		if status < 3 {
			status += 1
			continue
		}
		if err2 != nil || io.EOF == err2 {
			break
		}
		tmp, err := makedatawindows(line)
		if err != nil {
			continue
		} else {
			data = append(data, tmp)
		}
	}
	//b, _ := json.Marshal(data)

	//err = writefile(b)
	cmd.Wait()
	return
}

func makedatawindows(data string) (models.ProcessWindows, error) {
	var tmp = []rune(data) //生成对应的列表
	var status, key = 0, 0 // 状态
	var p_d = ""
	newdata := models.ProcessWindows{}
	for i := 0; i < len(tmp); i++ {
		if string(tmp[i]) != " " {
			p_d += string(tmp[i])
			status = 0
		} else {
			if key == 0 && string(tmp[i+1]) != " " {
				continue
			}
			if status == 0 {

				switch key {
				case 0:
					newdata.Command = p_d
				case 1:
					newdata.Pid = p_d
				case 2:
					newdata.SessionName = p_d
				case 3:
					newdata.Session = p_d
				case 4:
					newdata.Mem = strings.Replace(p_d, ",", "", -1)
				}
				status = 1
				key += 1
				p_d = ""
			}
		}
	}

	return newdata, nil

}
