package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//BasePath 基础路径
var BasePath, _ = os.Getwd()

//Sep 文件分隔符
var Sep string

//APP 应用名称
const APP = "myblog"

//Cache 缓存
type Cache struct {
	BasePath string
}

//GlobalCache 缓存实例
var GlobalCache Cache

//NewCache 新缓存
func NewCache(basepath string) Cache {
	if GlobalCache.BasePath == basepath {
		return GlobalCache
	}
	GlobalCache = Cache{basepath}
	return GlobalCache

}

//GetNowInt 获取当前的时间戳
func GetNowInt() int {
	return int(time.Now().Unix())
}

//PathExists 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//MakeDir 穿件目录 basepath=F:/ addpath=/add/dsd 得到F:/add/dsd
func MakeDir(basepath string, addpath string) string {
	dirlist := strings.Split(addpath, "/")
	templist := make([]string, 0)
	templist = append(templist, basepath+APP)
	os.Mkdir(basepath+APP, 0777)
	// templist = append(templist, APP)
	for _, value := range dirlist {
		if value != "" {
			templist = append(templist, value)
			currentpath := strings.Join(templist, Sep)
			fmt.Println(currentpath, Sep)
			err := os.Mkdir(currentpath, 0777)
			if err != nil {
				panic(err)
			}
		}
	}
	return strings.Join(templist, Sep)
}

//GetPath 获取路径
func GetPath(basepath string, key string) string {
	templist := make([]string, 0)
	templist = append(templist, basepath+APP)
	keylist := strings.Split(key, "/")
	for _, value := range keylist {
		if value != "" {
			templist = append(templist, value)
		}
	}
	return strings.Join(templist, Sep)
}

//Set 设置缓存
func (cache *Cache) Set(key string, expire int, value []byte) error {
	targetdir := GetPath(cache.BasePath, key)
	expiretime := GetNowInt() + expire
	b1, err := PathExists(targetdir)
	if err != nil {
		fmt.Println(err)
	}
	if !b1 {
		fmt.Println("创建目录")
		MakeDir(cache.BasePath, key)
	}
	err = ioutil.WriteFile(targetdir+Sep+strconv.Itoa(expiretime), value, 0777)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

//Get 获取缓存
func (cache *Cache) Get(key string) []byte {
	targetdir := GetPath(cache.BasePath, key)
	b1, _ := PathExists(targetdir)
	if b1 {
		infolist, err := ioutil.ReadDir(targetdir)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		//是否是文件
		if !infolist[len(infolist)-1].IsDir() {
			name := infolist[len(infolist)-1].Name()
			expiretime, err := strconv.Atoi(name)
			if GetNowInt() < expiretime && err == nil {
				content, err := ioutil.ReadFile(targetdir + Sep + name)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				return content
			}

		}

	}
	return nil
}

//LRC 定时清除过期缓存
func (cache *Cache) LRC(key string) error {
	c := time.Tick(20 * time.Second)
	targetdir := GetPath(cache.BasePath, key)
	b1, _ := PathExists(targetdir)
	if !b1 {
		MakeDir(cache.BasePath, key)
	}
	for now := range c {
		fmt.Println(now.Unix())
		infolist, _ := ioutil.ReadDir(targetdir)
		if len(infolist) > 2 {
			targetlist := infolist[0:(len(infolist) - 1)]
			for _, value := range targetlist {
				os.Remove(targetdir + Sep + value.Name())
			}
		}
	}
	return nil
}
func main() {
	if os.IsPathSeparator('\\') {
		Sep = "\\"
	} else {
		Sep = "/"
	}

	cache := NewCache(BasePath + Sep)
	cache.Set("/test2/nihao", 20*3000, []byte("sfsdf2"))
	res := cache.Get("/test2/nihao")
	if res != nil {
		fmt.Println(string(res))
	}
	cache.LRC("/test2/nihao")

}
