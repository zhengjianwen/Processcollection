package funcs

import (
	"os/exec"
	"fmt"
	"bufio"
	"io"
	"strings"
	"syscall"
	"github.com/zhengjianwen/Processcollection/models"

	"strconv"
	"errors"
	"unsafe"
)

func StartWindowscollect() (data []models.Process ) {
	prot_data := windows_prot()
	porcess_data := windows_process()
	for k,v := range porcess_data{
		_,ok := prot_data[k]
		if ok{
			tmp := prot_data[k]
			v.Proto = tmp.Proto
			v.LocalAddr = tmp.LocalAddr
			v.ForeignAddr = tmp.ForeignAddr
			v.State = tmp.State
		}
		data = append(data,v)
	}
	return
}

func windows_prot() (data map[int64]models.Process) {
	cmd := exec.Command("netstat","-ano")
	//显示运行的命令
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("命令错误", err)
		return nil
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	data = make(map[int64]models.Process)
	var status = 0
	for {
		line, err2 := reader.ReadString('\n')
		if status < 4 {
			status += 1
			continue
		}
		if err2 != nil || io.EOF == err2 {
			break
		}
		tmp,errs := makeport(line) //处理信息
		if errs != nil{
			continue
		}
		if _,ok := data[tmp.Pid];ok{
			atpm := data[tmp.Pid]
			atpm.Proto +=  "<br>"+tmp.Proto
			atpm.ForeignAddr +=  "<br>"+tmp.ForeignAddr
			atpm.LocalAddr +=  "<br>"+tmp.LocalAddr
			atpm.State +=  "<br>"+tmp.State
			atpm.Vsz += tmp.Vsz
			tmp = atpm
		}
		data[tmp.Pid] = tmp


	}
	cmd.Wait()
	return
}

func makeport(data string) (newdata models.Process,err error) {
	var tmp = []rune(data) //生成对应的列表
	var status, key = 1, 0 // 状态
	var p_d = ""
	for i := 0; i < len(tmp); i++ {
		if tmp[i] != 32 {
			p_d += string(tmp[i])
			status = 0
		} else {
			if status == 0 {
				switch key {
				case 0:
					newdata.Proto = p_d
				case 1:
					newdata.LocalAddr = p_d
				case 2:
					newdata.ForeignAddr = p_d
				case 3:
					newdata.State = p_d
				}
				p_d = ""
				key += 1
			}
			status = 1
		}
	}

	if key == 4{
		p_d = strings.Replace(p_d,"\r\n","",1)
		pid, err :=strconv.ParseInt(p_d, 10, 64)
		if err != nil{
			return newdata,err
		}
		newdata.Pid = pid
	}else {
		return newdata,errors.New("数据缺少")
	}
	return newdata,nil

}

//处理进程信息
func windows_process() (data map[int64]models.Process) {
	cmd := exec.Command("tasklist")
	//显示运行的命令
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("命令错误", err)
		return nil
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	allmem := GetMemory()
	user := GetUserName()
	data = make(map[int64]models.Process)
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
		tmp,errs := makedatawindows(line) //处理信息
		if errs != nil{
			continue
		}
		tmp.User = user
		tmp.Mem = MemOccupancy(allmem,tmp.Vsz)
		data[tmp.Pid] = tmp
	}
	cmd.Wait()
	return
}
// 处理进程单条信息
func makedatawindows(data string) (newdata models.Process,err error) {
	var tmp = []rune(data) //生成对应的列表
	var status, key = 0, 0 // 状态
	var p_d = ""

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
					pid, err :=strconv.ParseInt(p_d, 10, 64)
					if err != nil{
						return newdata,errors.New("pid转换失败")
					}
					newdata.Pid = pid
				case 2:
					newdata.Tty = p_d
				case 3:
					newdata.Stat = p_d
				case 4:
					p_d = strings.Replace(p_d, ",", "", -1)
					vsz,err :=strconv.ParseInt(p_d, 10, 64)
					if err != nil{
						return newdata,errors.New("vsz转换失败")
					}
					newdata.Vsz = vsz
				}
				status = 1
				key += 1
				p_d = ""
			}
		}
	}

	return newdata,nil

}

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 // in bytes
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var kernel = syscall.NewLazyDLL("Kernel32.dll")

//内存信息
func GetMemory() int64 {
	GlobalMemoryStatusEx := kernel.NewProc("GlobalMemoryStatusEx")
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	mem, _, _ := GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if mem == 0 {
		return 0
	}
	tmp := fmt.Sprint(memInfo.ullTotalPhys)
	mem1,_ := strconv.ParseInt(tmp,10,64)
	return mem1
}

// 用户账号
func GetUserName() string {
	var size uint32 = 128
	var buffer = make([]uint16, size)
	user := syscall.StringToUTF16Ptr("USERNAME")
	r, err := syscall.GetEnvironmentVariable(user, &buffer[0], size)
	if err != nil {
		return ""
	}
	buffer[r] = '@'
	if r >= size {
		return syscall.UTF16ToString(buffer[:r])
	}
	return syscall.UTF16ToString(buffer[:r])
}

func MemOccupancy(allmem,mem int64) float64{
	all:= float64(allmem)/1000
	m := float64(mem)
	ret := m/all
	tmp := fmt.Sprintf("%.2f", ret*100)
	ret ,_ = strconv.ParseFloat(tmp,64)
	return ret
}