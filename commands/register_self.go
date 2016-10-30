package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ice-stuff/ice-agent/ice"
	"github.com/ice-stuff/ice-agent/network"
	"github.com/ice-stuff/ice-agent/ssh"
	"github.com/ice-stuff/ice-agent/state"
	"github.com/urfave/cli"
)

// RegisterSelfCommand implements the register-self subcommand of the agent.
var RegisterSelfCommand = cli.Command{
	Name:        "register-self",
	Usage:       "register-self [options]",
	Description: "Registers virtual instance that calls this script to iCE.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "api-endpoint",
			Usage: "The iCE REST API endpoint URL",
		},
		cli.StringFlag{
			Name:  "session-id",
			Usage: "The session id",
		},
		cli.StringSliceFlag{
			Name:  "tag",
			Usage: "An instnace tag in the form of <key>=<value>",
		},
	},

	Action: func(ctx *cli.Context) error {
		apiEndpoint := ctx.String("api-endpoint")
		sessionID := ctx.String("session-id")
		tagsList := ctx.StringSlice("tag")

		tags, err := parseTags(tagsList)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
		}

		inst := ice.Instance{
			SessionID: sessionID,
			Tags:      tags,
		}

		// SSH
		sshAuthorizedFingerprint, err := ssh.AuthorizedFingerprint(context.TODO())
		if err == nil {
			inst.SSHAuthorizedFingerprint = sshAuthorizedFingerprint
		}

		iceClient := ice.NewClient(apiEndpoint)
		// Network
		inst.Networks, err = network.Networks(context.TODO())
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
		}
		inst.PublicIPAddr, err = iceClient.MyIP(context.TODO())
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
		}
		publicReverseDNS, err := network.ReverseDNS(
			context.TODO(), inst.PublicIPAddr.String(),
		)
		if err == nil {
			inst.PublicReverseDNS = publicReverseDNS
		}

		// Store the instance over
		instID, err := iceClient.StoreInstance(context.TODO(), inst)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
		}

		// Write the instance ID
		if err := state.WriteInstanceID(context.TODO(), instID); err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
		}

		// Print output
		inst.ID = instID
		if err := json.NewEncoder(os.Stdout).Encode(inst); err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
		}

		return nil
	},
}

func parseTags(tagsList []string) (map[string]string, error) {
	m := make(map[string]string)

	for _, tag := range tagsList {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag provided `%s`", tag)
		}

		m[parts[0]] = parts[1]
	}

	return m, nil
}
