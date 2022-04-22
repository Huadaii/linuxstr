package model

//服务器档案
type HostInfoService struct {
	HardwareInfoService    HardwareInfo    `xlsx:"硬件信息"`
	RunTimeInfoService     RunTimeInfo     `xlsx:"运行时信息"`
	ApplicationInfoService ApplicationInfo `xlsx:"应用信息"`
	NetworkInfoService     NetworkInfo     `xlsx:"网络信息"`
	LogInfoService         LogInfo         `xlsx:"日志信息"`
	PolicyInfoService      PolicyInfo      `xlsx:"策略信息"`
}
