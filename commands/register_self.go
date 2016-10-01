package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/glestaris/ice-agent/ice"
	"github.com/glestaris/ice-agent/network"
	"github.com/glestaris/ice-agent/ssh"
	"github.com/glestaris/ice-agent/state"
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
	},

	Action: func(ctx *cli.Context) error {
		apiEndpoint := ctx.String("api-endpoint")
		sessionID := ctx.String("session-id")

		inst := ice.Instance{
			SessionID: sessionID,
		}

		// SSH
		var err error
		inst.SSHUsername, err = ssh.Username(context.TODO())
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
		}
		inst.SSHAuthorizedFingerprint, err = ssh.AuthorizedFingerprint(
			context.TODO(), inst.SSHUsername,
		)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
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
		inst.PublicReverseDNS, err = network.ReverseDNS(
			context.TODO(), inst.PublicIPAddr.String(),
		)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("ERROR: %s", err), 1)
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
