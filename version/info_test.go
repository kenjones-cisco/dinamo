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
	GitDescribe = "e42813d"

	tests := []struct {
		name string
		want string
	}{
		{
			name: "Git Variables defined",
			want: GitDescribe,
		},
	}
	for _, tt := range tests {
		if got := getHumanVersion(); got != tt.want {
			t.Errorf("%q. getHumanVersion() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
