package comandante

// printableCommand holds the basic printable information about a command for the main
// help text. It is also sortable by PaddedName.
type printableCommand struct {
	PaddedName  string
	Description string
}

type printableCommands []*printableCommand

func (p printableCommands) Len() int      { return len(p) }
func (p printableCommands) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type PrintableCommandsByName struct{ printableCommands }

func (s PrintableCommandsByName) Less(i, j int) bool {
	return s.printableCommands[i].PaddedName < s.printableCommands[j].PaddedName
}
