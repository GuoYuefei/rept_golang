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
const (
	fileJsonType = ".json"
)
//默认在当前工作目录下
var JsonPathRoot string = "./testSinaGo/"
var datafile string = JsonPathRoot + "dataFromSina.txt"

//股票结构体，只需要记录名字和代码
//这种做法导致库无法使用，可以使用mao代替
type Stock struct {
	Name string
	Code string
}
func NewStock(name, code string) *Stock {
	return &Stock{
		name,
		code,
	}
}

//这个map中只存两样东西，key分别：一个是Name一个是Code
//type Stock map[string]string
//func (s Stock)Name() string {
//	return s["Name"]
//}
//func (s Stock)Code() string {
//	return s["Code"]
//}


func (s *Stock)String() string {
	return "{\n\tName:"+s.Name+",\n\tCode:"+s.Code+"\n}"
}



//一个用于记录需要查询的股票集合
//以后可能还有增添选项
//type Stocks struct {
//	StockArr []Stock
//}


//sina的股票数据统一从数据接口中获取，直接爬取的网页内容无法获取数据——可能原因：数据是异步获取的，爬虫环境不比游览器环境强大无法执行js

func main() {

	//var url string = "http://211.159.178.124/#/code.html"
	//fileName := "sina.xml"
	var dataInterface string
	var destJSONName string
	//临时文件
	//GetPage(url, fileName)

	useFlag(code, &dataInterface)
	useFlag(fileName, &destJSONName)
	flag.Parse()
	dataInterface = dataInterBase + dataInterface
	data := GetPage(dataInterface, datafile)
	//fmt.Println(string(data))
	dataStruct := ParseData(string(data))

	StoreJson(dataStruct, destJSONName)
	//defer log.Fatalln("json文件生成成功！")

	jsonBytes, _ := ioutil.ReadFile(JsonPathRoot+"123.json")
	GetDatas(jsonBytes)

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

func StoreJson(a interface{}, storeFile string) {
	flag := false		//如果执行过创建目录一次就置true
	data, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	storeFile = JsonPathRoot + storeFile
FLAGCREATE:
	file, err := os.Create(storeFile)
	if err != nil && !flag {
		e := os.MkdirAll(JsonPathRoot,666)			//三组用户都有读写权
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

//得到止一个数据时可以调用该函数
//该函数需要一个reader接口类型的参数,这个参数可以是各种输入，建议使用提前使用好的文件
//文件中提供需要查询的各种股票的代码，根据这些进行进行爬取数据
func GetDatas(whos []byte) {

	fmt.Println("read informations: " + string(whos))

	str := string(whos)
	strs := strings.Split(str, ";")
	var stocks []*Stock = make([]*Stock,len(strs))			//strs的长度就是将来map数组的长度
	//将文件读出来的byte切片分别转换成对象  嗯。。。因为这个json文件是非标准的，在解析前需要用过“；”分离找到对象
	//也是因为这个官方库不好用的原因。。。它无法解析数组中带有对象 也就是它不能解析引用类型中有引用类型的情况
	for i, v := range strs {
		fmt.Println(i,v)
		stocks[i] = NewStock("","")
		err := json.Unmarshal([]byte(v),stocks[i])
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("stocks's len:", len(stocks))
	for _, v := range stocks {
		fmt.Println(v)
		data := GetPage(dataInterBase+v.Code, datafile)
		dataStruct := ParseData(string(data))
		StoreJson(dataStruct, v.Name +fileJsonType)
	}
}
