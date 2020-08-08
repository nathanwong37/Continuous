package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/temp"
)

func main() {
	fmt.Println("Welcome to Temp Name, please enter your configuration \n 1. Local \n 2. Wide \n 3. Personal \n (Default is Local)")
	var input string
	fmt.Scanln(&input)
	var conf *memberlist.Config
	switch input {
	case "Wide":
		fmt.Println("Wide area Network Chosen")
		conf = memberlist.DefaultWANConfig()
		conf.BindPort = 8301
		conf.AdvertisePort = 8301
	case "Personal":
		fmt.Println("Personal area network chosen")
		conf = memberlist.DefaultLocalConfig()
	default:
		fmt.Println("Local Area Network chosen")
		conf = memberlist.DefaultLANConfig()
	}
	fmt.Println("Please enter a text file with nodes to connect to")
	fmt.Scanln(&input)
	lines, err := readLines(input)
	if err != nil {
		panic(err)
	}
	messeng := temp.NewMessenger(conf)
	_, err = messeng.Join(lines)
	if err != nil {
		fmt.Println("Error in joining")
	}
	done := make(chan bool)
	go forever()
	<-done
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func forever() {
	for {
		fmt.Println("Force quit the program to exit")
		time.Sleep(100 * time.Second)
	}
}
