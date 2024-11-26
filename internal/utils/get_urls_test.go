package utils

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	type args struct {
		htmlBody   string
		rawBaseURL string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"absolute and relative urls",
			args{
				`<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>`,
				"https://blog.boot.dev",
			},
			[]string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
			false,
		},

		{"no urls",
			args{
				`<html>
	<body>
		<p>Hello, World!</p>
	</body>
</html>`,
				"https://blog.boot.dev",
			},
			[]string{},
			false,
		},

		{"complex A tag",
			args{
				`<html>
	<body>
		<a id="jo" href="../python/getting-started">
			<span>Python getting started</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="/">
			<span>Root</span>
		</a>
	</body>
</html>`,
				"https://blog.boot.dev/tutorials/go/",
			},
			[]string{"https://blog.boot.dev/tutorials/python/getting-started", "https://other.com/path/one", "https://blog.boot.dev/"},
			false,
		},

		{"relative urls",
			args{
				`<html>
	<body>
		<a href="other">
			<span>Other</span>
		</a>
		<a href="..">
			<span>up</span>
		</a>
		<a href="../down">
			<span>up and down</span>
		</a>
		<a href="./here">
			<span>here</span>
		</a>
	</body>
</html>`,
				"https://blog.boot.dev/playground/go/",
			},
			[]string{"https://blog.boot.dev/playground/go/other", "https://blog.boot.dev/playground/", "https://blog.boot.dev/playground/down", "https://blog.boot.dev/playground/go/here"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseUrl, err := url.Parse(tt.args.rawBaseURL)
			if err != nil {
				t.Errorf("GetURLsFromHTML() invalid base url %v in test setup: %v", tt.args.rawBaseURL, err)
				return
			}
			got, err := GetURLsFromHTML(tt.args.htmlBody, baseUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetURLsFromHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetURLsFromHTML() got = %v, want %v", got, tt.want)
			}
		})
	}
}
