package models

type Process struct {
	Agent   		string
	User    		string // 执行用户
	Pid     		int64	// 进程id
	Cpu     		float64	//cpu 使用率
	Mem     		float64 //内存使用率
	Vsz     		int64 // 内存使用大小
	Rss     		int64	//
	Tty    			string	// 当前窗口也就是 tty
	Stat    		string
	Start   		string
	Stime    		string	//运行时长
	Command 		string	// 使用的命令
	Runtime 		int64	// 执行时间
	Proto 			string	// 当前
	Recvq 			string
	Sendq			string
	LocalAddr		string	// 内网监控ip
	ForeignAddr		string	//外网监控ip
	State			string
	Program_name 	string	// 进程名称
}

type PortData struct {
	Proto 			string
	Recvq 			string
	Sendq			string
	LocalAddr		string
	ForeignAddr		string
	State			string
	Pid 			int64
	Program_name 	string
}