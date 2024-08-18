# dnsping-go

`dnsping-go` 是一个用 Go 语言编写的程序，通过指定的 DNS 服务器查询某个域名，并输出查询结果和查询时间，然后 `tcpping` 该结果并显示延时。

## 参数

- -d    要测试的域名 (必须)
- -s    dns服务器 可多个用“,”连接 (必须)
- -t    4/6 (ipv4/ipv6，默认4)
- -v   显示版本信息

## 安装

1. 克隆仓库：
    ```sh
    git clone https://github.com/charley008/dnsping-go.git
    ```
2. 进入项目目录：
    ```sh
    cd dnsping-go
    ```
3. 初始化模块：
    ```sh
    go mod tidy
    ```

## 使用

运行以下命令来查询域名并显示 `tcpping` 延时：

```sh
go run main.go -d www.qq.com -s 1.1.1.1,8.8.8.8,223.5.5.5 -4

