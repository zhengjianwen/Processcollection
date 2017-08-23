package models

type ProcessLinux struct {
	Agent   		string
	User    		string
	Pid     		int64
	Cpu     		float64
	Mem     		float64
	Vsz     		int64
	Rss     		int64
	Tty    			string
	Stat    		string
	Start   		string
	Stime    		string
	Command 		string
	Runtime 		int64
	Proto 			string
	Recvq 			string
	Sendq			string
	LocalAddr		string
	ForeignAddr		string
	State			string
	Program_name 	string
}

type PortLiunx struct {
	Proto 			string
	Recvq 			string
	Sendq			string
	LocalAddr		string
	ForeignAddr		string
	State			string
	Pid 			int64
	Program_name 	string
}

type ProcessWindows struct {
	Command     string
	Pid         string
	SessionName string
	Session     string
	Mem         string
}