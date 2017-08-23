package funcs

import (
	"strconv"
	"os/exec"
	"fmt"
	"bufio"
	"io"
	"errors"
	"strings"
	"github.com/liunxprocess/models"
)

func StartLiunxcollect() (data []models.ProcessLinux ) {
	portdata := liunx_dk()
	processdata := liunx_process()
	for k,v := range processdata{
		_, ok := portdata[k]
		if ok{
			tmp := portdata[k]
			v.Proto = tmp.Proto
			v.Recvq = tmp.Recvq
			v.Sendq = tmp.Sendq
			v.LocalAddr = tmp.LocalAddr
			v.ForeignAddr = tmp.ForeignAddr
			v.State = tmp.State
			v.Program_name = tmp.Program_name
		}
		data = append(data,v)
	}
	return
}


//处理端口信息
func liunx_dk() (data map[int64]models.PortLiunx) {
	cmd := exec.Command("netstat", "-apn")
	//显示运行的命令
	//fmt.Println("运行的命令", cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("命令错误", err)
		return nil
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	data = make(map[int64]models.PortLiunx)
	var status = 2
	for {
		line, err2 := reader.ReadString('\n')
		if status != 0 {
			status -= 1
			continue
		}
		if err2 != nil || io.EOF == err2 {
			break
		}
		tmp, err := makeprotliunx(line)
		if err != nil {
			//fmt.Println(err)
			break
		} else {
			data[tmp.Pid] = tmp
		}

	}
	cmd.Wait()
	return data

}

func makeprotliunx(data string) (models.PortLiunx, error) {
	var tmp = []rune(data) //生成对应的列表
	var status, key = 0, 0 // 状态
	var p_d = ""
	nub := len(tmp)
	newdata := models.PortLiunx{}
	for i := 0; i < nub; i++ {
		if string(tmp[i]) != " " {
			p_d += string(tmp[i])
			status = 0
		} else {
			if status == 0 && key < 6{
				switch key {
				case 0:
					if p_d == "Active"{
						return models.PortLiunx{}, errors.New("执行完毕")
					}
					if p_d == "unix"{
						return models.PortLiunx{}, errors.New("执行完毕")
					}
					newdata.Proto = p_d
				case 1:
					newdata.Recvq = p_d
				case 2:
					newdata.Sendq = p_d
				case 3:
					newdata.LocalAddr = p_d
				case 4:
					newdata.ForeignAddr = p_d
				case 5:
					newdata.State = p_d
				}
				status = 1
				key += 1
				p_d = ""
			} else {
				if key == 6 && string(tmp[i]) != " "{
					p_d += string(tmp[i])
					status = 0
				}else {
					status = 1
				}
			}
		}
	}
	dd := stringsplit(p_d)
	if len(dd) == 2{
		pid,err :=strconv.ParseInt(dd[0], 10, 64)
		if err != nil{
			return models.PortLiunx{},errors.New("分割错误")
		}
		newdata.Pid = pid
		newdata.Program_name = dd[1]
	}

	return newdata, nil
}

func stringsplit(data string) (datalist []string)  {
	rs := []rune(data) //转换成列表
	tmp := ""
	status := 0
	for _,v := range rs{
		t := string(v)
		if t == "/" && status == 0{
			_,err := strconv.ParseInt(tmp,10,64)
			if err == nil{
				datalist = append(datalist,tmp)
			}else {
				return
			}

			tmp = ""
			status += 1
		}else {
			if t != " "{
				tmp += string(v)
			}
		}
	}
	tmp = strings.Replace(tmp, "\n", "", -1)
	datalist = append(datalist,tmp)
	return
}

//处理进程信息

func liunx_process() (data map[int64]models.ProcessLinux) {
	cmd := exec.Command("ps", "aux")
	//显示运行的命令
	//fmt.Println("运行的命令", cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		//fmt.Println("命令错误", err)
		return nil
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	data = make(map[int64]models.ProcessLinux)
	var status = 1
	for {
		line, err2 := reader.ReadString('\n')
		if status == 1 {
			status = 0
			continue
		}
		if err2 != nil || io.EOF == err2 {
			break
		}
		tmp, err := makedataliunx(line)
		if err != nil {
			continue
		} else {
			data[tmp.Pid] = tmp
		}

	}
	//b, _ := json.Marshal(data)

	//err = writefile(b)
	//fmt.Println(data)
	cmd.Wait()
	return
}

func makedataliunx(data string) (models.ProcessLinux, error) {
	var tmp = []rune(data) //生成对应的列表
	var status, key = 0, 0 // 状态
	var p_d = ""
	newdata := models.ProcessLinux{}
	for i := 0; i < len(tmp); i++ {
		if string(tmp[i]) != " " {
			if string(tmp[i]) == "[" || string(tmp[i]) == "]" {
				continue
			} else {
				p_d += string(tmp[i])
			}
			status = 0
		} else {
			if status == 0 && key < 10 {
				switch key {
				case 0:
					newdata.User = p_d
				case 1:
					pid,err :=strconv.ParseInt(p_d, 10, 64)
					if err != nil{
						pid = 0.0
					}
					newdata.Pid = pid
				case 2:
					cpu,err :=strconv.ParseFloat(p_d,  64)
					if err != nil{
						cpu = 0.0
					}
					newdata.Cpu = cpu
				case 3:
					mem,err :=strconv.ParseFloat(p_d,  64)
					if err != nil{
						mem = 0.0
					}
					newdata.Mem = mem
				case 4:
					vsz,err :=strconv.ParseInt(p_d, 10, 64)
					if err != nil{
						vsz = 0.0
					}
					newdata.Vsz = vsz
				case 5:
					rss,err :=strconv.ParseInt(p_d, 10, 64)
					if err != nil{
						rss = 0.0
					}
					newdata.Rss = rss
				case 6:
					newdata.Tty = p_d
				case 7:
					newdata.Stat = p_d
				case 8:
					newdata.Start = p_d
				case 9:
					newdata.Time = p_d
				}
				status = 1
				key += 1
				p_d = ""
			} else {
				if key == 10 {
					p_d += string(tmp[i])
					status = 0
				}

			}
		}
	}

	newdata.Command = strings.Replace(p_d, "\n", "", -1)
	return newdata, nil

}

// From hairui