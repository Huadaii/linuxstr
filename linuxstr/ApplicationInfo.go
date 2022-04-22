package linuxstr

import (
	"log"
	"os/exec"
	"strings"

	"github.com/Huadaii/linuxstr/model"
	"github.com/lizongshen/gocommand"
)

func GetApplicationInfo() (myApplicationInfos model.ApplicationInfo) {
	myApplicationInfos.ServiceList = GetServiceInfo()
	myApplicationInfos.CronTaskList = GetCronTask()
	myApplicationInfos.UserList = GetAppUser()
	return myApplicationInfos
}

func GetServiceInfo() (myServices []model.Service) {
	output, _ := exec.Command("service", "--status-all").Output()
	outputString := string(output)
	userSlice := strings.Split(outputString, "\n")
	for _, v := range userSlice {
		if v == "" || !strings.Contains(v, "is") {
			continue
		}
		ServiceSlice := strings.Split(v, "is")
		var myService model.Service
		if len(ServiceSlice) < 2 {
			continue
		}
		myService.ServiceName = strings.ReplaceAll(ServiceSlice[0], " ", "")
		if strings.Contains(ServiceSlice[0], "(") {
			myService.ServiceName = strings.Split(strings.ReplaceAll(ServiceSlice[0], " ", ""), "(")[0]
		}
		if strings.Contains(ServiceSlice[1], "running") {
			myService.IsServiceRunning = true
		} else {
			myService.IsServiceRunning = false
		}
		myService.IsAutoRun = getServiceAutoRun(myService.ServiceName)
		myServices = append(myServices, myService)
	}
	return myServices
}

func getServiceAutoRun(ServiceName string) bool {
	_, output, err := gocommand.NewCommand().Exec("chkconfig --list")
	if err != nil {
		log.Println("getServiceAutoRun Error", err)
	}
	for _, val := range strings.Split(output, "\n") {
		result := strings.Fields(val)
		if len(result) == 0 {
			continue
		}
		if result[0] == ServiceName {
			newResutlt := strings.ReplaceAll(result[6], "5:", "")
			if newResutlt == "on" || newResutlt == "启用" {
				return true
			}
		}
	}
	return false
}

func GetCronTask() (myCrons []model.CronTask) {
	output, _ := exec.Command("crontab", "-l").Output()
	outputString := string(output)
	userSlice := strings.Split(outputString, "\n")
	for _, v := range userSlice {
		if strings.HasPrefix(v, "#") || v == "" {
			continue
		}
		CronSplit := strings.Split(v, " ")
		var myCron model.CronTask
		myCron.CronValue = strings.Join(CronSplit[0:4], " ")
		myCron.Shell = strings.Join(CronSplit[5:len(CronSplit)-1], " ")
		myCrons = append(myCrons, myCron)
	}
	return myCrons
}

func GetAppUser() (myUsers []model.AppUser) {
	output, _ := exec.Command("cat", "/etc/passwd").Output()

	outputString := string(output)
	userSlice := strings.Split(outputString, "\n")
	for k, v := range userSlice {
		if k == len(userSlice)-1 {
			break
		}
		userSplit := strings.Split(v, string(':'))
		var myUser model.AppUser
		myUser.Name = userSplit[0]
		myUser.UID = userSplit[2]
		myUser.GID = userSplit[3]
		myUser.UserGroup = getGroupName(myUser.GID)
		myUser.Permission = userSplit[6]
		if userSplit[1] == "x" {
			myUser.EmptyPassword = false
		} else {
			myUser.EmptyPassword = true
		}
		myUsers = append(myUsers, myUser)
	}
	return myUsers
}

func getGroupName(gid string) (groupName string) {
	output, err := exec.Command("cat", "/etc/group", "|", "grep", gid).Output()
	if err != nil {
		return ""
	}
	outputString := string(output)
	groupName = strings.Split(outputString, ":")[0]
	return
}
