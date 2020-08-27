package xshell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func handleList(args []string) string {
	var buffer bytes.Buffer
	var resources []*resource
	if args == nil {
		resources = cache.query()
	} else {
		resources = cache.match(args[0])
	}

	if len(resources) == 0 {
		buffer.WriteString("No resources found")
	} else {
		buffer.WriteString("(Resources List): \r\n")
		maxUsernameLen, maxPasswordLen := 10, 10
		for _, resource := range resources {
			if len(resource.Username) > maxUsernameLen {
				maxUsernameLen = len(resource.Username)
			}
			if len(resource.Password) > maxPasswordLen {
				maxPasswordLen = len(resource.Password)
			}
		}
		columnLenLimit["username"] = maxUsernameLen
		columnLenLimit["password"] = maxPasswordLen

		fmt.Println(columnLenLimit["password"])
		for _, title := range []string{"no", "host", "port", "username", "password", "count", "last-login-time"} {
			paddingStringValue(&buffer, fmt.Sprintf("+%s", "-"), "-", columnLenLimit[title])
		}
		buffer.WriteString("+\r\n")

		for _, title := range []string{"no", "host", "port", "username", "password", "count", "last-login-time"} {
			paddingStringValue(&buffer, fmt.Sprintf("| %s", title), " ", columnLenLimit[title])
		}
		buffer.WriteString("|\r\n")

		for _, title := range []string{"no", "host", "port", "username", "password", "count", "last-login-time"} {
			paddingStringValue(&buffer, fmt.Sprintf("+%s", "-"), "-", columnLenLimit[title])
		}
		buffer.WriteString("+\r\n")

		for i, resource := range resources {
			buffer.WriteString(formatResourceToTable(resource, maxUsernameLen, maxPasswordLen))
			if i != len(resources)-1 {
				buffer.WriteString("\r\n")
			}
		}
		buffer.WriteString("\r\n")

		for _, title := range []string{"no", "host", "port", "username", "password", "count", "last-login-time"} {
			paddingStringValue(&buffer, fmt.Sprintf("+%s", "-"), "-", columnLenLimit[title])
		}
		buffer.WriteString("+\r\n")
	}
	return buffer.String()
}

func handleAdd(args []string) (string, error) {
	if args == nil {
		return "", fmt.Errorf("Missing resource configuration information")
	}

	var resource = &resource{}
	err := json.Unmarshal([]byte(args[0]), resource)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal resource configuration information, nest error: %v", err)
	}

	if err := resource.verify(); err != nil {
		return "", fmt.Errorf("Failed to verify resource configuration information, nest error: %v", err)
	}

	affected := cache.insert(resource)
	if affected != 1 {
		return "FAILURE", nil
	}

	if err := cache.dump(); err != nil {
		return fmt.Sprintf("Warning: Add resource on cache success, but dump to file failure, nest error: %v", err), nil
	}

	return "SUCCESS", nil
}

func handleDel(args []string) (string, error) {
	if args == nil {
		return "", fmt.Errorf("Missing resource-no")
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("Please enter resource-no")
	}

	affected := cache.delete(index)
	if affected == 0 {
		return "FAILURE", nil
	}

	if err := cache.dump(); err != nil {
		return fmt.Sprintf("Warning: Del resource on cache success, but dump to file failure, nest error: %v", err), nil
	}

	return "SUCCESS", nil
}

func handleMod(args []string) (string, error) {
	if args == nil || len(args) <= 1 {
		return "", fmt.Errorf("Missing resource-no or resource configuration information")
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("Please enter resource-no")
	}

	var resource = &resource{}
	err = json.Unmarshal([]byte(args[1]), resource)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal resource configuration information, nest error: %v", err)
	}

	if err := resource.verify(); err != nil {
		return "", fmt.Errorf("Failed to verify resource configuration information, nest error: %v", err)
	}

	affected := cache.update(index, resource)
	if affected == 0 {
		return "FAILURE", nil
	}

	if err := cache.dump(); err != nil {
		return fmt.Sprintf("Warning: Update resource on cache success, but dump to file failure, nest error: %v", err), nil
	}

	return "SUCCESS", nil
}

func handleLogin(args []string) error {
	if args == nil {
		return fmt.Errorf("Missing resource-no")
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Please enter resource-no")
	}

	resource := cache.get(index)
	if resource == nil {
		return fmt.Errorf("No resource found")
	}

	session, err := newSessionSSH(resource.Host, resource.Port, resource.Username, resource.Password, 20*time.Second)
	if err != nil {
		return fmt.Errorf("Failed to login resource [%v:%v/%v], nest error: %v", resource.Host, resource.Port, resource.Username, err)
	}
	defer session.Close()

	console := buildConsole(session)

	fmt.Fprintln(os.Stdout, fmt.Sprintf("Logging on resource [%v:%d/%s] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", resource.Host, resource.Port, resource.Username))
	if err := console.interactiveSession(); err != nil {
		return fmt.Errorf("Warning: interactive with ssh connection failure, nest error: %v", err)
	}
	console.close()
	console.complete()
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Logout from resource [%v:%d/%s] <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", resource.Host, resource.Port, resource.Username))
	resource.Count = resource.Count + 1
	resource.LastLoginTime = time.Now()

	affected := cache.update(index, resource)
	if affected == 0 {
		log.Printf("Warning: Update resource on cache failure")
	} else {
		if err := cache.dump(); err != nil {
			log.Printf("Warning: Update resource on cache success, but dump to file failure, nest error: %v\r\n", err)
		}
	}
	return nil
}

func handleHelp(args []string) string {
	var buffer bytes.Buffer
	buffer.WriteString(`
Usage:
	command [args]

Available Commands:
	list[l]:                            展示已存储资源配置信息
	list[l] | grep <args>:              过滤已存储资源配置信息
	add <resource-json>:                存储资源配置信息
	del <resource-no>:                  删除资源配置信息
	mod <resource-no> <resource-json>:  更新资源配置信息
	login <resource-no>:                登录资源
	quit:                               退出 xshell-go

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
	
	`)
	return buffer.String()
}

func formatResourceToTable(data *resource, usernameLen, passwordLen int) string {
	var buffer bytes.Buffer
	paddingStringValue(&buffer, fmt.Sprintf("| %d", data.No), " ", maxNoLen)
	paddingStringValue(&buffer, fmt.Sprintf("| %s", data.Host), " ", maxHostLen)
	paddingStringValue(&buffer, fmt.Sprintf("| %d", data.Port), " ", maxPortLen)
	paddingStringValue(&buffer, fmt.Sprintf("| %s", data.Username), " ", usernameLen)
	paddingStringValue(&buffer, fmt.Sprintf("| %s", data.Password), " ", passwordLen)
	paddingStringValue(&buffer, fmt.Sprintf("| %d", data.Count), " ", maxCountLen)
	paddingStringValue(&buffer, fmt.Sprintf("| %s", data.LastLoginTime.Format("2006-01-02 15:04:05")), " ", maxTimeLen)
	buffer.WriteString("|")
	return buffer.String()
}

func paddingStringValue(buffer *bytes.Buffer, value string, delim string, maxLen int) {
	buffer.WriteString(value)
	for i := 0; i <= maxLen-len(value)+1; i++ {
		buffer.WriteString(delim)
	}
}

var (
	columnLenLimit = map[string]int{
		"no":              maxNoLen,
		"host":            maxHostLen,
		"port":            maxPortLen,
		"username":        10,
		"password":        10,
		"count":           maxCountLen,
		"last-login-time": maxTimeLen,
	}
)

const (

	//10
	maxNoLen = 2 + 1

	// 123.123.123.123
	maxHostLen = 15 + 1

	// 65535
	maxPortLen = 5 + 1

	// 1234
	maxCountLen = 6 + 1

	// 2006-01-02 15:04:05
	maxTimeLen = 20 + 1
)
