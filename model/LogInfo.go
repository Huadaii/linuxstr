package model

var ResultCount = 50

// LogInfo /**
type LogInfo struct {
	HostLog    HostLog    `xlsx:"主机日志"`
	WebLog     WebLog     `xlsx:"网络日志"`
	HistoryLog HistoryLog `xlsx:"历史日志"` //历史命令日志
	LoginLog   []LoginLog `xlsx:"登录日志"` //登录信息
}

type LoginLog struct {
	Username string `xlsx:"用户名"`
	Remote   string `xlsx:"登录形式"`
	Time     string `xlsx:"登录时间"`
	Status   string `xlsx:"登录状态"`
}

// HostLog 主机Log
type HostLog struct {
	HostLogList []string `xlsx:"主机日志列表"`
}

type HistoryLog struct {
	HistoryLogList []string `xlsx:"历史日志列表"`
}

// WebLog web中间件信息
type WebLog struct {
	WebLogInfoList []WebLogInfo `xlsx:"Web日志信息"`
}

type WebLogInfo struct {
	Name    string   `xlsx:"WEB服务名称"`
	LogPath []string `xlsx:"日志路径"` //日志所在路径
	LogList []string `xlsx:"日志列表"` //日志文件名
}

// OtherLog 其他日志
type OtherLog struct {
	LogType  string
	LogValue []string //存放关键的日志信息
}
