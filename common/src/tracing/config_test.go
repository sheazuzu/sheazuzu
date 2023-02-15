/*
 *  config_test.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package tracing

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindConfig(t *testing.T) {
	t.Parallel()

	type input struct {
		args        []string
		defaultName string
	}

	type output struct {
		enabled      bool
		serviceName  string
		agentURI     string
		collectorURI string
	}

	type test struct {
		input  input
		output output
	}

	cases := map[string]test{
		"everything set": {
			input: input{
				args: []string{
					"--tracing.enabled",
					"--tracing.serviceName",
					"testService",
					"--tracing.agent.uri",
					"testAgentUri",
					"--tracing.collector.uri",
					"testCollectorUri",
				},
			},
			output: output{
				serviceName:  "testService",
				agentURI:     "testAgentUri",
				collectorURI: "testCollectorUri",
				enabled:      true,
			},
		},
		"nothing set": {
			input: input{
				args:        []string{},
				defaultName: "common",
			},
			output: output{
				serviceName: "common",
				enabled:     false,
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := flag.NewFlagSet("", flag.ContinueOnError)
			cfg := &Config{}

			BindConfig(cfg, fs, tc.input.defaultName)

			err := fs.Parse(tc.input.args)
			assert.NoError(t, err)

			assert.Equal(t, tc.output.enabled, cfg.Enabled)
			assert.Equal(t, tc.output.serviceName, cfg.ServiceName, "wrong service name")
			assert.Equal(t, tc.output.agentURI, cfg.AgentURI, "wrong agent uri")
			assert.Equal(t, tc.output.collectorURI, cfg.CollectorURI, "wrong collector uri")
		})
	}
}

func TestConfig_IsValid(t *testing.T) {
	t.Parallel()

	type input struct {
		cfg *Config
	}

	type output struct {
		expected bool
	}

	type test struct {
		input  input
		output output
	}

	cases := map[string]test{
		"everything set": {
			input: input{
				cfg: &Config{
					Enabled:      true,
					AgentURI:     "http://test.com",
					CollectorURI: "http://test.com",
				},
			},
			output: output{
				expected: true,
			},
		},
		"nothing set": {
			input: input{
				cfg: &Config{},
			},
			output: output{
				expected: true,
			},
		},
		"enabled, but no agent or collector": {
			input: input{
				cfg: &Config{
					Enabled: true,
				},
			},
			output: output{
				expected: false,
			},
		},
		"invalid collector url": {
			input: input{
				cfg: &Config{
					Enabled:      true,
					CollectorURI: "invalid",
				},
			},
			output: output{
				expected: false,
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := tc.input.cfg.IsValid()

			assert.Equal(t, tc.output.expected, result)
		})
	}
}
