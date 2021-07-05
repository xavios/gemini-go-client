package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

type request struct {
	url     url
	timeout time.Duration
	conn    *connection
}

func (r *request) Make() (*connection, error) {
	conf := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// for the POC, just skip the TLS verification
		InsecureSkipVerify: true,
	}

	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: r.timeout},
		"tcp",
		r.url.ServerAddress(),
		conf)
	r.conn = &connection{conn}

	if err != nil {
		return nil, errors.Wrap(err, "tcp connection")
	}

	request := fmt.Sprintf("%v%v", r.url.String(), clrf)
	_, err = fmt.Fprint(conn, request)
	if err != nil {
		defer conn.Close()
		return nil, errors.Wrap(err, "writing url to connection")
	}

	return &connection{conn}, err
}

func newRequest(opts ...requestOptions) *request {
	r := &request{}
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

type requestOptions func(*request)

func withUrl(url url) requestOptions {
	return func(r *request) {
		r.url = url
	}
}

func withTimeout(timeout time.Duration) requestOptions {
	return func(r *request) {
		r.timeout = timeout
	}
}

func (r *request) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

type response struct {
	// TODO
	status int
	head   string
	body   []string
}

//go:generate moq -out mocks_test.go . reader
type reader interface {
	Read([]byte) (int, error)
	Close() error
}

type connection struct {
	reader
}

func (c *connection) ReadResponse() (*response, error) {
	resp := &response{
		body: make([]string, 0),
	}
	buff := make([]byte, 1)
	lineBuff := make([]byte, 0)
	newLineDelimiter := []byte(clrf)
	for {
		readCount, err := c.Read(buff)
		if err == io.EOF && readCount <= 0 {
			break
		} else if err != nil && err != io.EOF {
			return nil, errors.Wrap(err, "read conn")
		}
		lineBuff = append(lineBuff, buff...)
		if bytes.HasSuffix(lineBuff, newLineDelimiter) {
			// write out the line and flush the line buffer
			resp.body = append(resp.body, string(lineBuff[:len(lineBuff)-len(newLineDelimiter)]))
			lineBuff = make([]byte, 0)
		}
	}
	if len(lineBuff) > 0 {
		// write out the last line if any
		resp.body = append(resp.body, string(lineBuff))
	}
	return resp, nil
}
