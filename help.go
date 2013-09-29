package comandante

import (
	"fmt"
	"io"
	"os"
)

func createHelpCommand(com *comandante, w io.Writer) *Command {
	action := func() error {
		// The first parameter passed to the help command should be the
		// command for requested documentation.
		if len(os.Args) < 3 {
			com.printDefaultHelp(w)
		} else {
			cmdName := os.Args[2]
			cmd := com.getCommand(cmdName)

			if cmd != nil && cmdName != "help" {
				fmt.Fprintf(w, "%s\n\n%s\n", cmdName, cmd.Documentation)

				if cmd.FlagInit != nil {
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
