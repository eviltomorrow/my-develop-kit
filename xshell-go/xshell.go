package xshell

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	nmPath    = "path"
	nmVersion = "version"
)

const (
	prompt = "[xshell-go]$ "
)

var (
	path    = flag.String("path", "resource.db", "resource path for xshell-go")
	version = flag.Bool("version", false, "xshell-go version")
)

// Run 入口启动
func Run() {
	flag.Parse()

	printWelcomeInformation()
	loadResourceCacheFromFile(*path)

	var reader = bufio.NewReader(os.Stdin)
loop:
	for {
		fmt.Fprintf(os.Stdout, prompt)

		data, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintln(os.Stdout, "Panic: reade data from os.Stdin failure, nest error: ", err)
			continue
		}

		if len(data) <= 1 {
			continue
		}

		request, args, err := parseStdinRequest(data[:len(data)-1])
		if err != nil {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%v\r\n", err))
			continue
		}

		switch request {
		case list, l:
			fmt.Fprintln(os.Stdout, handleList(args))

		case add:
			result, err := handleAdd(args)
			if err != nil {
				fmt.Fprintln(os.Stdout, err)
			} else {
				fmt.Fprintln(os.Stdout, result)
			}

		case del:
			result, err := handleDel(args)
			if err != nil {
				fmt.Fprintln(os.Stdout, err)
			} else {
				fmt.Fprintln(os.Stdout, result)
			}

		case mod:
			result, err := handleMod(args)
			if err != nil {
				fmt.Fprintln(os.Stdout, err)
			} else {
				fmt.Fprintln(os.Stdout, result)
			}

		case help, h:
			fmt.Fprintln(os.Stdout, handleHelp(args))

		case login:
			err := handleLogin(args)
			if err != nil {
				fmt.Fprintln(os.Stdout, err)
			}
		case quit, exit:
			break loop

		default:

		}
		fmt.Println()
	}

	fmt.Fprintf(os.Stdout, "Bye\r\n")
}

func printWelcomeInformation() {
	var version = `Welcome to the xshell-go.  Commands end with \n
	
Application version: v2.0.1 package
	
Copyright (c) 2020 github/eviltomorrow and others.
	
Type help \n for help. Have a good luck.
`
	fmt.Fprintln(os.Stdout, version)
}

func loadResourceCacheFromFile(path string) {
	cache = &db{
		Path:      path,
		Resources: make([]*resource, 0, 20),
	}

	info, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("The system can not run xhell-go, nest error: %v", err)
			os.Exit(0)
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Printf("Failed to create xhell-go resource file, nest error: %v", err)
			os.Exit(0)
		}
		defer file.Close()
		return
	}

	if info.IsDir() {
		log.Printf("Failed to create xhell-go resource file, nest error: Already exist same folder")
		os.Exit(0)
	}

	if err := cache.load(); err != nil {
		_, err = os.OpenFile(path, os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("Failed to truncate xhell-go resource file, nest error: %v", err)
			os.Exit(0)
		}
	}
}
