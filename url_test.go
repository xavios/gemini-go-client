package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Url(t *testing.T) {
	testCases := map[string]struct {
		url            string
		wantUrl        url
		wantStringRepr string
	}{
		"have all parts": {
			url: "gemini://gemini.conman.org:1965/news.txt?hello=bello&test=test1",
			wantUrl: url{
				scheme:  "gemini",
				address: "gemini.conman.org",
				port:    1965,
				path:    "news.txt",
				query: [][]string{
					{"hello", "bello"},
					{"test", "test1"},
				},
			},
			wantStringRepr: "gemini://gemini.conman.org/news.txt?hello=bello&test=test1",
		},
		"missing scheme, port defualts": {
			url: "gemini.conman.org/news.txt?hello=bello&test=test1",
			wantUrl: url{
				scheme:  "gemini",
				address: "gemini.conman.org",
				port:    1965,
				path:    "news.txt",
				query: [][]string{
					{"hello", "bello"},
					{"test", "test1"},
				},
			},
			wantStringRepr: "gemini://gemini.conman.org/news.txt?hello=bello&test=test1",
		},
		"no path": {
			url: "gemini.conman.org?hello=bello&test=test1",
			wantUrl: url{
				scheme:  "gemini",
				address: "gemini.conman.org",
				port:    1965,
				path:    "",
				query: [][]string{
					{"hello", "bello"},
					{"test", "test1"},
				},
			},
			wantStringRepr: "gemini://gemini.conman.org/?hello=bello&test=test1",
		},
		"pure url": {
			url: "gemini.conman.org",
			wantUrl: url{
				scheme:  "gemini",
				address: "gemini.conman.org",
				port:    1965,
				path:    "",
				query:   [][]string{},
			},
			wantStringRepr: "gemini://gemini.conman.org/",
		},
		"scheme, address, path": {
			url: "gemini://gemini.conman.org/news.txt",
			wantUrl: url{
				scheme:  "gemini",
				address: "gemini.conman.org",
				port:    1965,
				path:    "news.txt",
				query:   [][]string{},
			},
			wantStringRepr: "gemini://gemini.conman.org/news.txt",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			url := newUrl(tc.url)
			require.NotNil(t, url)
			require.Equal(t, tc.wantUrl.scheme, url.scheme)
			require.Equal(t, tc.wantUrl.address, url.address)
			require.Equal(t, tc.wantUrl.port, url.port)
			require.Equal(t, tc.wantUrl.path, url.path)
			require.True(t, reflect.DeepEqual(tc.wantUrl.query, url.query))
			require.Equal(t, tc.wantStringRepr, tc.wantUrl.String())
		})
	}
}
