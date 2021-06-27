package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func getProcessRunningStatus(pid int) (*os.Process, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return proc, nil
	}

	if err == syscall.ESRCH {
		return nil, errors.New("process not running")
	}

	// default
	return nil, errors.New("process running but query operation not permitted")
}

func getMessageProcess(pid int)(bool, error){
	check := true

	if runtime.GOOS != "windows" {
		_, err := getProcessRunningStatus(pid)

		if err != nil {
			return check, err
		}
	}

	p, err := ps.FindProcess(pid)

	if err != nil {
		fmt.Println("loop")
		return check, err
	}else{
		if p == nil {
			fmt.Println(p)
			check = false
		}
	}

	return check, err
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("-> ", netData)

		netData = strings.Replace(netData, "\n", "", -1)
		netData = strings.Replace(netData, "\r", "", -1)

		pid, err := strconv.Atoi(netData)

		if err != nil{
			c.Write([]byte(err.Error() + "\n"))
		}

		check, error1 := getMessageProcess(pid)

		if error1 != nil{
			c.Write([]byte(error1.Error() + "\n"))
		}else{
			if check == true{
				c.Write([]byte("Process maybe still alive" + "\n"))
			}else{
				c.Write([]byte(strconv.FormatBool(check) + "\n"))
			}
		}
	}
}

