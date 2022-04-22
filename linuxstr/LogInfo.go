package linuxstr

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/Huadaii/linuxstr/model"
	"github.com/Huadaii/linuxstr/tools"
)

var ResultCount = 50

func GetLogInfo(num int, hostStr, historyStr, lastStr string) (s model.LogInfo) {
	var (
		ResultNum int
		wg        sync.WaitGroup
	)
	if num != 0 {
		ResultNum = num
	} else {
		ResultNum = ResultCount
	}
	wg.Add(4)
	go func() {
		s.HostLog = GetHostLog(ResultNum, hostStr)
		wg.Done()
	}()
	go func() {
		result, err := GetWebLog()
		if err != nil {
			log.Println("GetWebLog() Error", err)
		}
		s.WebLog = result
		wg.Done()
	}()
	go func() {
		s.HistoryLog = GetHistoryLog(ResultNum, historyStr)
		wg.Done()
	}()
	go func() {
		s.LoginLog = GetLastLog(ResultNum, lastStr)
		wg.Done()
	}()
	wg.Wait()
	return s
}

// GetLastLog 通过命令 last -n 50 --time-format iso -R 获取登录日志
//传入参数为tail数量、grep的字符串
func GetLastLog(num int, str string) (result []model.LoginLog) {
	var logList []string
	if str != "" {
		psCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("last -n %d -R", num))
		psStdout, _ := tools.SetCommandStd(psCmd)
		err := psCmd.Run()
		if err != nil {
			log.Println("GetLastLog Error2", err)
		}
		grepCmd := exec.Command("grep", str)
		grepCmd.Stdin = psStdout
		grepStdout, _ := tools.SetCommandStd(grepCmd)
		_ = grepCmd.Run()
		logList = tools.RemoveDuplicates(strings.Split(grepStdout.String(), "\n"))
	} else {
		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("last -n %d -R", num))
		output, err := cmd.Output()
		if err != nil {
			log.Println("GetLast Error1", err)
		}
		logList = tools.RemoveDuplicates(strings.Split(string(output), "\n"))
	}
	if len(logList) == 0 {
		return result
	}
	for _, v := range logList[0 : len(logList)-1] {
		reg := regexp.MustCompile("\\s+")
		v = reg.ReplaceAllString(strings.TrimSpace(v), " ")
		s := strings.Split(v, " ")
		if len(s) < 4 {
			continue
		}
		name, remote, time, status := GetLoginInfo(s)
		var m model.LoginLog
		m.Username = name
		m.Remote = remote
		if m.Remote == "" {
			continue
		}
		m.Time = time
		m.Status = status
		result = append(result, m)
	}
	return result
}

func GetLoginInfo(list []string) (name, remote, time, status string) {
	name = list[0]
	var index int
	for key, val := range list {
		if key > 4 {
			if len(val+list[key-1]+list[key-2]+list[key-3]) == 13 {
				time = list[key-3] + " " + list[key-2] + list[key-1] + " " + val
				for k := 1; k <= key-4; k++ {
					remote = remote + list[k]
				}
				index = key
				break
			}
		}
	}
	for k := index + 1; k <= len(list)-1; k++ {
		status = status + list[k]
	}
	status = strings.ReplaceAll(status, "-", "")
	return
}

// GetHostLog 通过命令 cat /var/log/secure | tail -n 50 获取Secure
//传入参数为tail数量、grep的字符串
func GetHostLog(num int, str string) (s model.HostLog) {
	if str != "" {
		psCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cat /var/log/secure | tail -n %d", num))
		psStdout, _ := tools.SetCommandStd(psCmd)
		err := psCmd.Run()
		if err != nil {
			log.Println("GetLastLog Error2", err)
		}
		grepCmd := exec.Command("grep", str)
		grepCmd.Stdin = psStdout
		grepStdout, _ := tools.SetCommandStd(grepCmd)
		_ = grepCmd.Run()
		s.HostLogList = tools.RemoveDuplicates(strings.Split(grepStdout.String(), "\n"))
		return s
	}
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cat /var/log/secure | tail -n %d", num))
	output, err := cmd.Output()
	if err != nil {
		log.Println("GetHostLog Error", err)
	}
	s.HostLogList = tools.RemoveDuplicates(strings.Split(string(output), "\n"))
	return s

}

// GetHistoryLog 通过命令 cat ~/.bash_history | tail -n 50 获取历史命令行输入
//传入参数为tail数量、grep的字符串
func GetHistoryLog(num int, str string) (historyLog model.HistoryLog) {
	if str != "" {
		psCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cat ~/.bash_history | tail -n %d", num))
		psStdout, _ := tools.SetCommandStd(psCmd)
		err := psCmd.Run()
		if err != nil {
			log.Println("GetHistoryLog1 Error", err)
		}
		// 筛选
		grepCmd := exec.Command("grep", str)
		grepCmd.Stdin = psStdout
		grepStdout, _ := tools.SetCommandStd(grepCmd)
		_ = grepCmd.Run()
		historyLog.HistoryLogList = tools.RemoveDuplicates(strings.Split(grepStdout.String(), "\n"))
		return
	}
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cat ~/.bash_history | tail -n %d", num))
	// 执行命令，并返回结果
	output, err := cmd.Output()
	if err != nil {
		log.Println("GetHistoryLog3 Error", err)
	}
	historyLog.HistoryLogList = tools.RemoveDuplicates(strings.Split(string(output), "\n"))
	return
}

// GetWebLog Apache Nginx Tomcat
func GetWebLog() (s model.WebLog, err error) {
	//sync
	var wg sync.WaitGroup
	wg.Add(3)
	//Apache
	go func() {
		var webLogInfo model.WebLogInfo
		outApache := tools.Cmdexec("httpd -S")
		regexApache, _ := regexp.Compile(`Main ErrorLog: "(.*?)/error_log"`)
		resultApache := regexApache.FindStringSubmatch(outApache)
		if len(resultApache) >= 2 {
			webLogInfo.Name = "Apache"
			webLogInfo.LogPath = append(webLogInfo.LogPath, resultApache[1])
			webLogInfo.LogList = append(webLogInfo.LogList, tools.ScanDirBySuffix(resultApache[1], ".log")...)
		}
		outApache2 := tools.Cmdexec("httpd -S")
		regexApache2, _ := regexp.Compile(`PidFile: "(.*?)/httpd.pid"`)
		resultApache2 := regexApache2.FindStringSubmatch(outApache2)
		if len(resultApache2) >= 2 {
			webLogInfo.Name = "Apache"
			webLogInfo.LogPath = append(webLogInfo.LogPath, resultApache2[1])
			webLogInfo.LogList = append(webLogInfo.LogList, tools.ScanDirBySuffix(resultApache2[1], ".log")...)
			webLogInfo.LogList = append(webLogInfo.LogList, tools.ScanDirBySuffix(resultApache[1], "_log")...)
		}
		if webLogInfo.LogList != nil {
			s.WebLogInfoList = append(s.WebLogInfoList, webLogInfo)
		}
		wg.Done()
	}()
	//Nginx
	go func() {
		outNginx := tools.Cmdexec("nginx -V")
		regex, _ := regexp.Compile(`\-\-prefix\=(.*?)[ |$]`)
		resultNginx := regex.FindStringSubmatch(outNginx)
		if len(resultNginx) >= 2 {
			var webLogInfo model.WebLogInfo
			webLogInfo.Name = "Nginx"
			configFilePath := resultNginx[1] + "/conf/nginx.conf"
			webLogInfo.LogPath = append(webLogInfo.LogPath, resultNginx[1]+"/logs/")
			dat, err := ioutil.ReadFile(configFilePath)
			if err != nil {
				log.Println(err)
			}
			webLogInfo.LogList = tools.ScanDirBySuffix(resultNginx[1]+"/logs", ".log")
			pathRegex, _ := regexp.Compile(`access_log  (.*?)/access.log\;`)
			pathResult := pathRegex.FindStringSubmatch(string(dat))
			if len(pathResult) >= 2 {
				webLogInfo.LogPath = append(webLogInfo.LogPath, pathResult[1])
				webLogInfo.LogList = append(webLogInfo.LogList, tools.ScanDirBySuffix(pathResult[1], ".log")...)
			}
			s.WebLogInfoList = append(s.WebLogInfoList, webLogInfo)
		}
		wg.Done()
	}()
	//Tomcat
	go func() {
		outTomcat := tools.Cmdexec("find / -name tomcat")
		regex1, _ := regexp.Compile(`/(.*?)tomcat`)
		resultNginx1 := regex1.FindAllString(outTomcat, -1)
		for _, val := range resultNginx1 {
			if tools.IsDir(val + "/logs") {
				var webLogInfo model.WebLogInfo
				webLogInfo.Name = "Tomcat"
				webLogInfo.LogPath = append(webLogInfo.LogPath, val+"/logs")
				webLogInfo.LogList = append(webLogInfo.LogList, tools.ScanDirBySuffix(val+"/logs", ".log")...)
				s.WebLogInfoList = append(s.WebLogInfoList, webLogInfo)
			}
		}
		wg.Done()
	}()
	wg.Wait()
	return s, nil
}
