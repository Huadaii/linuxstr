package model

//应用信息
type ApplicationInfo struct {
	ServiceList  []Service  `xlsx:"服务信息"`
	CronTaskList []CronTask `xlsx:"定时任务信息"`
	UserList     []AppUser  `xlsx:"用户信息"`
}

//服务器信息
type Service struct {
	ServiceName      string `xlsx:"服务名称"`
	IsServiceRunning bool   `xlsx:"运行状态"`
	IsAutoRun        bool   `xlsx:"是否自启动"`
}

//定时任务信息
type CronTask struct {
	CronValue string `xlsx:"定时任务频率"`     //定时任务频率
	Shell     string `xlsx:"执行脚本或者脚本路径"` //执行脚本或者脚本路径
}

type AppUser struct {
	Name          string `xlsx:"应用名称"`
	UID           string `xlsx:"linxu Uid"` //linxu Uid
	GID           string `xlsx:"linxu GID"`
	UserGroup     string `xlsx:"用户组"`                //用户组
	Permission    string `xlsx:"权限说明,运行bash命令,特权用户"` //权限说明，运行bash命令、特权用户
	EmptyPassword bool   `xlsx:"空口令"`                //空口令
}
