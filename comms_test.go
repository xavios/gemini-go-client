package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_communication(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:1965")
	require.NoError(t, err, "tcp lister creation")
	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	ts.Listener.Close()
	ts.Listener = l
	ts.StartTLS()
	defer ts.Close()

	var r *request
	u := newUrl(ts.URL)
	r = newRequest(
		withUrl(u),
		withTimeout(defaultConncetionTimeout),
	)
	defer func() {
		require.NoError(t, r.Close(), "request close")
	}()
	require.NotNil(t, r)
	require.Equal(t, r.url, u)
	require.Equal(t, r.timeout, defaultConncetionTimeout)

	var conn *connection
	conn, err = r.Make()
	require.NoError(t, err, "connecting")
	require.NotNil(t, conn, "connection is nil")

	var resp *response
	resp, err = conn.ReadResponse()
	require.NotNil(t, resp, "response is nil")
	require.NoError(t, err, "reading response")
	// TODO: bad request - this is an HTTP server not a gemini one
	// require.ElementsMatch(t, []string{"Hello, client"}, resp.body)
}

func Test_readResponse(t *testing.T) {
	testCases := map[string]struct {
		resp      string
		wantLines []string
	}{
		"two lines received": {
			resp:      "This is a line\r\nthis is another",
			wantLines: []string{"This is a line", "this is another"},
		},
		"empty row represented": {
			resp:      "test\r\n\r\ntest1",
			wantLines: []string{"test", "", "test1"},
		},
		"line with breakline is handled as a single line": {
			resp:      "test\r\n",
			wantLines: []string{"test"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := readerMock{}
			i := 0
			r.ReadFunc = func(bytes []byte) (int, error) {
				// read 1 character from the response at a time
				if i < len(tc.resp) {
					bytes[0] = byte(tc.resp[i])
					i++
					return 1, nil
				}
				return 0, io.EOF
			}
			conn := connection{&r}
			resp, err := conn.ReadResponse()
			require.NoError(t, err)
			require.ElementsMatch(t, tc.wantLines, resp.body)
		})
	}
}
