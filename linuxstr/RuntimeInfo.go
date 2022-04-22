package linuxstr

import (
	"context"
	"sort"
	"time"

	"github.com/Huadaii/linuxstr/model"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	ps "github.com/shirou/gopsutil/process"
)

func GetRunTimeInfo() (s model.RunTimeInfo, err error) {
	//Task
	taskInfo, err := GetTaskInfo()
	if err != nil {
		return s, err
	}
	s.TaskList = taskInfo
	//UpTime
	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)
	s.UpTime = time.Now().Sub(t).Seconds()
	//LoadAverage
	avg, err := load.Avg()
	if err != nil {
		return s, err
	}
	s.LoadAverage = avg.String()
	//TaskTotal
	misc, err := load.Misc()
	if err != nil {
		return s, err
	}
	s.RunningTask = misc.ProcsRunning
	//CPUUsePercent
	cpuInfo, err := cpu.Percent(0, false)
	if err != nil || cpuInfo == nil {
		return s, err
	}
	_, S, T, I, Z, _, _, count := GetSpecialStatus()
	s.SleepTask = S
	s.StoppedTask = T
	s.IdleTask = I
	s.ZombieTask = Z
	s.TaskTotal = count + misc.ProcsRunning
	s.CPUUsePercent = cpuInfo[0]
	memory, err := mem.VirtualMemory()
	if err != nil {
		return s, err
	}
	s.MEMFree = memory.Free
	swapMemory, err := mem.SwapMemory()
	if err != nil {
		return s, err
	}
	s.SwapFree = swapMemory.Free
	return s, err
}

var ProcessorStat = make(map[string]int)

// GetSpecialStatus 状态返回进程状态。返回值可能是其中之一。R：运行 S：睡眠 T：停止 I：空闲 Z：僵尸 W：等待 L：锁定
func GetSpecialStatus() (R, S, T, I, Z, W, L, count int) {
	count = ProcessorStat["R"] + ProcessorStat["S"] + ProcessorStat["T"] + ProcessorStat["I"] + ProcessorStat["Z"] + ProcessorStat["E"] + ProcessorStat["L"]
	return ProcessorStat["R"], ProcessorStat["S"], ProcessorStat["T"], ProcessorStat["I"], ProcessorStat["Z"], ProcessorStat["E"], ProcessorStat["L"], count
}

func GetTaskInfo() (taskList []model.Task, err error) {
	var processorStat = make(map[string]int)
	ProcessorStat = processorStat
	processors, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range processors {
		name, _ := p.Name()
		cT, _ := p.CreateTime()
		cpuPercent, _ := p.CPUPercent()
		memPercentage, _ := p.MemoryPercentWithContext(context.Background())
		user, _ := p.Username()
		cmd, _ := p.Cmdline()
		timestamp := time.Unix(0, cT*int64(time.Millisecond))
		var tempTask = model.Task{
			PID:        p.Pid,
			Name:       name,
			User:       user,
			CPUPercent: cpuPercent,
			MEMPercent: memPercentage,
			CreateTime: timestamp,
			Command:    cmd,
		}
		status, err := p.Status()
		if err != nil || status == "" {
			return nil, err
		}
		ProcessorStat[status]++
		taskList = append(taskList, tempTask)
	}
	sort.Slice(taskList, func(i, j int) bool {
		return taskList[i].CPUPercent > taskList[j].CPUPercent
	})
	return taskList[:5], nil
}
