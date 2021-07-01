package main

import (
	"strconv"
	"strings"
)

type url struct {
	scheme  string
	address string
	port    int
}

func newUrl(urlStr string) url {
	const protocolDelimiter = "://"
	const addressDelimiter = ":"
	const portDelimiter = "/"
	url := url{}
	url.scheme = getUrlPart(&urlStr, protocolDelimiter)
	url.address = getUrlPart(&urlStr, addressDelimiter)
	port_string := getUrlPart(&urlStr, portDelimiter)
	port, _ := strconv.Atoi(port_string)
	url.port = port

	return url
}

func getUrlPart(urlStr *string, delimiter string) string {
	var part string
	protocolIdx := strings.Index(*urlStr, delimiter)
	if protocolIdx > 0 {
		part = (*urlStr)[:protocolIdx]
		*urlStr = (*urlStr)[protocolIdx+len(delimiter):]
	}
	return part
}
