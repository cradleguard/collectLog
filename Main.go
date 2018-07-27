package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

const tmpDir = "Tmp"
const logName = "ServerLog.txt"

func createTmpDir() {
	dir := tmpDir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}

func clearTmpDir() {
	dir := tmpDir
	if err := os.RemoveAll(dir); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func copyLogsToTmpDir(serverPath string, serverName string) {
	logPath := (serverPath + "\\" + logName)
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		from, err := os.Open(logPath)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer from.Close()

		newLogName := tmpDir + "\\ServerLog_" + serverName + ".txt"
		to, err := os.OpenFile(newLogName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}

func copyLogsToOneFile(serverPath string, serverName string) {
	logPath := (serverPath + "\\" + logName)
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		from, err := os.Open(logPath)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer from.Close()

		newLogName := tmpDir + "\\ServerLog_All.txt"
		to, err := os.OpenFile(newLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer to.Close()

		bo := binary.LittleEndian // windows
		scanner := bufio.NewScanner(from)
		splitFunc, orderFunc := ScanUTF16LinesFunc(bo)
		scanner.Split(splitFunc)
		for scanner.Scan() {
			b := scanner.Bytes()
			s := UTF16BytesToString(b, orderFunc())
			text := serverName + " " + s + "\r\n"
			if _, err := to.WriteString(text); err != nil {
				fmt.Println(err)
				panic(err)
			}
		}
	} else {
		fmt.Println("Log not found: path " + logPath)
	}
}

func collectLogsToTmpDir(servers string, workerFunc func(serverPath string, serverName string)) {
	if servers == "" {
		return
	}

	serverList := strings.Split(servers, ",")
	for _, server := range serverList {
		serverPath := strings.TrimSpace(server)
		serverNameParts := strings.Split(serverPath, string(os.PathSeparator))
		serverName := ""
		for _, part := range serverNameParts {
			if part != "" {
				serverName = part
				break
			}
		}

		workerFunc(serverPath, serverName)
	}

}

func main() {
	args := os.Args[1:]
	fmt.Println(args[0])
	clearTmpDir()
	createTmpDir()
	//copyLogsToTmpDir("\\\\localhost\\Share")
	// collectLogsToTmpDir("\\\\localhost\\Share", copyLogsToOneFile)
	if len(args) > 0 {
		collectLogsToTmpDir(args[0], copyLogsToOneFile)
	}
}
