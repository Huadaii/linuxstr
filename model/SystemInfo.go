package model

//系统文件
type SystemFiles struct {
	FileName   string
	FileType   string //文件类型：文件夹、可执行文件等
	FileAttr   string //文件属性
	Permission string //文件操作权限，保存类型,777,755之类
	FileUser   string //文件使用者
	FileGroup  string //文件使用组
	FileSize   int    //文件大小
	CreateTime string
}

//主机配置参数信息
type SystemInfo struct {
	Os                  string            `json:"os"`         //系统类型
	PlatformVersion     string            `json:"操作信息版本"`     //操作信息版本
	RegionInfo          string            `json:"地区"`         //地区
	HostName            string            `json:"主机名"`        //主机名
	EnvironmentVariable map[string]string `json:"环境参数列表"`     //环境参数列表
	SystemFiles         []SystemFiles     `json:"系统文件例如host"` //系统文件例如host
	HardwareInfo        HardwareInfo      `json:"硬件信息"`       //硬件信息
}
