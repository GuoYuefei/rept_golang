## 使用golang写的爬虫

> 暂且只写了对sina股票数据的爬取

1. sina股票单次爬取

   [sina.go]:(https://github.com/GuoYuefei/rept_golang/blob/master/src/cmd/sina.go)

   对应的可执行文件在release文件夹中。使用方式：命令行程序。

   切换到可执行文件所在文件夹后

   使用实例如下：

   ```shell
   # 注意目录
   Administrator@PC-201809211459 MINGW64 ~/Desktop/rept
   $ ./sina.exe -c "sz000950" -d "重药控股.json"
   Administrator@PC-201809211459 MINGW64 ~/Desktop/rept
   $ ./sina.exe -c "sh600176" -d "巨石集团.json"
   ```

   该命令行程序暂且只有两个可选选项

   -c：可选选项，默认是“sh000001”代表沪指。需要一个string参数，是需要用户提供一个需要爬取股票信息的代码。前面两个字母表示上市交易所，上交所则是sh，深交所则是sz。后面的数字是对应的股票代码。**比如**要爬取“巨石集团”的相关数据，我们需要了解到它在上交所上市，所以前两字母是sh，它的股票代码是600176，结合起来就是sh600176.

   -d：可选选项，默认是“default.json”。需要一个string参数，爬取的数据会以json形式保存在指定的的文件中。若文件不存在便会创建文件，若存在则覆盖。


