package linuxstr

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Huadaii/linuxstr/model"
	"github.com/lizongshen/gocommand"
)

func GetPolicyInfo() (p model.PolicyInfo, err error) {
	p.PassWordPolicy, err = GetPasswordPolicy()
	if err != nil {
		return p, err
	}
	p.Ufw, err = GetUfwInfo()
	if err != nil {
		return p, err
	}
	if p.Ufw.Status == "inactive" {
		p.IpTable, err = GetIPTableInfo()
		if err != nil {
			return p, err
		}
	}

	p.AuditPolicy, err = getAuditPolicy()
	if err != nil {
		return p, err
	}
	p.LogSizePolicy, err = getLogSizePolicy()
	if err != nil {
		return p, err
	}
	p.TimeOutPolicy, err = getTimeOutPolicy()
	if err != nil {
		return p, err
	}
	return p, nil
}

func getTimeOutPolicy() (int, error) {
	//grep "TMOUT" /etc/profile
	_, all, err := gocommand.NewCommand().Exec(`grep -E -o "TMOUT=[0-9]+" /etc/profile`)
	if err != nil {
		return 0, err
	}
	all = strings.Split(all, "\n")[0]
	if all == "" {
		return 0, nil
	}
	result, err := strconv.Atoi(strings.Split(all, "=")[1])
	if err != nil {
		return 0, err
	}
	return result, nil

}

func getLogSizePolicy() ([]string, error) {
	_, all, err := gocommand.NewCommand().Exec(`more /etc/logrotate.conf | grep -v "^#\|^$"`)
	if err != nil {
		return nil, err
	}
	allSlice := strings.Split(all, "\n")
	return allSlice, nil

}

func getAuditPolicy() (map[string]bool, error) {
	//var myPolicy map[string]bool
	myPolicy := make(map[string]bool)
	_, all, err := gocommand.NewCommand().Exec(`service auditd status`)
	if err != nil {
		return nil, err
	}
	if strings.Contains(all, "active") {
		myPolicy["auditd"] = true
	} else {
		myPolicy["auditd"] = false
	}
	_, all, err = gocommand.NewCommand().Exec(`service rsyslog status`)
	if err != nil {
		return nil, err
	}
	if strings.Contains(all, "active") {
		myPolicy["rsyslog"] = true
	} else {
		myPolicy["rsyslog"] = false
	}
	_, all, err = gocommand.NewCommand().Exec(`service syslog status`)
	if err != nil {
		return nil, err
	}
	if strings.Contains(all, "active") {
		myPolicy["syslog"] = true
	} else {
		myPolicy["syslog"] = false
	}
	return myPolicy, nil
}

func GetIPTableInfo() (i model.Iptable, err error) {
	//command := fmt.Sprintf("grep -E -o %s /etc/pam.d/login",KeyName)
	var (
		indexes []int
	)
	_, all, err := gocommand.NewCommand().Exec(`sudo iptables -L`)
	if err != nil {
		return i, err
	}
	allSlice := strings.Split(all, "\n")
	_, out, err := gocommand.NewCommand().Exec(`sudo iptables -L | grep -n "target"`)
	if err != nil {
		return i, err
	}
	outSlice := strings.Split(out, "\n")
	for _, v := range outSlice {
		index, _ := strconv.Atoi(strings.Split(v, ":")[0])
		indexes = append(indexes, index)
	}
	i.InPutRule = allSlice[indexes[0] : indexes[1]-2]
	i.ForwardRule = allSlice[indexes[1] : indexes[2]-2]
	i.OutPutRule = allSlice[indexes[2]:len(allSlice)]
	return i, nil
}

func GetUfwInfo() (u model.Ufw, err error) {
	_, all, err := gocommand.NewCommand().Exec(`sudo ufw status`)
	if err != nil {
		return u, err
	}
	if strings.Contains(all, "inactive") {
		return model.Ufw{
			Status: "inactive",
			Rules:  nil,
		}, nil
	}
	allSlice := strings.Split(all, "\n")
	flag := 0
	for k, v := range allSlice {
		if strings.HasPrefix(v, "--") {
			flag = k + 1
			break
		}
	}
	u.Rules = allSlice[flag : len(allSlice)-1]
	u.Status = "active"
	return u, nil
}

func GetPasswordPolicy() (p model.PasswordPolicy, err error) {
	p.MaxDay, err = getInfoFromLoginDefs("PASS_MAX_DAY")
	if err != nil {
		return p, err
	}
	p.MinLen, err = getInfoFromLoginDefs("PASS_MIN_LEN")
	p.WarnAge, err = getInfoFromLoginDefs("PASS_WARN_AGE")
	p.MinDay, err = getInfoFromLoginDefs("PASS_MIN_DAY")
	p.MinLen, err = getInfoFromLoginDefs("PASS_MIN_LEN")
	p.DenyTime, err = getInfoFromPamd("deny=[0-9]+")
	p.LockTime, err = getInfoFromPamd("lock_time=[0-9]+")
	p.UnLockTime, err = getInfoFromPamd("unlock_time=[0-9]+")
	p.PassWordComplexity.Difok, _ = getInfoFromSysAuth("difok=[0-9]+")
	p.PassWordComplexity.Minlen, _ = getInfoFromSysAuth("minlen=[0-9]+")
	p.PassWordComplexity.Ucredit, _ = getInfoFromSysAuth("ucredit=[-0-9]+")
	p.PassWordComplexity.Lcredit, _ = getInfoFromSysAuth("lcredit=[-0-9]+")
	p.PassWordComplexity.Dcredit, _ = getInfoFromSysAuth("dcredit=[-0-9]+")
	p.PassWordComplexity.Retry, _ = getInfoFromSysAuth("retry=[0-9]+")
	return p, nil
}

func getInfoFromPamd(KeyName string) (value string, err error) {
	KeyName = string('"') + KeyName + string('"')
	command := fmt.Sprintf("grep -E -o %s /etc/pam.d/login", KeyName)
	_, out, err := gocommand.NewCommand().Exec(command)
	if err != nil {
		return "", err
	}
	if out == "" {
		return "nil", nil
	}
	out = strings.Split(out, "\n")[0] //remove the empty line

	value = strings.Split(out, "=")[1] //grep the value
	return value, nil
}

func getInfoFromLoginDefs(KeyName string) (value string, err error) {
	KeyName = string('"') + KeyName + string('"')
	command := fmt.Sprintf("grep %s /etc/login.defs", KeyName)
	_, out, err := gocommand.NewCommand().Exec(command)
	if err != nil {
		return "", err
	}

	outSlice := strings.Split(out, "\n")
	for _, v := range outSlice {
		if strings.HasPrefix(v, "#") == false {
			out = v
			break
		}
	}
	if out == "" {
		return "nil", nil
	}
	value = strings.Fields(out)[1]
	return value, nil
}

func getInfoFromSysAuth(KeyName string) (value string, err error) {
	KeyName = string('"') + KeyName + string('"')
	command := fmt.Sprintf("grep -E -o %s /etc/pam.d/system-auth", KeyName)
	_, out, err := gocommand.NewCommand().Exec(command)
	if err != nil {
		return "", err
	}
	if out == "" {
		return "nil", nil
	}
	out = strings.Split(out, "\n")[0] //remove the empty line

	value = strings.Split(out, "=")[1] //grep the value
	return value, nil
}
