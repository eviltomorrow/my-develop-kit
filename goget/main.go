package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

var reg = regexp.MustCompile(`.*cannot find package ".*" in any of.*`)

var f = flag.String("f", "", "path of file")

func main() {
	flag.Parse()
	if *f == "" {
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.OpenFile(*f, os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("[fatal] open file failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// setenv()

	var packageList = make([]string, 0, 32)
	var scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var text = strings.TrimSpace(scanner.Text())
		if text != "" && reg.MatchString(text) {
			left, right := strings.Index(text, `"`), strings.LastIndex(text, `"`)
			if left == -1 || right == -1 {
				continue
			}
			packageList = append(packageList, string(text)[left+1:right])
		}
	}
	if len(packageList) == 0 {
		log.Printf("[warning] no package will be [go get]")
		return
	}

	for _, pname := range packageList {
		log.Printf("[info] go-cmd: go get -v %s(downloading...)\r\n", pname)
	}
	for _, pname := range packageList {
		execcmd(fmt.Sprintf("go get -v %s", pname))
	}
}

// func setenv() {
// 	var env = os.Getenv("http_proxy")
// 	if env == "" {
// 		os.Setenv("http_proxy", "socks5://127.0.0.1:1080")
// 	}

// 	env = os.Getenv("http_proxys")
// 	if env == "" {
// 		os.Setenv("https_proxy", "socks5://127.0.0.1:1080")
// 	}
// }

// ExecuteCmd 执行 command
func execcmd(c string) error {
	var cmd = exec.Command("/bin/sh", "-c", c)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start execute cmd failure, nest error: %v", err)
	}
	defer syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

	go func() {
		io.Copy(os.Stdout, stdout)
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	return cmd.Wait()
}
