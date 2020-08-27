package xshell

import (
	"bytes"
	"fmt"
	"strings"
)

/*
*	list														-- list all
*	list | grep 123 f											-- grep
*	add {"host":"192.168.180.67"}								-- add
*	del 2														-- del
*	mod 3 {"host":"192.168.180.67"}								-- mod
*	login 3														-- login
*   quit														-- quit
 */

const (
	list  = "list"
	l     = "l"
	add   = "add"
	del   = "del"
	mod   = "mod"
	login = "login"
	quit  = "quit"
	help  = "help"
	h     = "h"
	exit  = "exit"
)

func parseStdinRequest(text []byte) (string, []string, error) {
	var buf = bytes.TrimSpace(text)

	var request string
	var args []byte
	for i, b := range buf {
		if b == 32 {
			request = string(buf[:i])
			args = bytes.TrimSpace(buf[i:])
			break
		}
		if i == len(buf)-1 {
			request = string(buf)
		}
	}

	switch request {
	case list, l:
		if len(args) == 0 {
			return request, nil, nil
		}
		if !bytes.HasPrefix(args, []byte("|")) {
			return "", nil, fmt.Errorf("Unsupported parameter format")
		}

		var n = bytes.Index(args, []byte("grep "))

		if n == -1 {
			return "", nil, fmt.Errorf("Unsupported parameter format")
		}

		if strings.TrimSpace(string(args[:n])) != "|" {
			return "", nil, fmt.Errorf("Unsupported parameter format")
		}

		return request, []string{strings.TrimSpace(string(args[n+5:]))}, nil
	case add, del, login:
		var data []string
		if args != nil {
			data = []string{strings.TrimSpace(string(args))}
		}
		return request, data, nil

	case mod:
		var data []string
		for i, b := range args {
			if b == 32 && i != len(args)-1 {
				data = append(data, string(args[:i]))
				data = append(data, string(args[i+1:]))
				break
			}
			if i == len(args)-1 {
				data = append(data, string(args))
			}
		}
		return mod, data, nil

	case help, h:
		return help, nil, nil

	case quit, exit:
		return request, nil, nil

	default:
		return "", nil, fmt.Errorf("Unsupported operation instructions-[%v]", request)

	}
}
