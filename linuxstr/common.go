package linuxstr

import (
	"strings"

	"github.com/Huadaii/linuxstr/tools"
)

var (
	LinuxSystem string = "centos"
)

func SetLinuxSystem(s string) {
	if s == "redhat" {
		LinuxSystem = "redhat"
	}
}

//SetLinuxStr 每次配置本项目版本
func SetLinuxStr() {
	if strings.Contains(tools.Cmdexec("cat /proc/version"), "redhat") || strings.Contains(tools.Cmdexec("cat /proc/version"), "Red Hat") {
		SetLinuxSystem("redhat")
	}
}
