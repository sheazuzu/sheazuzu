/*
 * config_test.go Created on 13.12.2021Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package mongo

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
		err    bool
		assert func(cfg *Config)
	}

	cases := map[string]struct {
		input  input
		output output
	}{
		"happy path": {
			input: input{
				args: []string{
					"--mongo.uri", "mongodb://test.com",
					"--mongo.database", "testDB",
					"--mongo.timeout", "20",
					"--mongo.useSSL",
					"--mongo.sslClientCertFile", "myCertificate",
					"--mongo.sslClientKeyFile", "myKey",
				},
			},
			output: output{
				err: false,
				assert: func(cfg *Config) {
					assert.EqualValues(t, "mongodb://test.com", cfg.URI)
					assert.EqualValues(t, "testDB", cfg.Database)
					assert.EqualValues(t, 20, cfg.Timeout)
					assert.EqualValues(t, true, cfg.UseSSL)
					assert.EqualValues(t, "myCertificate", cfg.SSLClientCertFile)
					assert.EqualValues(t, "myKey", cfg.SSLClientKeyFile)
				},
			},
		},
		"error - unknown variable": {
			input: input{
				args: []string{
					"--unknownParameter", "test",
				},
			},
			output: output{
				err: true,
				assert: func(cfg *Config) {
					assert.Empty(t, cfg.URI)
				},
			},
		},
	}

	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			config := &Config{}
			fs := flag.NewFlagSet("", flag.ContinueOnError)

			BindConfig(config, fs)

			err := fs.Parse(tc.input.args)
			if tc.output.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tc.output.assert(config)
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

	cases := map[string]struct {
		input  input
		output output
	}{
		"happy path": {
			input: input{
				cfg: Config{
					URI:               "mongodb://test.com",
					Database:          "test",
					Timeout:           10,
					UseSSL:            true,
					SSLClientCertFile: "myCertificate",
					SSLClientKeyFile:  "myKey",
				},
			},
			output: output{
				isValid: true,
			},
		},
		"error - no uri": {
			input: input{
				cfg: Config{
					URI:               "",
					Database:          "test",
					Timeout:           10,
					UseSSL:            true,
					SSLClientCertFile: "myCertificate",
					SSLClientKeyFile:  "myKey",
				},
			},
			output: output{
				isValid: false,
			},
		},
		"error - no database": {
			input: input{
				cfg: Config{
					URI:               "mongodb://test.com",
					Database:          "",
					Timeout:           10,
					UseSSL:            true,
					SSLClientCertFile: "myCertificate",
					SSLClientKeyFile:  "myKey",
				},
			},
			output: output{
				isValid: false,
			},
		},
		"error - no ssl cert": {
			input: input{
				cfg: Config{
					URI:               "mongodb://test.com",
					Database:          "test",
					Timeout:           10,
					UseSSL:            true,
					SSLClientCertFile: "",
					SSLClientKeyFile:  "myKey",
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

			isValid := tc.input.cfg.IsValid()
			assert.EqualValues(t, tc.output.isValid, isValid)
		})
	}
}
