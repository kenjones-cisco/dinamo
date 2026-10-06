package version

import (
	"strings"
	"testing"
)

func TestGetVersionDisplay(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Display Version",
			want: ProductName + "\n version\t" + Version + "\n",
		},
	}
	for _, tt := range tests {
		if got := GetVersionDisplay(); !strings.HasPrefix(got, tt.want) {
			t.Errorf("%q. GetVersionDisplay() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getHumanVersion(t *testing.T) {
	previous := GitDescribe

	t.Cleanup(func() { GitDescribe = previous })

	tests := []struct {
		name     string
		describe string
		want     string
	}{
		{
			name: "Development fallback",
			want: "dev",
		},
		{
			name:     "Exact release tag",
			describe: "v0.4.0",
			want:     "v0.4.0",
		},
		{
			name:     "Quoted metadata",
			describe: "'v0.4.0'",
			want:     "v0.4.0",
		},
	}
	for _, tt := range tests {
		GitDescribe = tt.describe
		if got := getHumanVersion(); got != tt.want {
			t.Errorf("%q. getHumanVersion() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
