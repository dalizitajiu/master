package lib

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

//BasePath 基础路径
const BasePath = "F:\\"

//APP 应用名称
const APP = "myblog"

//Cache 缓存
type Cache struct {
	BasePath string
}

//Set 设置缓存
func (cache *Cache) Set(key string, expire int, value []byte) error {
	expiretime := GetNowInt() + expire
	err := ioutil.WriteFile(cache.BasePath+APP+"\\"+key+"\\"+strconv.Itoa(expiretime), value, 0644)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
