package linuxstr

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"github.com/Huadaii/linuxstr/model"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetHardwareInfo() (h model.HardwareInfo, err error) {
	h.Cpu, err = GetCpuInfo()
	if err != nil {
		return h, err
	}
	h.Disks, err = GetDisksInfo()
	if err != nil {
		return h, err
	}

	h.Memory, err = GetMemoryInfo()
	if err != nil {
		return h, err
	}

	h.Swap, err = GetSwapInfo()
	if err != nil {
		return h, err
	}

	return h, nil
}

func GetCpuInfo() (c model.Cpu, err error) {
	info, err := cpu.Info()
	if err != nil {
		return c, err
	}
	c.ModelName = info[0].ModelName
	c.CPUNum = len(info)
	c.CPUMHZ = info[0].Mhz
	return c, nil
}

func GetDisksInfo() (d []model.Disk, err error) {
	var myDisks []model.Disk
	DefUsageList, err := GetSysMTabInfo()
	if err != nil {
		return []model.Disk{}, err
	}
	for _, v := range DefUsageList {
		var myDisk model.Disk
		myDisk.FileSystem = v.DevName
		myDisk.Size = v.TotalBlock
		myDisk.Used = v.UsedBlock
		myDisk.Avail = v.TotalBlock - v.UsedBlock
		myDisk.UsePercent = v.UsedBlock / v.TotalBlock
		myDisk.Mount = v.MountPoint
		myDisks = append(myDisks, myDisk)
	}
	return myDisks, nil
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// disk usage of path/disk
func DiskUsage(path string) (disk model.DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

func AnalyzeLine2HDDevUsageInfo(lineInfo string) (model.HDDevUsage, error) {

	var HDDev model.HDDevUsage
	arr := strings.Split(lineInfo, " ")
	disk := DiskUsage(arr[1])
	HDDev.DevName = arr[0]
	HDDev.FSName = arr[1]
	HDDev.UsedBlock = float64(disk.Used) / float64(GB)
	HDDev.TotalBlock = float64(disk.All) / float64(GB)
	HDDev.MountPoint = arr[1]
	return HDDev, nil
}

func GetSysMTabInfo() ([]model.HDDevUsage, error) {

	MTabFilePath := "/etc/mtab"

	var HDDevList []model.HDDevUsage
	fp, err := os.Open(MTabFilePath)

	if err != nil {
		return nil, err
	}
	defer fp.Close()

	br := bufio.NewReader(fp)

	for {
		bc, _, errR := br.ReadLine()

		if errR != nil {
			break
		}
		Devinfo, errA := AnalyzeLine2HDDevUsageInfo(string(bc))

		if errA != nil {
			continue
		}

		if !(strings.HasPrefix(Devinfo.DevName, "/dev/vd")) && !(strings.HasPrefix(Devinfo.DevName, "/dev/sd")) {
			continue
		}
		HDDevList = append(HDDevList, Devinfo)
	}

	return HDDevList, nil

}

func GetMemoryInfo() (m model.Memory, err error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return m, err
	}
	m.Total = v.Total
	m.Used = v.Used
	m.Free = v.Free
	m.Shared = v.Shared
	m.Buffers = v.Buffers
	m.Cached = v.Cached
	return m, nil
}

func GetSwapInfo() (s model.Swap, err error) {
	sw, err := mem.SwapMemory()
	if err != nil {
		return s, err
	}
	s.Total = sw.Total
	s.Used = sw.Used
	s.Free = sw.Free

	return s, nil
}
