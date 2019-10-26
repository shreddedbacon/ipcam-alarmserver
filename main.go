/*
A simple TCP server to listen to events from ONVIF/IP Camera that have an "AlarmServer" that can be configured

*/
package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"

	shinobiclient "github.com/shreddedbacon/shinobi-client"
)

type AlarmEvent struct {
	Address   string `json:"Address"`
	Channel   int    `json:"Channel"`
	Descrip   string `json:"Descrip"`
	Event     string `json:"Event"`
	SerialID  string `json:"SerialID"`
	StartTime string `json:"StartTime"`
	Status    string `json:"Status"`
	Type      string `json:"Type"`
}

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 15002, "The port to listen on; default is 8000.")
var config = flag.String("config", "config.json", "The config file to use; default is config.json.")

func main() {
	flag.Parse()
	shinobiServerConfig, err := ioutil.ReadFile(*config) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	sa := shinobiclient.New(string(shinobiServerConfig))

	fmt.Println("Starting server...")
	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}
		go handleConnection(conn, sa)
	}
}

func handleConnection(conn net.Conn, sa shinobiclient.ShinobiClient) {
	remoteAddr := conn.RemoteAddr().String()
	u, err := url.Parse("camera://" + remoteAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	scanner := bufio.NewScanner(conn)
	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		handleMessage(scanner.Text(), conn, sa)
	}
}

func handleMessage(message string, conn net.Conn, sa shinobiclient.ShinobiClient) {
	alarmMessage := strings.TrimSpace(message[20:])
	event := AlarmEvent{}
	json.Unmarshal([]byte(alarmMessage), &event)
	event.Descrip = strings.Replace(event.Descrip, "\n", "", -1) // replace new lines
	host := getIPAddress(event.Address)
	fmt.Println("Address:", host, event.Address, "Description:", event.Descrip, "Event:", event.Event, "Type:", event.Type, "Status:", event.Status)
	if event.Event == "MotionDetect" && event.Status == "Start" {
		str, err := sa.TriggerMotion(host)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(str)
		}
	}
}

func getIPAddress(hexstring string) string {
	var newAddress []string
	for _, sub := range splitSubN(reverse(hexstring)[:8], 2) {
		newAddress = append(newAddress, reverse(sub))
	}
	a, _ := hex.DecodeString(strings.Join(newAddress, ""))
	return fmt.Sprintf("%v.%v.%v.%v", a[0], a[1], a[2], a[3])
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func splitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}
