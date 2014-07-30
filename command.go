package comandante

import (
	"flag"
)

type actionFunc func() error
type flagInitFunc func(*flag.FlagSet)
type flagPostParseFunc func(*flag.FlagSet)

// Command details a command that can be run from the command line
type Command struct {
	// Name is used for invoking the command. If Name is "sayhello" the command is
	// invoked with: your-binary sayhello
	Name string

	// ShortDescription is a one line description of the command that is displayed in
	// the default list of commands
	ShortDescription string

	// Documentation is a longer form of help text that is displayed when help
	// information is queried about a specific command.
	Documentation string

	// Action is the function called when the command is invoked
	Action actionFunc

	// FlagInit is the function called to handle delaing with flags sent to the command
	FlagInit flagInitFunc

	// FlagPostParse is the function called after the flagset has been parsed
	FlagPostParse flagPostParseFunc

	// flagset is for handling command lines flags passed into the command
	flagSet flag.FlagSet
}

// NewCommand creates a new command with a name, a short description and an action that runs
// when the command is invoked.
func NewCommand(name, shortDescription string, action actionFunc) *Command {
	cmd := &Command{
		Name:             name,
		ShortDescription: shortDescription,
		Action:           action,
	}

	return cmd
}

// parseFlags parses the flagset for the command
func (c *Command) parseFlags(args []string) error {
	return c.flagSet.Parse(args)
}
