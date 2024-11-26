package utils

import "testing"

func TestNormalizeURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"root", args{"https://blog.boot.dev/"}, "blog.boot.dev"},
		{"root without slash", args{"https://blog.boot.dev"}, "blog.boot.dev"},
		{"root http", args{"http://blog.boot.dev/"}, "blog.boot.dev"},
		{"simple", args{"https://blog.boot.dev/path/"}, "blog.boot.dev/path"},
		{"missing trailing slash", args{"https://blog.boot.dev/path"}, "blog.boot.dev/path"},
		{"simple http", args{"http://blog.boot.dev/path/"}, "blog.boot.dev/path"},
		{"missing trailing slash with http", args{"http://blog.boot.dev/path"}, "blog.boot.dev/path"},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.args.url)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tt.name, err)
				return
			}
			if got != tt.want {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tt.name, tt.want, got)
			}
		})
	}
}
