/*
 *  config_test.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package logging

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindConfig(t *testing.T) {
	t.Parallel()

	type input struct {
		args []string
	}

	type output struct {
		assert func(t *testing.T, cfg *Config)
	}

	type test struct {
		input  input
		output output
	}

	cases := map[string]test{
		"testLevel and Format": {
			input: input{
				args: []string{
					"--logging.level",
					"testLevel",
					"--logging.format",
					"testFormat",
				},
			},
			output: output{
				assert: func(t *testing.T, cfg *Config) {
					t.Helper()
					assert.Equal(t, "testFormat", cfg.Format, "wrong format")
					assert.Equal(t, "testLevel", cfg.Level, "wrong level")
				},
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := flag.NewFlagSet("", flag.ContinueOnError)
			cfg := Config{}

			BindConfig(&cfg, fs)

			err := fs.Parse(tc.input.args)
			assert.NoError(t, err)

			tc.output.assert(t, &cfg)
		})
	}
}

func TestConfig_IsValid(t *testing.T) {
	t.Parallel()

	type input struct {
		cfg Config
	}

	type output struct {
		isValid bool
	}

	type test struct {
		input  input
		output output
	}

	cases := map[string]test{
		"level info": {
			input: input{
				cfg: Config{
					Level:  "info",
					Format: "json",
				},
			},
			output: output{
				isValid: true,
			},
		},
		"level debug": {
			input: input{
				cfg: Config{
					Level:  "debug",
					Format: "json",
				},
			},
			output: output{
				isValid: true,
			},
		},
		"level warn": {
			input: input{
				cfg: Config{
					Level:  "warn",
					Format: "json",
				},
			},
			output: output{
				isValid: true,
			},
		},
		"level error": {
			input: input{
				cfg: Config{
					Level:  "error",
					Format: "json",
				},
			},
			output: output{
				isValid: true,
			},
		},
		"format console": {
			input: input{
				cfg: Config{
					Level:  "info",
					Format: "console",
				},
			},
			output: output{
				isValid: true,
			},
		},
		"wrong level": {
			input: input{
				cfg: Config{
					Level:  "thisleveldoesnotexist",
					Format: "json",
				},
			},
			output: output{
				isValid: false,
			},
		},
		"wrong format": {
			input: input{
				cfg: Config{
					Level:  "info",
					Format: "thisformatdoesnotexist",
				},
			},
			output: output{
				isValid: false,
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := tc.input.cfg.IsValid()
			assert.Equal(t, tc.output.isValid, result)
		})
	}
}

func ExampleConfig_IsValid_valid() {
	cfg := Config{
		Level:  "info",
		Format: "console",
	}
	fmt.Println(cfg.IsValid())
	// Output: true
}

func ExampleConfig_IsValid_invalidFormat() {
	cfg := Config{
		Level:  "info",
		Format: "thisformatdoesnotexist",
	}
	fmt.Println(cfg.IsValid())
	// Output: log format must either be 'json' 'console', or 'minimal'
	//false
}

func ExampleConfig_IsValid_invalidLevel() {
	cfg := Config{
		Level:  "everything",
		Format: "console",
	}
	fmt.Println(cfg.IsValid())
	// Output: log level must either be 'info', 'debug', 'warn' or 'error'
	//false
}
