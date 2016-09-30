package ssh

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
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

func makeAuthorizedKeysFile() (io.WriteCloser, error) {
	if err := os.Mkdir("/home/ice/.ssh", 0700); err != nil {
		return nil, err
	}

	return os.OpenFile("/home/ice/.ssh/authorized_keys", os.O_CREATE|os.O_WRONLY, 0600)
}

func populateAuthorizedKeys() error {
	f, err := makeAuthorizedKeysFile()
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
		t.Fatalf("Failed to populate the `/home/.ssh/authorized_keys` file: %s", err)
	}
	defer func() {
		if err := os.RemoveAll("/home/ice/.ssh"); err != nil {
			t.Fatalf("Failed to remove `/home/ice/.ssh`: %s", err)
		}
	}()

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

func TestAuthorizedFingerprintReturnsErrorWhenNoKeysFileFound(t *testing.T) {
	_, err := AuthorizedFingerprint(nil, "ice")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "open /home/ice/.ssh/authorized_keys: no such file"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf(
			"Expected error message to contain `%s`, got `%s`",
			expectedError, err.Error(),
		)
	}
}

func TestAuthorizedFingerprintReturnsErrorWhenNoKeysDefined(t *testing.T) {
	f, err := makeAuthorizedKeysFile()
	if err != nil {
		t.Fatalf("Failed to make the `/home/.ssh/authorized_keys` file: %s", err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll("/home/ice/.ssh"); err != nil {
			t.Fatalf("Failed to remove `/home/ice/.ssh`: %s", err)
		}
	}()

	_, err = AuthorizedFingerprint(nil, "ice")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "file `/home/ice/.ssh/authorized_keys` is empty"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf(
			"Expected error message to contain `%s`, got `%s`",
			expectedError, err.Error(),
		)
	}
}

func TestAuthorizedFingerprintReturnsErrorWhenKeysFileIsInvalid(t *testing.T) {
	f, err := makeAuthorizedKeysFile()
	if err != nil {
		t.Fatalf("Failed to make the `/home/.ssh/authorized_keys` file: %s", err)
	}
	_, err = f.Write([]byte("hello world 123\n"))
	if err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll("/home/ice/.ssh"); err != nil {
			t.Fatalf("Failed to remove `/home/ice/.ssh`: %s", err)
		}
	}()

	_, err = AuthorizedFingerprint(nil, "ice")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "no key found"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf(
			"Expected error message to contain `%s`, got `%s`",
			expectedError, err.Error(),
		)
	}
}
