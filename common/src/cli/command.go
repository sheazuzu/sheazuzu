/*
 *  command.go
 *  Created on 29.08.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package cli

import (
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Command struct {
	Name        string
	Usage       string
	Version     string
	Commit      string
	Flags       *flag.FlagSet
	SubCommands []Command
	Run         func(*Command, ...string)
	Validate    func() bool
}

// we use this, to manipulate the values for testing
var prgName = os.Args[0]
var exit = os.Exit
var out = os.Stdout

// performs setup and and validation and execution of the command
func (cmd *Command) Execute(args ...string) {

	if cmd.Flags != nil {
		cmd.Flags.SetOutput(ioutil.Discard)
	}

	if len(args) == 0 {
		if cmd.Run != nil {
			cmd.parseAndValidate(args...)
			cmd.Run(cmd)
			return
		}

		_, _ = fmt.Fprintf(out, "\nUSAGE\n%s", cmd.Help())
		return
	}

	if args[0] == "--help" {
		_, _ = fmt.Fprintf(out, "\nUSAGE\n%s", cmd.Help())
		return
	}

	if args[0] == "--version" {
		_, _ = fmt.Println(cmd.Name)
		if cmd.Usage != "" {
			_, _ = fmt.Fprintln(out, cmd.Usage)
		}
		_, _ = fmt.Fprintln(out)

		_, _ = fmt.Fprintln(out, "VERSION")
		_, _ = fmt.Fprintln(out, cmd.Version)
		_, _ = fmt.Fprintln(out)

		_, _ = fmt.Fprintln(out, "COMMIT")
		_, _ = fmt.Fprintln(out, cmd.Commit)
		return
	}

	for _, subCmd := range cmd.SubCommands {
		if subCmd.Name == args[0] {
			subCmd.Version = cmd.Version
			subCmd.Commit = cmd.Commit
			subCmd.Execute(args[1:]...)
			return
		}
	}

	if cmd.Run != nil {
		cmd.parseAndValidate(args...)
		cmd.Run(cmd, args...)
		return
	}

	_, _ = fmt.Fprintf(out, "Unknown command.\n\nUSAGE\n%s", cmd.Help())
}

func (cmd *Command) PrintUsage() string {
	return fmt.Sprintf("%s  %s%s", cmd.Name, strings.Repeat(" ", 30-len(cmd.Name)), cmd.Usage)
}

func (cmd *Command) Help() string {
	var msg strings.Builder

	subCommand := cmd.Name
	if subCommand == "" && len(cmd.SubCommands) > 0 {
		subCommand = "[command]"
	}

	_, binaryName := path.Split(strings.ReplaceAll(prgName, "\\", "/"))

	_, _ = fmt.Fprintf(&msg, "  %s %s [flags]\n\n", binaryName, subCommand)

	if cmd.Usage != "" {
		_, _ = fmt.Fprintf(&msg, "  %s\n", cmd.Usage)
		_, _ = fmt.Fprintln(&msg)
	}

	if cmd.Version != "" {
		_, _ = fmt.Fprintln(&msg, "VERSION")
		_, _ = fmt.Fprintf(&msg, "  %s\n", cmd.Version)
		_, _ = fmt.Fprintln(&msg)
	}

	if cmd.Commit != "" {
		_, _ = fmt.Fprintln(&msg, "COMMIT")
		_, _ = fmt.Fprintf(&msg, "  %s\n", cmd.Commit)
		_, _ = fmt.Fprintln(&msg)
	}

	if len(cmd.SubCommands) > 0 {
		_, _ = fmt.Fprintln(&msg, "COMMANDS")

		for _, cmd := range cmd.SubCommands {
			_, _ = fmt.Fprintln(&msg, "  "+cmd.PrintUsage())
		}
		_, _ = fmt.Fprintln(&msg)
	}

	if cmd.Flags != nil {
		_, _ = fmt.Fprintln(&msg, "FLAGS")

		// grouped flags are appended at the end of the message to keep the ungrouped flags together
		var groupedFlagsMsg strings.Builder
		caption := ""

		cmd.Flags.VisitAll(func(f *flag.Flag) {

			nameParts := strings.Split(f.Name, ".")
			if len(nameParts) > 1 {
				if caption != nameParts[0] {
					_, _ = fmt.Fprintf(&groupedFlagsMsg, "\n  %s\n", strings.ToUpper(nameParts[0]))
				}

				caption = nameParts[0]
				_, _ = fmt.Fprintf(&groupedFlagsMsg, "  --%s%s%s (default '%s')\n", f.Name, strings.Repeat(" ", 45-len(f.Name)), f.Usage, f.DefValue)

				return
			}

			caption = ""
			_, _ = fmt.Fprintf(&msg, "  --%s%s%s (default '%s')\n", f.Name, strings.Repeat(" ", 45-len(f.Name)), f.Usage, f.DefValue)
		})

		msg.WriteString(groupedFlagsMsg.String())
	}

	return msg.String()
}

func (cmd *Command) parseAndValidate(args ...string) {

	err := ff.Parse(cmd.Flags, args, ff.WithEnvVarNoPrefix())
	if err != nil {
		_, _ = fmt.Fprintln(out, "startup failed")
		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintln(out, err.Error())
		_, _ = fmt.Fprintln(out, "\nUSAGE\n\n"+cmd.Help())
		exit(1)
		return
	}

	if cmd.Validate != nil && !cmd.Validate() {
		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintln(out, "\nUSAGE\n\n"+cmd.Help())
		exit(1)
		return
	}
}
