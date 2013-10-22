package comandante

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestComandanteNew(t *testing.T) {
	c := New("binaryName", "binaryDescription")

	// the commandante should have binaryName and binaryDescription set
	if c.binaryName != "binaryName" {
		t.Error("comandante binaryName is incorrect")
	}
	if c.description != "binaryDescription" {
		t.Error("comandante description is incorrect")
	}

	// the length of registered commands should be 0
	if len(c.registeredCommands) != 0 {
		t.Error("registeredCommands should initially be zero")
	}
}

func TestRunCommand(t *testing.T) {
	c := New("binaryName", "")

	a := false
	cmd := NewCommand("test", "short description", func() error { a = true; return nil })
	c.RegisterCommand(cmd)

	oldArgs := os.Args
	os.Args = []string{"bin", "test"}

	c.Run()

	if a != true {
		t.Error("a should be true. It looks like the command didn't run")
	}

	os.Args = oldArgs
}

func TestRegisterCommand(t *testing.T) {
	c := New("binaryName", "")

	cmd := NewCommand("test", "short description", func() error { return nil })
	c.RegisterCommand(cmd)

	// one command should be registered
	if len(c.registeredCommands) != 1 {
		t.Error("registeredCommands should be 1")
	}

	// adding a second command with the same name should fail
	err := c.RegisterCommand(cmd)
	if err == nil {
		t.Error("RegisterCommand should fail when adding two commands of the same name")
	}
}

func TestCommandArgs(t *testing.T) {
	c := New("binaryName", "")

	testFlag := ""
	cmd := NewCommand("test", "short description", func() error { return nil })
	cmd.FlagInit = func(fs *flag.FlagSet) {
		fs.StringVar(&testFlag, "testing", "", "This is the usage")
	}
	cmd.FlagInit(&cmd.flagSet)
	c.RegisterCommand(cmd)
	cmd.parseFlags([]string{"--testing=value"})

	if testFlag != "value" {
		t.Error("parsing a sub command flag did not work", testFlag)
	}
}

func TestDefaultHelp(t *testing.T) {
	c := New("binaryName", "Handles things")
	b := bytes.NewBufferString("")

	cmd := createHelpCommand(c, b)
	c.RegisterCommand(cmd)

	othercmd := NewCommand("other", "does stuff", func() error { return nil })
	c.RegisterCommand(othercmd)

	c.printDefaultHelp(b)

	if !strings.Contains(b.String(), "other") {
		t.Error("Missing the 'other' command")
	}

	if !strings.Contains(b.String(), "help") {
		t.Error("Missing the 'help' command")
	}

	if !strings.Contains(b.String(), "more information about a command") {
		t.Error("Missing extra help text")
	}
}

func TestHelpCommand(t *testing.T) {
	c := New("binaryName", "")
	b := bytes.NewBufferString("")

	cmd := createHelpCommand(c, b)
	c.RegisterCommand(cmd)

	othercmd := NewCommand("other", "short", func() error { return nil })
	doc := "The documentation"
	othercmd.Documentation = doc
	c.RegisterCommand(othercmd)

	// the command should be found
	oldArgs := os.Args
	os.Args = []string{"bin", "help", "other"}
	found := c.getCommand("help")
	if found == nil {
		t.Error("Help command was not registered")
	}

	// running the command should write help text to our buffer
	found.Action()
	if !strings.Contains(b.String(), doc) {
		t.Error("printing documentation from the help command did not work")
	}
	os.Args = oldArgs
}
