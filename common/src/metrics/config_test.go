/*
 * config_test.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
 */

package metrics

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
		serviceName string
	}

	type test struct {
		input  input
		output output
	}

	cases := map[string]test{
		"everything set": {
			input: input{
				args: []string{
					"--metrics.enabled",
					"--metrics.serviceName",
					"testService",
				},
			},
			output: output{
				serviceName: "testService",
			},
		},
		"nothing set": {
			input: input{
				args:        []string{},
				defaultName: "common",
			},
			output: output{
				serviceName: "common",
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

			BindConfig(&cfg, fs, tc.input.defaultName)

			err := fs.Parse(tc.input.args)
			assert.NoError(err)

			assert.True(cfg.Enabled)
			assert.Equal(tc.output.serviceName, cfg.ServiceName, "wrong service name")
		})
	}
}
