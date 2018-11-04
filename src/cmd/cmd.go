package main

import (
	"fmt"
	"strings"
)

func main() {
	sss := "12asd,ewq,rwe,tr,ytur,ur,sfe12"
	re := strings.Split(sss, ",")
	fmt.Println(re)
	re = strings.SplitAfter(sss, ",")
	fmt.Println(re)
	//前两个是分割后的字符串，后面一个是剩下没分割的
	//有After的字串里带“，”
	re = strings.SplitAfterN(sss, ",", 3)
	fmt.Println(re)
	re = strings.SplitN(sss, ",", 3)
	fmt.Println(re)
	reStr := strings.Trim(sss, "12")
	fmt.Println(reStr)

}
