package comandante

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"text/template"
)

type Comandante struct {
	// binaryName is the name of the binary that will be used to invoke all commands
	binaryName string

	// A short description of what the binary does.
	description string

	// registeredCommands holds a list of all commands registered to this instance
	// of comandante
	registeredCommands []*Command
}

// New creates a new comandante
func New(binaryName, description string) *Comandante {
	c := &Comandante{
		binaryName:         binaryName,
		description:        description,
		registeredCommands: make([]*Command, 0),
	}

	return c
}

// RegisterCommand tells Comandante about a command so that it can be used.
func (c *Comandante) RegisterCommand(cmd *Command) error {
	// return an error if a command of the same name already exists
	for _, registeredCmd := range c.registeredCommands {
		if cmd.Name == registeredCmd.Name {
			msg := fmt.Sprintf("A command with the name '%s' already exists", cmd.Name)
			return errors.New(msg)
		}
	}

	c.registeredCommands = append(c.registeredCommands, cmd)
	return nil
}

// Run finds a command based on the command line argument and invokes
// the command in the case that one is found.
func (c *Comandante) Run() error {
	cmdName, err := getCmdName(os.Args)
	if err != nil || cmdName == "--help" || cmdName == "-h" {
		c.printDefaultHelp(os.Stderr)
		return nil
	}

	// invoke the command
	cmd := c.getCommand(cmdName)
	if cmd != nil {
		if cmd.FlagInit != nil {
			cmd.FlagInit(&cmd.flagSet)
			flag.Parse()
			cmd.flagSet.Parse(flag.Args()[1:])

			if cmd.FlagPostParse != nil {
				cmd.FlagPostParse(&cmd.flagSet)
			}
		}

		return cmd.Action()
	}

	c.printDefaultHelp(os.Stderr)
	return nil
}

// IncludeHelp adds the built in help command.
func (c *Comandante) IncludeHelp() {
	cmd := createHelpCommand(c, os.Stderr)
	c.RegisterCommand(cmd)
}

// getCommand retrieves a registered command.
func (c *Comandante) getCommand(cmdName string) *Command {
	for _, cmd := range c.registeredCommands {
		if cmdName == cmd.Name {
			return cmd
		}
	}

	return nil
}

// printDefaultHelp prints the default help text
func (c *Comandante) printDefaultHelp(w io.Writer) {
	tpl := template.New("usage")

	data := struct {
		BinaryDescription string
		BinaryName        string
		ShowHelpCommand   bool
		Commands          []*printableCommand
	}{
		c.description,
		c.binaryName,
		(c.getCommand("help") != nil),
		c.collectCommandsForHelp(),
	}

	template.Must(tpl.Parse(usage))
	_ = tpl.Execute(w, data)
}

// collectCommands
func (c *Comandante) collectCommandsForHelp() []*printableCommand {
	// find the longest command
	longest := 0
	for _, cmd := range c.registeredCommands {
		if len(cmd.Name) > longest {
			longest = len(cmd.Name)
		}
	}

	// pad all commands
	commands := make([]*printableCommand, len(c.registeredCommands))
	formatter := "%-" + strconv.Itoa(longest) + "s"
	for i, cmd := range c.registeredCommands {
		commands[i] = &printableCommand{
			PaddedName:  fmt.Sprintf(formatter, cmd.Name),
			Description: cmd.ShortDescription,
		}
	}

	sort.Sort(PrintableCommandsByName{commands})

	return commands
}

// getCmdName returns a command name from an arg list.
func getCmdName(args []string) (string, error) {
	// command name should always be the second string in the process args
	if len(args) < 2 {
		return "", errors.New("Unable to find a command")
	}

	return args[1], nil
}

var usage = `{{.BinaryDescription}}

Usage:
	{{.BinaryName}} command [arguments]

Available commands: {{ range .Commands}}
{{.PaddedName}}  {{.Description}}{{ end }}
{{if .ShowHelpCommand}}
Use "{{.BinaryName}} help [command]" for more information about a command.
{{end}}
`
