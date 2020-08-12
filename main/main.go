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
	messeng := startMain()
	fmt.Println("Running on :" + messeng.M.LocalNode().Address())
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

//startMain creates and joins the messenger, default port for WAN is 8301, while default port for LAN is 8300
func startMain() *temp.Messenger {
	fmt.Println("Welcome to Temp Name, please enter your configuration \n 1. Local \n 2. Wide \n 3. Personal \n (Default is Local)")
	var input string
	fmt.Scanln(&input)
	var conf *temp.MessengerConfig
	switch input {
	case "Wide":
		fmt.Println("Wide area Network Chosen")
		config := memberlist.DefaultWANConfig()
		config.BindPort = 8301
		config.AdvertisePort = 8301
		addr := getOutboundIP()
		config.BindAddr = addr.String()
		conf = temp.CustomConfig(config, false)
		fmt.Println("Port number is 8301")
	case "Personal":
		fmt.Println("Personal area network chosen")
		config := memberlist.DefaultLocalConfig()
		fmt.Println("Enter a Name")
		fmt.Scanln(&input)
		config.Name = input
		fmt.Println("Enter a Port number")
		var port int
		fmt.Scanln(&port)
		config.BindPort = port
		config.AdvertisePort = port
		conf = temp.CustomConfig(config, true)
	default:
		fmt.Println("Local Area Network chosen")
		config := memberlist.DefaultLANConfig()
		config.BindPort = 8300
		config.AdvertisePort = 8300
		addr := getOutboundIP()
		config.BindAddr = addr.String()
		conf = temp.CustomConfig(config, false)
		fmt.Println("Port number is 8300")
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
		fmt.Println(err.Error())
		//os.Exit(8)
	}
	return messeng

}
