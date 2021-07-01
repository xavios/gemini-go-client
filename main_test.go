package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

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
			lines, err := readResponse(&r)
			require.NoError(t, err)
			require.ElementsMatch(t, tc.wantLines, lines)
		})
	}

}
