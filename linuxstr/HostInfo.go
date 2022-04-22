package linuxstr

import (
	"log"

	"github.com/Huadaii/linuxstr/model"
)

func GetHostInfo() (s model.HostInfoService, err error) {
	s.HardwareInfoService, err = GetHardwareInfo()
	s.RunTimeInfoService, err = GetRunTimeInfo()
	s.ApplicationInfoService = GetApplicationInfo()
	s.NetworkInfoService = GetNetworkInfo()
	s.LogInfoService = GetLogInfo(0, "", "", "")
	s.PolicyInfoService, err = GetPolicyInfo()
	if err != nil {
		return s, err
	}
	return s, nil
}

func GetHostInfoFunc(list []string) (s model.HostInfoService, cnName map[string]string, hostInfo []string) {
	var err error
	for _, val := range list {
		switch val {
		case "Hardware":
			s.HardwareInfoService, err = GetHardwareInfo()
			hostInfo = append(hostInfo, "Hardware")
			cnName["Hardware"] = "计算机硬件信息"
			if err != nil {
				log.Println("GetHostInfoFunc GetHardwareInfo Error", err)
			}
		case "RunTime":
			s.RunTimeInfoService, err = GetRunTimeInfo()
			hostInfo = append(hostInfo, "RunTime")
			cnName["Hardware"] = "计算机运行信息"
			if err != nil {
				log.Println("GetHostInfoFunc GetRunTimeInfo Error", err)
			}
		case "Application":
			s.ApplicationInfoService = GetApplicationInfo()
			hostInfo = append(hostInfo, "Application")
			cnName["Hardware"] = "计算机应用信息"
			if err != nil {
				log.Println("GetHostInfoFunc GetApplicationInfo Error", err)
			}
		case "Network":
			s.NetworkInfoService = GetNetworkInfo()
			hostInfo = append(hostInfo, "Network")
			cnName["Hardware"] = "计算机网络信息"
			if err != nil {
				log.Println("GetHostInfoFunc GetNetworkInfo Error", err)
			}
		case "Log":
			s.LogInfoService = GetLogInfo(0, "", "", "")
			hostInfo = append(hostInfo, "Log")
			cnName["Hardware"] = "计算机日志信息"
			if err != nil {
				log.Println("GetHostInfoFunc GetRunTimeInfo Error", err)
			}
		case "Policy":
			s.PolicyInfoService, err = GetPolicyInfo()
			hostInfo = append(hostInfo, "Policy")
			cnName["Hardware"] = "计算机策略信息"
			if err != nil {
				log.Println("GetHostInfoFunc GetPolicyInfo Error", err)
			}
		default:
			s.HardwareInfoService, err = GetHardwareInfo()
			s.RunTimeInfoService, err = GetRunTimeInfo()
			s.ApplicationInfoService = GetApplicationInfo()
			s.NetworkInfoService = GetNetworkInfo()
			s.LogInfoService = GetLogInfo(0, "", "", "")
			s.PolicyInfoService, err = GetPolicyInfo()
			if err != nil {
				log.Println("GetHostInfoFunc GetAll Error", err)
			}
		}
	}
	return s, cnName, hostInfo
}
