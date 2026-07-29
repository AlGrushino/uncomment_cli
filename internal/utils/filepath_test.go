package utils

import "testing"

func TestGetFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantRes string
		wantErr bool
	}{
		{"regular", "/some/new/path", "/some/new/", false},
		{"directory, no file", "/some/new/", "", false},
		{"no source directory", "some/new/path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GetFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf(`GetFilePath(%s) error = %v, wantErr = %v`, tt.path, err, tt.wantErr)
			}
			if res != tt.wantRes {
				t.Errorf(`GetFilePath(%s) = %q, want match for %#q`, tt.path, res, tt.wantRes)
			}
		})
	}
}
