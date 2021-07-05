package main

import (
	"fmt"
	"strconv"
	"strings"
)

type url struct {
	scheme  string
	address string
	port    int
	path    string
	query   [][]string
}

func (u *url) String() string {
	s := fmt.Sprintf("%s://%s/", u.scheme, u.address)
	if len(u.path) > 0 {
		s += u.path
	}
	var query string
	for _, kv := range u.query {
		query += "&" + kv[0] + "=" + kv[1]
	}
	if len(query) > 0 {
		query = "?" + query[1:]
		s += query
	}
	return s
}

func (u *url) ServerAddress() string {
	return fmt.Sprintf("%s:%d", u.address, u.port)
}

func newUrl(urlStr string) url {
	const schemeDelim = "://"
	const addressDelim = ":"
	const portDelim = "/"
	const pathDelim = "?"
	const queryDelim = "&"
	const queryKVDelim = "="

	const defaultScheme = "gemini"
	const defaultPort = 1965

	var port_string string

	url := url{}
	url.scheme = getUrlPart(&urlStr, schemeDelim)
	if url.scheme == "" {
		url.scheme = defaultScheme
	}
	url.address = getUrlPart(&urlStr, addressDelim)
	if url.address == "" {
		// no port in urlStr
		url.address = getUrlPart(&urlStr, portDelim)
		url.port = defaultPort
		if url.address == "" {
			// no path in urlStr
			url.address = getUrlPart(&urlStr, pathDelim)
		}
		if url.address == "" {
			// only address in urlStr
			url.address = urlStr
		}
	} else {
		fmt.Println(urlStr)
		port_string = getUrlPart(&urlStr, portDelim)
		if port_string == "" {
			port_string = urlStr
		}
		port, _ := strconv.Atoi(port_string)
		url.port = port
	}
	url.path = getUrlPart(&urlStr, pathDelim)
	if url.path == "" &&
		!strings.Contains(urlStr, queryDelim) &&
		urlStr != url.address &&
		urlStr != port_string {
		// at the end only the path remained
		url.path = urlStr
	}

	query := [][]string{}
	for _, queryPart := range strings.Split(urlStr, queryDelim) {
		keyAndValue := strings.Split(queryPart, queryKVDelim)
		if len(keyAndValue) > 1 {
			query = append(query, []string{keyAndValue[0], keyAndValue[1]})
		}
	}
	url.query = query

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
