package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const prompt = ">"
const geminiPort = 1965
const defaultConnectionTimeout = 15 * time.Second
const clrf = "\r\n"

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func tokenize(s string) []string {
	s = strings.TrimSpace(s)
	return strings.Split(s, " ")
}

func run() error {
	for {
		fmt.Printf("%v ", prompt)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		tokens := tokenize(text)
		switch tokens[0] {
		case "q":
			return nil
		case "visit":
			err := visit(tokens)
			if err != nil {
				fmt.Println("error: ", err.Error())
			}
			continue
		}
		fmt.Printf("%v You entered: %q\n", prompt, text)
	}
}

// command: visit gemini.circumlunar.space
func visit(tokens []string) error {
	destination := strings.Join(tokens[1:], " ")
	url := newUrl(destination)
	fmt.Printf("Attempting to visit --> %v... \n", url.ServerAddress())
	req := newRequest(
		withUrl(url),
		withTimeout(defaultConnectionTimeout),
	)
	defer req.Close()
	conn, err := req.Make()
	if err != nil {
		return err
	}
	resp, err := conn.ReadResponse()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", resp.body)
	return nil
}

func openConn(srvAddr string) (*tls.Conn, error) {
	conf := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// for the POC, just skip the TLS verification
		InsecureSkipVerify: true,
	}

	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: defaultConnectionTimeout},
		"tcp",
		srvAddr,
		conf)
	return conn, err
}

func closeConn(conn *tls.Conn) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			fmt.Print("could not close tcp connection")
		}
	}
}
