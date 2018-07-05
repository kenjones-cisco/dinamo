package commands

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var wantConfig = `apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-alln.cisco.com:443
  name: cae-alln-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-rcdn.cisco.com:443
  name: cae-rcdn-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-rtp.cisco.com:443
  name: cae-rtp-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://localhost:443
  name: localhost:443
contexts:
- context:
    cluster: cae-alln-cisco-com:443
    namespace: mynamespace
    user: fakeuser/cae-alln-cisco-com:443
  name: mynamespace/cae-alln-cisco-com:443/fakeuser
- context:
    cluster: cae-rcdn-cisco-com:443
    namespace: mynamespace
    user: fakeuser/cae-rcdn-cisco-com:443
  name: mynamespace/cae-rcdn-cisco-com:443/fakeuser
- context:
    cluster: cae-rtp-cisco-com:443
    namespace: mynamespace
    user: fakeuser/cae-rtp-cisco-com:443
  name: mynamespace/cae-rtp-cisco-com:443/fakeuser
- context:
    cluster: localhost:443
    namespace: mynamespace
    user: fakeuser/localhost:443
  name: mynamespace/localhost:443/fakeuser
`

func TestGenerate_error1(t *testing.T) {
	template := "fixtures/config.tmpl"

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("generate --template %s", template))
	if result.Error == nil {
		t.Error("Expected command to fail due to missing --file flag")
	}

	want := `required flag(s) "file" not set`
	assertResult(t, want, result.Error.Error())
}

func TestGenerate_error2(t *testing.T) {
	file := makeTmp("config.yaml")

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("generate --file %s", file))
	if result.Error == nil {
		t.Error("Expected command to fail due to missing --template flag")
	}

	want := `required flag(s) "template" not set`
	assertResult(t, want, result.Error.Error())
}

func TestGenerate_error3(t *testing.T) {
	template := "fixtures/config.tmpl"
	file := makeTmp("config.yaml")

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("generate --template %s --file %s", template, file))
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg(s)")
	}

	want := `template: config.tmpl:17:15: executing "config.tmpl" at <.INSTANCE>: map has no entry for key "INSTANCE"`
	assertResult(t, want, result.Error.Error())
}

func TestGenerate(t *testing.T) {
	template := "fixtures/config.tmpl"
	file := makeTmp("config.yaml")
	args := []string{"INSTANCE=https://localhost:443", "INSTANCE_NAME=localhost:443", "NAMESPACE=mynamespace", "USERNAME=fakeuser"}

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("generate --template %s --file %s %s", template, file, strings.Join(args, " ")))
	if result.Error != nil {
		t.Error(result.Error)
	}

	got := readFile(file)
	assertResult(t, wantConfig, got)
	want := "generated: " + file + "\n"
	assertResult(t, want, result.Output)
}

func TestGenerate_alias_short(t *testing.T) {
	template := "fixtures/config.tmpl"
	file := makeTmp("config.yaml")
	args := []string{"INSTANCE=https://localhost:443", "INSTANCE_NAME=localhost:443", "NAMESPACE=mynamespace", "USERNAME=fakeuser"}

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("gen -t %s -f %s %s", template, file, strings.Join(args, " ")))
	if result.Error != nil {
		t.Error(result.Error)
	}

	got := readFile(file)
	assertResult(t, wantConfig, got)
	want := "generated: " + file + "\n"
	assertResult(t, want, result.Output)
}

func TestGenerateEnv(t *testing.T) {
	envVars := []string{"INSTANCE", "INSTANCE_NAME", "NAMESPACE", "USERNAME"}
	existingEnv := make(map[string]string)
	for _, k := range envVars {
		existingEnv[k] = os.Getenv(k)
	}
	defer func() {
		for k, v := range existingEnv {
			_ = os.Setenv(k, v)
		}
	}()

	template := "fixtures/config.tmpl"
	file := makeTmp("config.yaml")

	_ = os.Setenv("INSTANCE", "https://localhost:443")
	_ = os.Setenv("INSTANCE_NAME", "localhost:443")
	_ = os.Setenv("NAMESPACE", "mynamespace")
	_ = os.Setenv("USERNAME", "fakeuser")

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("generate --env --template %s --file %s", template, file))
	if result.Error != nil {
		t.Error(result.Error)
	}

	got := readFile(file)
	assertResult(t, wantConfig, got)
	want := "generated: " + file + "\n"
	assertResult(t, want, result.Output)
}

func TestGenerateEnv_alias_short(t *testing.T) {
	envVars := []string{"INSTANCE", "INSTANCE_NAME", "NAMESPACE", "USERNAME"}
	existingEnv := make(map[string]string)
	for _, k := range envVars {
		existingEnv[k] = os.Getenv(k)
	}
	defer func() {
		for k, v := range existingEnv {
			_ = os.Setenv(k, v)
		}
	}()

	template := "fixtures/config.tmpl"
	file := makeTmp("config.yaml")

	_ = os.Setenv("INSTANCE", "https://localhost:443")
	_ = os.Setenv("INSTANCE_NAME", "localhost:443")
	_ = os.Setenv("NAMESPACE", "mynamespace")
	_ = os.Setenv("USERNAME", "fakeuser")

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("gen -e -t %s -f %s", template, file))
	if result.Error != nil {
		t.Error(result.Error)
	}

	got := readFile(file)
	assertResult(t, wantConfig, got)
	want := "generated: " + file + "\n"
	assertResult(t, want, result.Output)
}
