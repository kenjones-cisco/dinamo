package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type resulter struct {
	Error   error
	Output  string
	Command *cobra.Command
}

func runCmd(c *cobra.Command, input string) resulter {
	buf := new(bytes.Buffer)
	c.SetOutput(buf)
	c.SetArgs(strings.Split(input, " "))

	err := c.Execute()
	output := buf.String()

	return resulter{err, output, c}
}

func getRootCommand() *cobra.Command {
	return NewCommandCLI()
}

func assertResult(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if expectedValue != actualValue {
		t.Error("Expected <", expectedValue, "> but got <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func uniqueID() string {
	b := make([]byte, 36)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func makeTmp(name string) string {
	tmp := path.Join(os.TempDir(), uniqueID())
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return path.Join(os.TempDir(), name)
	}
	return path.Join(tmp, name)
}

func readFile(name string) string {
	rawData, err := ioutil.ReadFile(name)
	if err != nil {
		return ""
	}
	return string(rawData)
}
