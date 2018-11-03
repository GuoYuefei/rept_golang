package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/djimenez/iconv-go"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type DataStruct struct {
	StockCode string
	StockName string
	//开票价
	OpPrice float64
	//昨日收盘价
	Settlement float64
	//当前价格
	CurrentPrice float64
	//今日最高 The Highest today
	THT float64
	//今日最低
	TLT float64
	//当前买一竞买价 Current bid price
	CurBidP float64
	//当前卖一竞卖价
	CurBidingP float64
	//成交股市 Number of transactions,
	NumOfTrans float64
	//成交金额 Transaction amount
	TransAmount float64

	//5手报价和股数暂不分析
}

func NewDataStruct(sc, sn string, fs []float64) (ds *DataStruct) {
	return &DataStruct{
		sc,
		sn,
		fs[0],
		fs[1],
		fs[2],
		fs[3],
		fs[4],
		fs[5],
		fs[6],
		fs[7],
		fs[8],
	}
}

type Mark int

const (
	code Mark = iota
	fileName

)

const dataInterBase  = "http://hq.sinajs.cn/list="

//sina的股票数据统一从数据接口中获取，直接爬取的网页内容无法获取数据——可能原因：数据是异步获取的，爬虫环境不比游览器环境强大无法执行js

func main() {
	//var url string = "http://211.159.178.124/#/code.html"
	//fileName := "sina.xml"
	var dataInterface string
	var destJSONName string
	//临时文件
	datafile := "dataFromSina.txt"
	//GetPage(url, fileName)

	useFlag(code, &dataInterface)
	useFlag(fileName, &destJSONName)
	flag.Parse()
	dataInterface = dataInterBase + dataInterface
	data := GetPage(dataInterface, datafile)
	//fmt.Println(string(data))
	dataStruct := ParseData(string(data))
	data, err := json.Marshal(dataStruct)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(destJSONName)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(file, string(data))
	//defer log.Fatalln("json文件生成成功！")

}

//第一个标记我这次要输入的是什么
func useFlag(mark Mark, a *string) {
	switch mark {
	case code:
		flag.StringVar(a, "c", "sh000001","输入你要爬取的股票代码")
	case fileName:
		flag.StringVar(a, "d", "default.json", "输入你希望保存到的文件名")
	}

}

//根据url得到内容，并将网页主体存入fileName制定的文件中
//返回文件中存入的内容，也就是网页主体
//当fileName为空时不存入文件
func GetPage(url, fileName string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	bodyGB, err := ioutil.ReadAll(resp.Body)
	body := make([]byte, len(bodyGB))
	iconv.Convert(bodyGB, body, "gb2312", "utf-8")
	if err != nil {
		panic(err)
	}

	file.Write(body)
	//必须延时执行不然body无法正常返回
	//defer log.Fatalln("获取数据成功")
	return body
}

//分析从数据接口中获取的数据
func ParseData(str string) (result *DataStruct) {
	tempStr := strings.SplitN(str, "_", 3)[2]
	tempstrSlice := strings.Split(tempStr, "=")
	//分析出股票代码
	SCode := tempstrSlice[0]
	//修建掉两边的引号
	tempStr = strings.Trim(tempstrSlice[1], "\"")
	//分割成12项，前11项有作用 按照数据接口顺序 提取成交金额前的所有数据
	tempstrSlice = strings.SplitN(tempStr, ",", 11)

	var values []float64 = make([]float64, 0, 9)
	for i, value := range tempstrSlice {
		if i == 0 || i == (len(tempstrSlice)-1) {
			continue
		}
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}
		values = append(values, v)
	}

	result = NewDataStruct(SCode, tempstrSlice[0], values)
	//defer log.Fatalln("解析数据成功")
	return

}
