package linuxstr

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Huadaii/linuxstr/model"
	"github.com/jinzhu/copier"
	"github.com/lizongshen/gocommand"
	utilnet "github.com/shirou/gopsutil/v3/net"
)

/**
网络信息、具体结构体根据linux架构自行调整
以下结构体参考etc/sysconfig/network-scripts 参数
*/

func GetNetworkInfo() (s model.NetworkInfo) {
	s.NetCards, _ = GetNetCardInfo()
	s.ConnectionStat, _ = GetConnectionStatInfo(0)
	return s
}

func GetNetCardInfo() (NetCardList []model.NetCard, err error) {
	interfaces, err := utilnet.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, val := range interfaces {
		var netCard model.NetCard
		netCard.Name = val.Name
		netCard.Device = val.Index
		uuid, err := GetNetCardUUid(val.Name)
		if err != nil {
			return nil, err
		}
		netCard.UUID = uuid
		if len(val.Addrs) != 0 {
			netCard.IPAddr = strings.Split(val.Addrs[0].Addr, "/")[0]
			IP := net.ParseIP(netCard.IPAddr)
			netCard.NetMask = IP.DefaultMask().String()
		}
		netCard.Type = val.Flags[1]
		netCard.HWAddr = val.HardwareAddr
		netCard.OnBoote = true
		_, gateWay, err := gocommand.NewCommand().Exec(fmt.Sprintf(`ip route show | grep %s | grep default|awk '{print $3}'`, val.Name))
		if err != nil {
			return nil, err
		}
		netCard.GateWay = strings.ReplaceAll(gateWay, "\n", "")
		_, iPType, err := gocommand.NewCommand().Exec(fmt.Sprintf(`ip route show | grep %s | grep default|awk '{print $7}'`, val.Name))
		if err != nil {
			return nil, err
		}
		netCard.IPType = strings.ReplaceAll(iPType, "\n", "")
		NetCardList = append(NetCardList, netCard)
	}
	return NetCardList, err
}

// GetNetCardUUid 获取网卡UUid信息
func GetNetCardUUid(name string) (string, error) {
	cmd := exec.Command("nmcli", "con")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return "", err
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	netCardList := strings.Split(string(bytes), "\n")
	for _, val := range netCardList {
		if !strings.Contains(val, name) {
			continue
		}
		reg := regexp.MustCompile("\\s+")
		rex := reg.ReplaceAllString(strings.TrimSpace(val), " ")
		for _, v := range strings.Split(rex, " ") {
			if len(v) == 36 {
				return v, nil
			}
		}
	}
	return "", nil
}

// GetConnectionStatInfo 获取netstat信息
//max为获取数据限制数量 0=all
func GetConnectionStatInfo(max int) ([]model.ConnectionStat, error) {
	if os.Getenv("CI") != "" {
		return nil, errors.New("skip ci")
	}
	resp, err := utilnet.ConnectionsMax("all", max)
	if err != nil {
		return nil, err
	}
	var connectionStat []model.ConnectionStat
	err = copier.Copy(&connectionStat, &resp)
	if err != nil {
		return nil, err
	}
	return connectionStat, nil
}
