package model

//系统硬件信息
type HardwareInfo struct {
	Cpu    Cpu    `xlsx:"CPU信息"` //cpu信息
	Memory Memory `xlsx:"内存信息"`  //内存信息
	Disks  []Disk `xlsx:"磁盘信息"`  //磁盘信息
	Swap   Swap   `xlsx:"交换区"`   //交换区
}

type Cpu struct {
	ModelName string  `xlsx:"芯片名"`
	CPUNum    int     `xlsx:"内核数"` //内核数
	CPUMHZ    float64 `xlsx:"主频"`  //主频
}

type Disk struct {
	FileSystem string  `xlsx:"盘符"`    //盘符
	Size       float64 `xlsx:"大小"`    //大小
	Used       float64 `xlsx:"使用大小"`  //使用大小
	Avail      float64 `xlsx:"剩余空间"`  //剩余空间
	UsePercent float64 `xlsx:"使用百分百"` //使用百分百
	Mount      string  `xlsx:"挂载路径"`  //挂载路径
}

type DiskStatus struct {
	All  uint64 `json:"all" xlsx:"All"`
	Used uint64 `json:"used" xlsx:"Used"`
	Free uint64 `json:"free" xlsx:"Free"`
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type HDDevUsage struct {
	DevName    string  `xlsx:"DevName"`
	FSName     string  `xlsx:"FSName"`
	UsedBlock  float64 `xlsx:"UsedBlock"`
	TotalBlock float64 `xlsx:"TotalBlock"`
	MountPoint string  `xlsx:"MountPoint"`
}

type Memory struct {
	Total   uint64 `xlsx:"内存总量"`  //总量
	Used    uint64 `xlsx:"内存使用数"` //使用数
	Free    uint64 `xlsx:"内存空余"`  //空余
	Shared  uint64 `xlsx:"内存分享"`
	Buffers uint64 `xlsx:"内存Buffers"`
	Cached  uint64 `xlsx:"内存缓存"`
}

type Swap struct {
	Total uint64 `xlsx:"内存Total"`
	Used  uint64 `xlsx:"内存Used"`
	Free  uint64 `xlsx:"内存Free"`
}
