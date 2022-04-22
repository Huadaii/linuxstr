package model

import (
	"time"
)

// RunTimeInfo GetRunTimeInfo
type RunTimeInfo struct {
	UpTime        float64 `xlsx:"运行时间"`   //运行时间
	LoadAverage   string  `xlsx:"负载情况"`   //负载情况
	TaskTotal     int     `xlsx:"任务总数"`   //任务总数
	RunningTask   int     `xlsx:"运行任务总数"` //运行任务总数
	SleepTask     int     `xlsx:"睡眠进程"`
	StoppedTask   int     `xlsx:"停止进程"`
	IdleTask      int     `xlsx:"闲置进程"`
	ZombieTask    int     `xlsx:"僵尸进程"`  //僵尸进程
	CPUUsePercent float64 `xlsx:"CPU占用"` //使用CPU信息
	MEMFree       uint64  `xlsx:"内存空闲"`
	SwapFree      uint64  `xlsx:"SwapFree"`
	TaskList      []Task  `xlsx:"TaskList"`
}

//top运行任务
type Task struct {
	PID        int32     `xlsx:"PID"`
	Name       string    `xlsx:"进程名称"`
	User       string    `xlsx:"进程拥有者"`
	CPUPercent float64   `xlsx:"CPU占比"` //CPU占比
	MEMPercent float32   `xlsx:"内存占比"`  //内存占比
	RunTime    int       `xlsx:"运行时长"`  //运行时长
	Command    string    `xlsx:"运行命令"`  //运行命令
	CreateTime time.Time `xlsx:"创建时间"`
}
