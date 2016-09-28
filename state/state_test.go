package state

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestWriteInstanceID(t *testing.T) {
	if err := WriteInstanceID(nil, "hello-world"); err != nil {
		t.Fatal(err)
	}

	stateFileContents, err := ioutil.ReadFile("/home/ice/.ice_instance_id")
	if err != nil {
		t.Fatal(err)
	}

	actualInstanceID := strings.TrimSpace(string(stateFileContents))
	if actualInstanceID != "hello-world" {
		t.Fatalf(
			"Expected to read instance ID `hello-world`, instead read `%s`",
			actualInstanceID,
		)
	}
}
