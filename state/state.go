package state

import (
	"context"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const iceInstanceIDPath = "/var/run/ice_instance_id"

// WriteInstanceID persists the ID of the subnitted iCE instance to
// /var/run/ice_instance_id. Failing that, due to lack of permissions, it
// writes the ID to ~/.ice_instance_id.
func WriteInstanceID(ctx context.Context, instID string) error {
	err := ioutil.WriteFile(iceInstanceIDPath, []byte(instID), 0644)
	if !os.IsPermission(err) {
		return err
	}

	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(
		filepath.Join(currentUser.HomeDir, ".ice_instance_id"),
		[]byte(instID), 0644,
	)
}
