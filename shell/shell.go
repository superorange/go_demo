package main

import (
	"fmt"
	"syscall"
)

func main() {
	// 获取当前进程ID
	pid, _, errno := syscall.RawSyscall(syscall.SYS_GETPID, 0, 0, 0)
	if errno != 0 {
		fmt.Printf("RawSyscall error: %v\n", errno)
		return
	}
	fmt.Printf("Current PID: %d\n", pid)
}
