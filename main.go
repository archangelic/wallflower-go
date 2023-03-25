package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	timestamp string
	user      string
	message   string
}

func parseLine(data string) Message {
	var msg Message

	// get timestamp
	msg.timestamp = strconv.FormatInt(time.Now().UTC().Unix(), 10)

	// get user
	userReg := regexp.MustCompile(`(?P<nick>.+)!(?P<user>.+)@(?P<host>.+)`)
	userstring := strings.Trim(strings.Split(data, " ")[0], ": ")
	matches := userReg.FindStringSubmatch(userstring)
	msg.user = matches[userReg.SubexpIndex("user")]

	// get message
	msg.message = strings.SplitAfterN(data, ":", 3)[2]

	return msg
}

func main() {
	server := "localhost"
	port := "6667"
	nick := "wallflower"
	channel := "tildetown"
	logfile := "/home/archangelic/wallflower/irc.log"

	conn, err := net.Dial("tcp", server+":"+port)

	if err != nil {
		fmt.Println(err)
	} else {
		defer conn.Close()
	}

	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err)
	} else {
		defer f.Close()
	}

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	for {
		data, err := tp.ReadLine()
		if err == nil {
			switch {
			case strings.Contains(data, "Found your username"):
				fmt.Fprintf(conn, "NICK "+nick+"\r\n")
				fmt.Fprintf(conn, "USER "+nick+" \""+nick+".com\" \""+server+"\" :"+nick+" robot\r\n")
			case strings.Contains(data, "End of MOTD command"):
				fmt.Fprintf(conn, "JOIN #"+channel+"\r\n")
			case strings.Contains(data, "PING "):
				fmt.Fprintf(conn, strings.Replace(data, "PING", "PONG", 1)+"\r\n")
			case strings.Contains(data, "PRIVMSG #"+channel):
				msg := parseLine(data)
				if _, err := f.WriteString(msg.timestamp + "\t" + msg.user + "\t" + msg.message + "\n"); err != nil {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println("Disconnected")
			break
		}
	}
}
