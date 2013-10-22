package comandante

import (
	"fmt"
	"io"
	"os"
)

func createHelpCommand(com *Comandante, w io.Writer) *Command {
	action := func() error {
		// The first parameter passed to the help command should be the
		// command for requested documentation.
		if len(os.Args) < 3 {
			com.printDefaultHelp(w)
		} else {
			cmdName := os.Args[2]
			cmd := com.getCommand(cmdName)

			if cmd != nil && cmdName != "help" {
				fmt.Fprintf(w, "%s %s\n%s", com.binaryName, cmdName, cmd.Documentation)

				if cmd.FlagInit != nil {
					cmd.FlagInit(&cmd.flagSet)

					fmt.Fprintf(w, "\noptions\n")

					cmd.flagSet.SetOutput(w)
					cmd.flagSet.PrintDefaults()
				}
			} else {
				com.printDefaultHelp(w)
			}
		}

		return nil
	}

	cmd := NewCommand("help", "get more information about a command", action)
	return cmd
}
