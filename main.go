// docker-chown | (c) 2020 Icinga GmbH | GPLv2+

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type dirInfo interface {
	Name() string
	IsDir() bool
}

var _ dirInfo = os.FileInfo(nil)

type currentDir struct{}

var _ dirInfo = currentDir{}

func (currentDir) Name() string {
	return "."
}

func (currentDir) IsDir() bool {
	return true
}

var me = syscall.Getuid()
var us = syscall.Getgid()

func main() {
	// Rationale:
	// https://github.com/golang/go/blob/8cd75f3/src/syscall/syscall_linux.go#L960-L971

	{
		it := syscall.Geteuid()
		if _, _, errSU := syscall.RawSyscall(syscall.SYS_SETUID, uintptr(it), 0, 0); errSU != 0 {
			fmt.Fprintf(os.Stderr, "syscall.RawSyscall(syscall.SYS_SETUID, %d, 0, 0): %s\n", it, errSU.Error())
			os.Exit(1)
		}
	}

	{
		them := syscall.Getegid()
		if _, _, errSG := syscall.RawSyscall(syscall.SYS_SETGID, uintptr(them), 0, 0); errSG != 0 {
			fmt.Fprintf(os.Stderr, "syscall.RawSyscall(syscall.SYS_SETGID, %d, 0, 0): %s\n", them, errSG.Error())
			os.Exit(1)
		}
	}

	chownTree(currentDir{}, "/data")
}

func chownTree(root dirInfo, relativeTo string) {
	fullPath := path.Join(relativeTo, root.Name())
	if errCO := os.Chown(fullPath, me, us); errCO != nil {
		fmt.Fprintf(os.Stderr, "os.Chown(%s, %d, %d): %s\n", fullPath, me, us, errCO.Error())
	}

	if root.IsDir() {
		if branches, errRD := ioutil.ReadDir(fullPath); errRD == nil {
			for _, branch := range branches {
				chownTree(branch, fullPath)
			}
		} else {
			fmt.Fprintf(os.Stderr, "ioutil.ReadDir(%s): %s\n", fullPath, errRD.Error())
		}
	}
}
