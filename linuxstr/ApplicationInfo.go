package linuxstr

import (
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
	_, outputString, _ := gocommand.NewCommand().Exec(`systemctl | grep  "\.service"`)

	ServicesSlice := strings.Split(outputString, "\n")
	for _, v := range ServicesSlice {
		if v == "" {
			continue
		}
		ServiceSlice := strings.Fields(v)
		var myService model.Service
		myService.ServiceName = ServiceSlice[0]
		if strings.Contains(ServiceSlice[3], "running") {
			myService.IsServiceRunning = true
		} else {
			myService.IsServiceRunning = false
		}
		if strings.Contains(ServiceSlice[2], "active") {
			myService.IsAutoRun = true
		} else {
			myService.IsAutoRun = false
		}
		myServices = append(myServices, myService)
	}
	return myServices
}

func getServiceAutoRun(ServiceName string) bool {
	output, _ := exec.Command("service", ServiceName, "status").Output()
	outputString := string(output)
	if strings.Contains(outputString, "enabled;") {
		return true
	} else {
		return false
	}
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
