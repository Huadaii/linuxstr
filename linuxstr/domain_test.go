package linuxstr

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"

	"github.com/lizongshen/gocommand"
)

func TestCpu_GetCpuInfo(t *testing.T) {
	info, err := GetCpuInfo()
	if err != nil {
		return
	}
	fmt.Println(info)
}

func Test_GetSysMTabInfo(t *testing.T) {
	DefUsageList, err := GetSysMTabInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range DefUsageList {
		fmt.Println(v)

	}
}

func TestHardwareInfo_GetHardwareInfo(t *testing.T) {
	info, err := GetHardwareInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info)
}

func TestGetAppUser(t *testing.T) {
	myUsers := GetAppUser()
	fmt.Println(myUsers)
}

func TestCronTask_GetCronTask(t *testing.T) {
	myCrons := GetCronTask()
	fmt.Println(myCrons)
}

func TestService_GetServiceInfo(t *testing.T) {
	myServices := GetServiceInfo()
	fmt.Println(myServices)
}

func Test1test(t *testing.T) {
	command := "cat"
	arg := []string{`/etc/login.defs`}

	output, _ := exec.Command(command, arg...).Output()
	outputString := string(output)
	strings.Index(outputString, "PASS_MAX_DAYS")
	//strings.index
	fmt.Println(string(output))
}

func Test2test(t *testing.T) {
	_, out, err := gocommand.NewCommand().Exec(`grep "Three" /etc/login.defs`)
	if err != nil {
		log.Panic(err)
	}

	log.Println(out)
}

func TestPasswordPolicy_GetPasswordPolicy(t *testing.T) {
	asd, _ := GetPasswordPolicy()
	fmt.Println(asd)
}

func TestIptable_GetIPTableInfo(t *testing.T) {
	info, err := GetIPTableInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info)
}

func TestUfw_GetUfwInfo(t *testing.T) {
	info, err := GetUfwInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range info.Rules {
		fmt.Println(v)
	}
}

func TestPolicyInfo_GetPolicyInfo(t *testing.T) {
	info, err := GetPolicyInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info)
}

func TestTimeOut(t *testing.T) {
	policy, err := getTimeOutPolicy()
	if err != nil {
		return
	}
	fmt.Println(policy)
}

func TestGetLogSizePolicy(t *testing.T) {
	info, err := getLogSizePolicy()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

func TestGetAuditPolicy(t *testing.T) {
	info, err := getAuditPolicy()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}
