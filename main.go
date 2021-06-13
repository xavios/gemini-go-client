package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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
				fmt.Println("Error: ", err.Error())
			}
			continue
		}
		fmt.Printf("%v You entered: %q\n", prompt, text)
	}
}

func visit(tokens []string) error {
	dest := strings.Join(tokens[1:], " ")
	srvAddr := fmt.Sprintf("%v:1965", dest)
	fmt.Printf("Attempting to visit --> %v... \n", srvAddr)
	conn, err := net.Dial("tcp", srvAddr)
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Print("could not close tcp connection")
		}
	}()
	if err != nil {
		return errors.Wrap(err, "dialing tcp address")
	}
	written, err := fmt.Fprintf(conn, "%v\r\n", dest)
	if err != nil {
		return errors.Wrap(err, "writing to connection")
	}
	fmt.Printf("bytes written to connection: %d\n", written)
	var resp []byte
	read, err := conn.Read(resp)
	if err != nil {
		return errors.Wrap(err, "read from connection")
	}
	fmt.Printf("bytes read fromconnection: %d\n", read)

	fmt.Println(string(resp))
	return nil
}
