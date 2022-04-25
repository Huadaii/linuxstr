package model

//策略信息
type PolicyInfo struct {
	PassWordPolicy PasswordPolicy  `xlsx:"密码策略"`  //密码策略
	AuditPolicy    map[string]bool `xlsx:"审计策略"`  //审计策略：检测auditd 和 rsyslog或者syslog是否开启
	IpTable        Iptable         `xlsx:"防火墙策略"` //防火墙策略
	Ufw            Ufw             `xlsx:"Ufw"`
	LogSizePolicy  []string        `xlsx:"保持日志长度"` //保持日志长度
	TimeOutPolicy  int             `xlsx:"超时时间"`   //超时时间
}

type Iptable struct {
	Status      bool     `xlsx:"防火墙开启状态"` //防火墙开启状态
	InPutRule   []string `xlsx:"入站策略"`    //入站策略
	ForwardRule []string `xlsx:"转发策略"`    //转发策略
	OutPutRule  []string `xlsx:"出战策略"`    //出战策略
	DenyList    []string `xlsx:"拒绝IP列表"`  //拒绝IP列表
}

type Ufw struct {
	Status string   `xlsx:"状态"`
	Rules  []string `xlsx:"规则"`
}

type PasswordPolicy struct {
	MaxDay             string             `xlsx:"最大时间"` //最大时间
	MinDay             string             `xlsx:"最小时间"`
	MinLen             string             `xlsx:"最小密码长度"` //最小密码长度
	WarnAge            string             `xlsx:"警告时间"`   //警告时间
	PassWordComplexity PassWordComplexity `xlsx:"密码复杂度"`  //密码复杂度
	DenyTime           string             `xlsx:"失败拒绝次数"` //失败拒绝次数
	LockTime           string             `xlsx:"锁定次数"`   //锁定次数
	UnLockTime         string             `xlsx:"解锁时间"`   //解锁时间
}

type PassWordComplexity struct {
	Difok   string `xlsx:"本次密码与上次密码至少不同字符数"` //本次密码与上次密码至少不同字符数
	Minlen  string `xlsx:"密码最小长度"`           //密码最小长度，此配置优先于login.defs中的PASS_MAX_DAYS
	Ucredit string `xlsx:"最少大写字母"`           //最少大写字母
	Lcredit string `xlsx:"最少小写字母"`           //最少小写字母
	Dcredit string `xlsx:"最少数字"`             //最少数字
	Retry   string `xlsx:"重试多少次后返回密码修改错误"`   //重试多少次后返回密码修改错误
}
