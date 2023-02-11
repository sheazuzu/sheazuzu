/*
 *  command_test.go
 *  Created on 01.09.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package cli

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCommand_Execute(t *testing.T) {

	type test struct {
		cmd            *Command
		args           []string
		prgName        string
		expectedOutput string
		exitCode       int
	}

	cases := map[string]test{
		"execute with no args and subcommands": {
			cmd: &Command{
				Name:    "test",
				Usage:   "test usage",
				Version: "1.2.3",
				Flags:   flag.NewFlagSet("", flag.ContinueOnError),
				Validate: func() bool {
					return true
				},
				Run: func(command *Command, args ...string) {
					_, _ = fmt.Fprint(out, "Executed")
				},
			},
			expectedOutput: "Executed",
		},
		"execute subcommand": {
			args: []string{"test"},
			cmd: &Command{
				Name:    "",
				Usage:   "test usage",
				Version: "1.2.3",
				SubCommands: []Command{
					{
						Name:  "test",
						Usage: "test usage",
						Flags: flag.NewFlagSet("", flag.ContinueOnError),
						Validate: func() bool {
							return true
						},
						Run: func(command *Command, args ...string) {
							_, _ = fmt.Fprint(out, "Executed")
						},
					},
				},
			},
			expectedOutput: "Executed",
		},
		"execute run instead of subcommand": {
			args: []string{"foo"},
			cmd: &Command{
				Name:    "",
				Usage:   "test usage",
				Version: "1.2.3",
				Flags:   flag.NewFlagSet("", flag.ContinueOnError),
				SubCommands: []Command{
					{
						Name:  "test",
						Usage: "test usage",
					},
				},
				Validate: func() bool {
					return true
				},
				Run: func(command *Command, args ...string) {
					_, _ = fmt.Fprint(out, "Executed")
				},
			},
			expectedOutput: "Executed",
		},
		"execute unknown subcommand": {
			args: []string{"foo"},
			cmd: &Command{
				Name:    "",
				Usage:   "test usage",
				Version: "1.2.3",
				Flags:   flag.NewFlagSet("", flag.ContinueOnError),
				SubCommands: []Command{
					{
						Name:  "test",
						Usage: "test usage",
						Flags: flag.NewFlagSet("", flag.ContinueOnError),
					},
				},
			},
			expectedOutput: "Unknown command.\n\nUSAGE\n   [command] [flags]\n\n  test usage\n\nVERSION\n  1.2.3\n\nCOMMANDS\n  test                            test usage\n\nFLAGS\n",
		},
		"execute with no run": {
			cmd: &Command{
				Name:    "test",
				Usage:   "test usage",
				Version: "1.2.3",
				Flags:   flag.NewFlagSet("", flag.ContinueOnError),
				Validate: func() bool {
					return true
				},
				Run: nil,
			},
			expectedOutput: "\nUSAGE\n   test [flags]\n\n  test usage\n\nVERSION\n  1.2.3\n\nFLAGS\n",
		},
		"execute --help with subcommands and grouped fs": {
			args: []string{"--help"},
			cmd: &Command{
				Name:    "test",
				Usage:   "test usage",
				Version: "1.2.3",
				Commit:  "ababab",
				Flags: func() *flag.FlagSet {
					fs := flag.NewFlagSet("", flag.ContinueOnError)
					fs.String("foo.bar", "baz", "foobar")
					return fs
				}(),
				Validate: func() bool {
					return true
				},
				SubCommands: []Command{
					{
						Name:        "subtest",
						Usage:       "sub test",
						Flags:       nil,
						SubCommands: nil,
						Run:         nil,
						Validate:    nil,
					},
				},
				Run: func(command *Command, args ...string) {},
			},
			expectedOutput: "\nUSAGE\n   test [flags]\n\n  test usage\n\nVERSION\n  1.2.3\n\nCOMMIT\n  ababab\n\nCOMMANDS\n  subtest                         sub test\n\nFLAGS\n\n  FOO\n  --foo.bar                                      foobar (default 'baz')\n",
		},
		"execute --help with subcommands and ungrouped fs": {
			args: []string{"--help"},
			cmd: &Command{
				Name:    "",
				Usage:   "test usage",
				Version: "1.2.3",
				Commit:  "ababab",
				Flags: func() *flag.FlagSet {
					fs := flag.NewFlagSet("", flag.ContinueOnError)
					fs.String("foo", "baz", "foobar")
					return fs
				}(),
				Validate: func() bool {
					return true
				},
				SubCommands: []Command{
					{
						Name:        "subtest",
						Usage:       "sub test",
						Flags:       nil,
						SubCommands: nil,
						Run:         nil,
						Validate:    nil,
					},
				},
				Run: func(command *Command, args ...string) {},
			},
			expectedOutput: "\nUSAGE\n   [command] [flags]\n\n  test usage\n\nVERSION\n  1.2.3\n\nCOMMIT\n  ababab\n\nCOMMANDS\n  subtest                         sub test\n\nFLAGS\n  --foo                                          foobar (default 'baz')\n",
		},
		"execute --version": {
			args: []string{"--version"},
			cmd: &Command{
				Name:    "test",
				Usage:   "test usage",
				Version: "1.2.3",
				Flags:   flag.NewFlagSet("", flag.ContinueOnError),
				Validate: func() bool {
					return true
				},
				Run: func(command *Command, args ...string) {},
			},
			expectedOutput: "test usage\n\nVERSION\n1.2.3\n\nCOMMIT\n\n",
		},
	}

	for name, tc := range cases {

		tc := tc

		t.Run(name, func(t *testing.T) {

			// we can replace the os.Exit call with our own function to catch the exist and check the result exit code
			var resultCode int
			exit = func(code int) {
				resultCode = code
			}
			defer func() {
				exit = os.Exit
			}()

			prgName = tc.prgName

			// capture the output in a buffer
			reader, writer, _ := os.Pipe()
			out = writer
			buf := new(bytes.Buffer)

			// execute the method under test
			tc.cmd.Execute(tc.args...)

			// reset program name
			prgName = os.Args[0]

			// reset stdout and convert the captured output ito a string
			_ = writer.Close()
			_, _ = buf.ReadFrom(reader)
			out = os.Stdout
			capturedOutput := buf.String()

			assert.Equal(t, tc.expectedOutput, capturedOutput)
			assert.Equal(t, tc.exitCode, resultCode)
		})
	}
}

func TestCommand_parseAndValidate(t *testing.T) {

	type test struct {
		args     []string
		cmd      *Command
		exitCode int
	}

	cases := map[string]test{
		"valid command": {
			cmd: &Command{
				Name:    "test",
				Usage:   "test usage",
				Version: "1.2.3",
				Flags:   flag.NewFlagSet("", flag.ContinueOnError),
				Validate: func() bool {
					return true
				},
			},
			exitCode: 0,
		},
		"parse failed": {
			args: []string{
				"--foo",
			},
			cmd: &Command{
				Name:    "test",
				Usage:   "test usage",
				Version: "1.2.3",
				Flags:   flag.NewFlagSet("", flag.ContinueOnError),
				Validate: func() bool {
					return true
				},
			},
			exitCode: 1,
		},
		"validate failed": {
			cmd: &Command{
				Name:    "test",
				Usage:   "test usage",
				Version: "1.2.3",
				Flags:   flag.NewFlagSet("", flag.ContinueOnError),
				Validate: func() bool {
					return false
				},
			},
			exitCode: 1,
		},
	}

	for name, tc := range cases {

		tc := tc

		t.Run(name, func(t *testing.T) {

			// we can replace the os.Exit call with our own function to catch the exist and check the result exit code
			var resultCode int
			exit = func(code int) {
				resultCode = code
			}
			defer func() {
				exit = os.Exit
			}()

			tc.cmd.parseAndValidate(tc.args...)

			assert.Equal(t, tc.exitCode, resultCode)
		})
	}
}
