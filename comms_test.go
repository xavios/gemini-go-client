package main

import (
	"io"
	"testing"
	"time"

	"github.com/pitr/gig"
	"github.com/stretchr/testify/require"
)

func Test_communication(t *testing.T) {
	const testResponse = "Hello, client"

	g := gig.Default()
	g.Handle("/", func(c gig.Context) error {
		return c.Gemini(testResponse)
	})
	go g.Run("test/example.crt", "test/example.key")
	defer g.Close()

	// Wait a bit for gig server, so the request will be served
	time.Sleep(1 * time.Second)
	var r *request

	u := newUrl("localhost:1965")
	r = newRequest(
		withUrl(u),
		withTimeout(defaultConnectionTimeout),
	)
	defer func() {
		require.NoError(t, r.Close(), "request close")
	}()
	require.NotNil(t, r)
	require.Equal(t, r.url, u)
	require.Equal(t, r.timeout, defaultConnectionTimeout)

	var conn *connection
	conn, err := r.Make()
	require.NoError(t, err, "connecting")
	require.NotNil(t, conn, "connection is nil")

	var resp *response
	resp, err = conn.ReadResponse()
	require.NotNil(t, resp, "response is nil")
	require.NoError(t, err, "reading response")

	require.Equal(t, 20, resp.status)
	require.ElementsMatch(t, []string{testResponse}, resp.body)
}

func Test_readResponse(t *testing.T) {
	testCases := map[string]struct {
		resp       string
		wantLines  []string
		wantHeader string
		wantStatus int
	}{
		"two lines received": {
			resp:       "20 text/gemini\r\nThis is a line\r\nthis is another",
			wantHeader: "20 text/gemini",
			wantStatus: 20,
			wantLines:  []string{"This is a line", "this is another"},
		},
		"empty row represented": {
			resp:       "20 text/gemini\r\ntest\r\n\r\ntest1",
			wantHeader: "20 text/gemini",
			wantStatus: 20,
			wantLines:  []string{"test", "", "test1"},
		},
		"line with breakline is handled as a single line": {
			resp:       "20 text/gemini\r\ntest\r\n",
			wantHeader: "20 text/gemini",
			wantStatus: 20,
			wantLines:  []string{"test"},
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
			require.NoError(t, err, "reading response")
			require.NotNil(t, resp)
			require.ElementsMatch(t, tc.wantLines, resp.body)
			require.Equal(t, tc.wantHeader, resp.header)
			require.Equal(t, tc.wantStatus, resp.status)
		})
	}
}
