package main

import (
	"bufio"
	"bytes"
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
const geminiPort = 1965
const defaultConncetionTimeout = 15 * time.Second

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
	srvAddr := fmt.Sprintf("%v:%d", destination, geminiPort)
	fmt.Printf("Attempting to visit --> %v... \n", srvAddr)

	conn, err := openConn(srvAddr)
	if err != nil {
		return errors.Wrap(err, "dialing tcp address")
	}
	defer closeConn(conn)
	url := fmt.Sprintf("gemini://%v/\r\n", destination)
	written, err := fmt.Fprint(conn, url)
	if err != nil {
		return errors.Wrap(err, "writing url to connection")
	}
	fmt.Printf("bytes written to connection: %d\n", written)
	lines, err := readResponse(conn)
	if err != nil {
		return err
	}
	fmt.Printf("%s", lines)
	return nil
}

//go:generate moq -out mocks_test.go . reader
type reader interface {
	Read([]byte) (int, error)
}

func readResponse(r reader) (lines []string, err error) {
	buff := make([]byte, 1)
	lineBuff := make([]byte, 0)
	newLineDelimiter := []byte("\r\n")
	for {
		readCount, err := r.Read(buff)
		if err == io.EOF && readCount <= 0 {
			break
		} else if err != nil && err != io.EOF {
			return nil, errors.Wrap(err, "read conn")
		}
		lineBuff = append(lineBuff, buff...)
		if bytes.HasSuffix(lineBuff, newLineDelimiter) {
			// write out the line and flush the line buffer
			lines = append(lines, string(lineBuff[:len(lineBuff)-len(newLineDelimiter)]))
			lineBuff = make([]byte, 0)
		}
	}
	if len(lineBuff) > 0 {
		// write out the last line if any
		lines = append(lines, string(lineBuff))
	}
	return lines, nil
}

func openConn(srvAddr string) (*tls.Conn, error) {
	conf := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// for the POC, just skip the TLS verification
		InsecureSkipVerify: true,
	}

	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: defaultConncetionTimeout},
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
