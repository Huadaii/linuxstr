package model

//主机属性
type SystemProperties struct {
	HostID           string `json:"host_ID" gorm:"host_ID"`                     //主机编号
	RiskValue        string `json:"risk_value" gorm:"risk_value"`               //风险值
	Importance       string `json:"importance" gorm:"importance"`               //重要级
	Business         string `json:"business" gorm:"business"`                   //业务
	Department       string `json:"department" gorm:"department"`               //部门
	AdministratorUID int    `json:"administrator_uid" gorm:"administrator_uid"` //运维管理员
}
