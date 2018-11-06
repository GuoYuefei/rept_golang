package jsonhandle

import (
	"encoding/json"
	"os"
)

type JsonHandle struct {
	//json文件处理的根目录设置
	JsonFileRoot string
	MyJson
}

type MyJson interface {
	MyJsonSerializer
	MyJsonDeSerializer
}

//序列化接口
type MyJsonSerializer interface {
	StoreJson(a interface{}, storeFile string)
	//这个函数用于将一批对象序列化到同一个文件中，对象与对象之间使用分号“；”隔开
	StoreMyJson(a []interface{}, storeFile string)
}

//反序列化接口
type MyJsonDeSerializer interface {

}

//这里主要用于写正常格式的json格式的序列化
func (j *JsonHandle)StoreJson(a interface{}, storeFile string) {
	flag := false		//如果执行过创建目录一次就置true
	data, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	storeFile = j.JsonFileRoot + storeFile
FLAGCREATE:
	file, err := os.Create(storeFile)
	if err != nil && !flag {
		e := os.MkdirAll(j.JsonFileRoot,666)			//三组用户都有读写权
		flag = true
		if e != nil {
			panic(e)
		}
		goto FLAGCREATE
	}
	defer file.Close()
	if err != nil {
		panic(err)
	}
	file.Write(data)
}

func (j *JsonHandle)StoreMyJson(a []interface{}, storeFile string) {

}
