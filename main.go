package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const prompt = ">"

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
	geminiPort := 1965
	dest := strings.Join(tokens[1:], " ")
	srvAddr := fmt.Sprintf("%v:%d", dest, geminiPort)
	fmt.Printf("Attempting to visit --> %v... \n", srvAddr)
	conf := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 15 * time.Second},
		"tcp",
		srvAddr,
		conf)
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Print("could not close tcp connection")
		}
	}()
	if err != nil {
		return errors.Wrap(err, "dialing tcp address")
	}
	written, err := fmt.Fprintf(conn, "gemini://%v/\r\n", dest)
	if err != nil {
		return errors.Wrap(err, "writing to connection")
	}
	fmt.Printf("bytes written to connection: %d\n", written)
	buffer := make([]byte, 1)
	var line []byte
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF && n <= 0 {
			break
		} else if err != nil && err != io.EOF {
			return errors.Wrap(err, "read conn")
		}
		line = append(line, buffer...)
	}

	fmt.Println(string(line))
	return nil
}
