## PS
没有键盘监听事件，所以无法监听到键盘 <-  ->  等动作（左移，右移）

## Installation

To install xshell-go package, you need to install Go and set your Go workspace first.

1. The first need [Go](https://golang.org/) installed (**version 1.11+ is required**), then you can use the below Go command to install xshell-go.

```sh
$ go get -u github.com/eviltomorrow/xshell-go
```

2. Import it in your code:

```go
import "github.com/eviltomorrow/xshell-go"
```

## Quick start


```go
package main

import "github.com/eviltomorrow/xshell-go"

func main() {
	xshell.Run()
}
```

## Help Show

```text
Usage:   
    command [args]

Available Commands:
	list[l]:                            展示已存储资源配置信息
	list[l] | grep <args>:              过滤已存储资源配置信息
	add <resource-json>:                存储资源配置信息
	del <resource-no>:                  删除资源配置信息
	mod <resource-no> <resource-json>:  更新资源配置信息
	login <resource-no>:                登录资源
	quit :                              退出 xshell-go

Example:
  =>[xshell-go]$ list
	(Resources List): 
	+---+----------------+------+----------+----------+-------+---------------------+
	| no| host           | port | username | password | count | last-login-time     |
	+---+----------------+------+----------+----------+-------+---------------------+
	| 0 | 127.0.0.1      | 22   | root     | root     | 0     | 0001-01-01 00:00:00 |
	+---+----------------+------+----------+----------+-------+---------------------+

  =>[xshell-go]$ list | grep root
	(Resources List): 
	+---+----------------+------+----------+----------+-------+---------------------+
	| no| host           | port | username | password | count | last-login-time     |
	+---+----------------+------+----------+----------+-------+---------------------+
	| 0 | 127.0.0.1      | 22   | root     | root     | 0     | 0001-01-01 00:00:00 |
	+---+----------------+------+----------+----------+-------+---------------------+

  =>[xshell-go]$ add {"host":"127.0.0.1","port":22,"username":"root","password":"root"}
	SUCCESS

  =>[xshell-go]$ del 0
	SUCCESS

  =>[xshell-go]$ mod 0 {"host":"127.0.0.1","port":22,"username":"root","password":"root"}
	SUCCESS

  =>[xshell-go]$ login 0
	Logging on resource [127.0.0.1:22/root] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	Last login: Wed May 20 20:47:09 2020 from 10.10.10.10
	[root@localhost ~]# 
```
