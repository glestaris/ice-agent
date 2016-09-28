package commands

import (
	"context"

	"github.com/glestaris/ice-agent/ice"
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
		cli.StringSliceFlag{
			Name:  "tag",
			Usage: "Tags to include in the registered instance",
		},
	},

	Action: func(ctx *cli.Context) error {
		var err error

		inst := ice.Instance{
			SessionID: ctx.String("session-id"),
			Tags:      ctx.StringSlice("tag"),
		}

		// SSH
		inst.SSHUsername, err = ssh.Username(context.TODO())
		if err != nil {
			return cli.NewExitError("ERROR", 1)
		}
		inst.SSHAuthorizedFingerprint, err = ssh.AuthorizedFingerprint(context.TODO(), inst.SSHUsername)
		if err != nil {
			return cli.NewExitError("ERROR", 1)
		}

		// Write the fake instance ID
		if err := state.WriteInstanceID(context.TODO(), "fake-inst-id"); err != nil {
			return cli.NewExitError("ERROR", 1)
		}

		return nil
	},
}
