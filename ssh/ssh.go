package ssh

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/glestaris/passwduser"

	"golang.org/x/crypto/ssh"
)

// Username returns the username of the running user.
func Username(ctx context.Context) (string, error) {
	currentUser, err := passwduser.Current()
	if err != nil {
		return "", err
	}

	return currentUser.Username, nil
}

func firstAuthorizedKeysLine(authorizedKeysPath string) ([]byte, error) {
	authorizedKeysFile, err := os.Open(authorizedKeysPath)
	if err != nil {
		return nil, err
	}

	authorizedKeysReader := bufio.NewReader(authorizedKeysFile)
	line, isPrefix, err := authorizedKeysReader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("file `%s` is empty", authorizedKeysPath)
		}

		return nil, err
	}
	if isPrefix {
		return nil, errors.New("first line of ~/.ssh/authorized_keys is too long")
	}

	if err = authorizedKeysFile.Close(); err != nil {
		return nil, err
	}

	return line, nil
}

func firstAuthorizedKey(usr *passwduser.User) (ssh.PublicKey, error) {
	authorizedKeysPath := filepath.Join(usr.HomeDir, ".ssh", "authorized_keys")
	line, err := firstAuthorizedKeysLine(authorizedKeysPath)
	if err != nil {
		return nil, err
	}

	pk, _, _, _, err := ssh.ParseAuthorizedKey(line)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func fingerprint(pk ssh.PublicKey) string {
	fingerprintHash := md5.Sum(pk.Marshal())

	fingerprintHashHex := make([]byte, 47)

	hex.Encode(fingerprintHashHex, fingerprintHash[:])
	fingerprint := ""
	for i := 0; i < 15; i++ {
		c := i * 2
		fingerprint += string(fingerprintHashHex[c : c+2])
		fingerprint += ":"
	}
	fingerprint += string(fingerprintHashHex[30:32])

	return fingerprint
}

// AuthorizedFingerprint returns the MD5 fingerprint of the first public SSH
// key defined in ~/.ssh/authorized_keys.
func AuthorizedFingerprint(ctx context.Context, sshUsername string) (string, error) {
	currentUser, err := passwduser.Current()
	if err != nil {
		return "", err
	}

	pk, err := firstAuthorizedKey(currentUser)
	if err != nil {
		return "", err
	}

	return fingerprint(pk), nil
}
