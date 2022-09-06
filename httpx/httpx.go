package httpx

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return i
}

func ParseUriQueryToMap(query string) map[string]interface{} {
	queryMap := strings.Split(query, "&")
	a := make(map[string]interface{}, len(queryMap))
	for _, item := range queryMap {
		itemMap := strings.Split(item, "=")
		a[itemMap[0]] = itemMap[1]
	}
	return a
}

func MapToJson(data map[string]interface{}) string {
	jsonStr, err := json.Marshal(data)
	Must(err)
	return string(jsonStr)
}

// GetDaysAgoZeroTime 获得某一天0点的时间戳
func GetDaysAgoZeroTime(day int) int64 {
	date := time.Now().AddDate(0, 0, day).Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", date)
	return t.Unix()
}

// TimeToHuman 时间戳转人可读
func TimeToHuman(target int) string {
	var res = ""
	if target == 0 {
		return res
	}

	t := int(time.Now().Unix()) - target
	data := [7]map[string]interface{}{
		{"key": 31536000, "value": "年"},
		{"key": 2592000, "value": "个月"},
		{"key": 604800, "value": "星期"},
		{"key": 86400, "value": "天"},
		{"key": 3600, "value": "小时"},
		{"key": 60, "value": "分钟"},
		{"key": 1, "value": "秒"},
	}
	for _, v := range data {
		var c = t / v["key"].(int)
		if 0 != c {
			res = strconv.Itoa(c) + v["value"].(string) + "前"
			break
		}
	}

	return res
}

//获得当前绝对路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// RightAddPathPos 检测并补全路径右边的反斜杠
func RightAddPathPos(path string) string {
	if path[len(path)-1:len(path)] != "/" {
		path = path + "/"
	}
	return path
}

// UriToFilePathByDate url的path转文件名
func UriToFilePathByDate(uriPath string, dir string) string {
	pathArr := strings.Split(uriPath, "/")
	fileName := strings.Join(pathArr, "-")
	writePath := CreateDateDir(dir, "") //根据时间检测是否存在目录，不存在创建
	writePath = RightAddPathPos(writePath)
	fileName = path.Join(writePath, fileName[1:len(fileName)]+".log")
	return fileName
}

// CreateDateDir 根据当前日期，不存在则创建目录
func CreateDateDir(Path string, prex string) string {
	folderName := time.Now().Format("20060102")
	if prex != "" {
		folderName = prex + folderName
	}
	folderPath := filepath.Join(Path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(folderPath, 0755) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0755)
	}
	return folderPath
}
