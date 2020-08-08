package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
		addr := getOutboundIP()
		conf.BindAddr = addr.String()
	case "Personal":
		fmt.Println("Personal area network chosen")
		conf = memberlist.DefaultLocalConfig()
	default:
		fmt.Println("Local Area Network chosen")
		conf = memberlist.DefaultLANConfig()
	}
	config := temp.CustomConfig(conf)
	fmt.Println("Please enter a text file with nodes to connect to")
	fmt.Scanln(&input)
	lines, err := readLines(input)
	if err != nil {
		panic(err)
	}
	messeng := temp.NewMessenger(config)
	_, err = messeng.Join(lines)
	if err != nil {
		fmt.Println(err.Error())
		//os.Exit(8)
	}
	done := make(chan bool)
	go forever()
	<-done
}

//read address from files
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

//makes the function run forever
func forever() {
	for {
		fmt.Println("Force quit the program to exit")
		time.Sleep(100 * time.Second)
	}
}

//gets the ipaddress of the machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
