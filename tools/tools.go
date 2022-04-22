package tools

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/widuu/goini"
)

/**
工具类实现：
	//TODO
	1 切片处理工具：切片返回列表数组整理
	2 结构体工具：结构体复制，结构体反射
	3 配置文件读取工具：ini，yaml读取工具
*/

/*
时间常量
*/
const (
	//定义Format 日期显示
	FormatDay = "2006-01-02"
	//定义Format 日期时分秒显示
	FormatDayHourMinSecond = "2006-01-02 15:04:05"
)

//Time 转 String
func TimeToString(time time.Time) string {
	return time.Format(FormatDayHourMinSecond)
}

func TimeToStringDay(time time.Time) string {
	return time.Format(FormatDay)
}

//byte 转 int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

//string 转 int
func StringToInt(value string) (i int) {
	i, _ = strconv.Atoi(value)
	return
}

//int 转 int64
func IntToInt64(value int) int64 {
	i, _ := strconv.ParseInt(string(value), 10, 64)
	return i
}

//int64 转 byte
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

//获取字符串在文件内容
func GetStrLine(file, str string) ([]string, error) {
	var result []string
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), str) {
			result = append(result, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

//读取configIni文件内容
func ReadConfigIni(key, value string) (result string) {
	filePath := ""
	filePath = "config/config.ini"
	prefix := goini.SetConfig(filePath)
	result = prefix.GetValue(key, value)
	if result == "no value" {
		log.Println(key, "In config.ini not exist")
		return ""
	}
	return
}

// ScanDirBySuffix 扫描当前目录下文件，不递归扫描,指定后缀
func ScanDirBySuffix(dirName, suffix string) []string {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Println(err)
	}
	var fileList []string
	for _, file := range files {
		if strings.Contains(file.Name(), suffix) {
			fileList = append(fileList, dirName+string(os.PathSeparator)+file.Name())
		}
	}
	return fileList
}

// SingleSlashToDoubleSlash 把系统获取的路径/ 转换为//格式
func SingleSlashToDoubleSlash(path string) string {
	return strings.ReplaceAll(path, "\\", "\\\\")
}

// IsDir 判断目录是否存在
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// SetCommandStd cmd Grep
func SetCommandStd(cmd *exec.Cmd) (stdout, stderr *bytes.Buffer) {
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return
}

// FormatFileSize 字节的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// Cmdexec 执行系统命令
func Cmdexec(cmd string) string {
	var c *exec.Cmd
	var data string
	argArray := strings.Split(cmd, " ")
	c = exec.Command(argArray[0], argArray[1:]...)
	out, _ := c.CombinedOutput()
	data = string(out)
	return data
}

// RemoveDuplicates 去除重复字符串和空格
func RemoveDuplicates(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// 截取小数位数
func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

//数组去重
func DuplicateIntArray(m []int) []int {
	s := make([]int, 0)
	smap := make(map[int]int)
	for _, value := range m {
		//计算map长度
		length := len(smap)
		smap[value] = 1
		//比较map长度, 如果map长度不相等， 说明key不存在
		if len(smap) != length {
			s = append(s, value)
		}
	}
	return s
}

//数组取出不同元素 放入结果
//sourceList中的元素不在sourceList2中 则取到result中
func GetDifferentIntArray(sourceList, sourceList2 []int) (result []int) {
	for _, src := range sourceList {
		var find bool
		for _, target := range sourceList2 {
			if src == target {
				find = true
				continue
			}
		}
		if !find {
			result = append(result, src)
		}
	}
	return
}

//数组存在某个数字
func ExistIntArray(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//字符串数组存在某个字符串
func ExistStringArray(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/*
时间转换函数
*/
func ResolveTime(seconds uint64) string {
	//每分钟秒数
	var day = seconds / (24 * 3600)
	hour := (seconds - day*3600*24) / 3600
	minute := (seconds - day*24*3600 - hour*3600) / 60
	second := seconds - day*24*3600 - hour*3600 - minute*60
	return fmt.Sprintf("%vh%vm%vs", hour, minute, second)
}

//数组取出不同元素 放入结果
//sourceList中的元素不在sourceList2中 则取到result中
func GetDifferentStringArray(sourceList, sourceList2 []string) (result []string) {
	for _, src := range sourceList {
		var find bool
		for _, target := range sourceList2 {
			if src == target {
				find = true
				continue
			}
		}
		if !find {
			result = append(result, src)
		}
	}
	return
}

//字符串数组去重
func DuplicateStringArray(m []string) []string {
	s := make([]string, 0)
	smap := make(map[string]string)
	for _, value := range m {
		//计算map长度
		length := len(smap)
		smap[value] = value
		//比较map长度, 如果map长度不相等， 说明key不存在
		if len(smap) != length {
			s = append(s, value)
		}
	}
	return s
}

//捕获异常 error
func Catch(err error) {
	if err != nil {
		panic(err)
	}
}

//返回某一天的当地时区0点
func ParseMorningTime(t time.Time) time.Time {
	s := t.Format("19931226")
	result, _ := time.ParseInLocation("19931226", s, time.Local)
	return result
}

//返回传入值具体大小
func DisplaySize(raw float64) string {
	if raw < 1024 {
		return fmt.Sprintf("%.1fB", raw)
	}

	if raw < 1024*1024 {
		return fmt.Sprintf("%.1fK", raw/1024.0)
	}

	if raw < 1024*1024*1024 {
		return fmt.Sprintf("%.1fM", raw/1024.0/1024.0)
	}

	if raw < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.1fG", raw/1024.0/1024.0/1024.0)
	}

	if raw < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.1fT", raw/1024.0/1024.0/1024.0/1024.0)
	}

	if raw < 1024*1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.1fP", raw/1024.0/1024.0/1024.0/1024.0/1024.0)
	}

	return "TooLarge"
}

func RemoveElement(slice []interface{}, elem interface{}) []interface{} {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			return slice
		}
	}
	return slice
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
