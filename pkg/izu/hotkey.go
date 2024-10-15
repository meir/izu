package izu

import (
	"fmt"
	"strings"
)

// Hotkey is the type that defines a single hotkey
type Hotkey struct {
	Binding Part
	Flags   map[string][]string
	Command map[string]Part
}

func (hotkey Hotkey) String() string {
	binding := hotkey.Binding.String()

	flaglist := []string{}
	for flag, values := range hotkey.Flags {
		flaglist = append(flaglist, fmt.Sprintf("%s[%s]", flag, strings.Join(values, " ")))
	}
	flags := strings.Join(flaglist, " ")
	if flags != "" {
		flags = " | " + flags
	}

	commandlist := []string{}
	for command, parts := range hotkey.Command {
		pre := ""
		if command != "default" {
			pre = fmt.Sprintf("%s | ", command)
		}
		commandlist = append(commandlist, fmt.Sprintf("  %s%s", pre, parts.String()))
	}
	commands := strings.Join(commandlist, "\n")
	if commands != "" {
		commands = "\n" + commands
	}

	return fmt.Sprintf("%s%s%s\n", binding, flags, commands)
}
