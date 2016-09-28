package ssh

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestUsername(t *testing.T) {
	username, err := Username(nil)
	if err != nil {
		t.Fatal(err)
	}
	if username != "ice" {
		t.Fatalf("Username `%s` does not match `ice`", username)
	}
}

func populateAuthorizedKeys() error {
	if err := os.Mkdir("/home/ice/.ssh", 0700); err != nil {
		return err
	}

	f, err := os.OpenFile("/home/ice/.ssh/authorized_keys", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	testAPub, err := ioutil.ReadFile("../test_assets/test_1.pub")
	if err != nil {
		return err
	}
	if _, err = f.Write(testAPub); err != nil {
		return err
	}

	testBPub, err := ioutil.ReadFile("../test_assets/test_1.pub")
	if err != nil {
		return err
	}
	if _, err := f.Write(testBPub); err != nil {
		return err
	}

	return f.Close()
}

func TestAuthorizedFingerprint(t *testing.T) {
	if err := populateAuthorizedKeys(); err != nil {
		t.Fatalf("Failed to populate the ~/.ssh/authorized_keys file: %s", err)
	}

	fingerprint, err := AuthorizedFingerprint(nil, "ice")
	if err != nil {
		t.Fatal(err)
	}

	expectedFingerprint := "30:b6:cb:7e:0b:a3:5a:56:b2:f2:c7:c3:16:1d:2f:db"
	if fingerprint != expectedFingerprint {
		t.Fatalf(
			"Wrong fingerprint for test_assets/test_1.pub: expected `%s`, got `%s`",
			expectedFingerprint, fingerprint,
		)
	}
}
