/*
 * config_test.go
 * Created on 25.02.2020
 * Copyright (C) 2020 Volkswagen AG, All rights reserved
 *
 */

package server

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindConfig(t *testing.T) {
	t.Parallel()

	type input struct {
		args []string
	}

	type output struct {
		assert func(assert *assert.Assertions, cfg *Config)
	}

	type test struct {
		input  input
		output output
	}

	cases := map[string]test{
		"happy path": {
			input: input{
				args: []string{
					"--server.contextPath=test",
					"--profiling.enabled",
				},
			},
			output: output{
				assert: func(assert *assert.Assertions, cfg *Config) {
					assert.True(cfg.ProfilingEnabled)
					assert.Equal("test", cfg.ContextPath)
					assert.Equal(8080, cfg.Port)
				},
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			cfg := Config{}

			BindConfig(&cfg, fs)

			err := fs.Parse(tc.input.args)
			assert.NoError(err)

			tc.output.assert(assert, &cfg)
		})
	}
}

func TestConfig_GetContextPath(t *testing.T) {

	t.Parallel()

	type test struct {
		cfg      *Config
		expected string
	}

	cases := map[string]test{
		"empty": {
			cfg: &Config{
				ContextPath: "",
			},
			expected: "",
		},
		"just a slash": {
			cfg: &Config{
				ContextPath: "/",
			},
			expected: "",
		},
		"slash tss slash": {
			cfg: &Config{
				ContextPath: "/tss/",
			},
			expected: "tss",
		},
	}

	for name, tc := range cases {

		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := tc.cfg.GetContextPath()
			assert.Equal(t, tc.expected, result)
		})
	}
}
