package model

//策略信息
type PolicyInfo struct {
	PassWordPolicy PasswordPolicy  //密码策略
	AuditPolicy    map[string]bool //审计策略：检测auditd 和 rsyslog或者syslog是否开启
	IpTable        Iptable         //防火墙策略
	Ufw            Ufw
	LogSizePolicy  []string //保持日志长度
	TimeOutPolicy  int      //超时时间
}

type Iptable struct {
	Status      bool     //防火墙开启状态
	InPutRule   []string //入站策略
	ForwardRule []string //转发策略
	OutPutRule  []string //出战策略
	DenyList    []string //拒绝IP列表
}

type Ufw struct {
	Status string
	Rules  []string
}

type PasswordPolicy struct {
	MaxDay             string //最大时间
	MinDay             string
	MinLen             string             //最小密码长度
	WarnAge            string             //警告时间
	PassWordComplexity PassWordComplexity //密码复杂度
	DenyTime           string             //失败拒绝次数
	LockTime           string             //锁定次数
	UnLockTime         string             //解锁时间
}

type PassWordComplexity struct {
	Difok   string //本次密码与上次密码至少不同字符数
	Minlen  string //密码最小长度，此配置优先于login.defs中的PASS_MAX_DAYS
	Ucredit string //最少大写字母
	Lcredit string //最少小写字母
	Dcredit string //最少数字
	Retry   string //重试多少次后返回密码修改错误
}
