package model

/**
网络信息、具体结构体根据linux架构自行调整
以下结构体参考etc/sysconfig/network-scripts 参数
*/

type NetworkInfo struct {
	NetCards       []NetCard        `json:"网卡状况" xlsx:"网卡信息"` //网卡状况
	ConnectionStat []ConnectionStat `json:"netstat" xlsx:"连接信息"`
}

//网卡信息
type NetCard struct {
	Device  int    `xlsx:"网卡设备"` //网卡设备
	Name    string `xlsx:"网卡名称"`
	Type    string `xlsx:"网络类型"`     //网络类型
	UUID    string `xlsx:"标识符"`      //标识符
	IPAddr  string `xlsx:"IP地址"`     //ip地址
	HWAddr  string `xlsx:"Mac地址"`    //Mac地址
	NetMask string `xlsx:"子网掩码"`     //子网掩码
	GateWay string `xlsx:"网关"`       //网关
	IPType  string `xlsx:"网卡获取ip方式"` //网卡获取ip方式  dhcp ,none ,static
	OnBoote bool   `xlsx:"开机自启动"`    //开机自启动
}

type ConnectionStat struct {
	Fd         uint32  `xlsx:"Fd"`
	Family     uint32  `xlsx:"Family"` // tcp 0x2 udp 0x2 tcp6 0xa unix 0x1
	Type       uint32  `xlsx:"套接字类型"`  //Socket 1TCP 2UDP
	LocalAddr  Addr    `xlsx:"LocalAddr"`
	RemoteAddr Addr    `xlsx:"RemoteAddr"`
	Status     string  `xlsx:"状态"`
	Uids       []int32 `xlsx:"Uids"`
	Pid        int32   `xlsx:"Pid"`
}

type Addr struct {
	Ip   string `xlsx:"Ip"`
	Port int    `xlsx:"Port"`
}
