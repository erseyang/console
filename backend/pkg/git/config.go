// Copyright 2022 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package git

import (
	"flag"
	"fmt"
	"time"
)

// Config for Git Service
type Config struct {
	Enabled bool `yaml:"enabled" json:"enabled"`

	// AllowedFileExtensions specifies file extensions that shall be picked up. If at least one is specified all other
	// file extensions will be ignored.
	AllowedFileExtensions []string `yaml:"-" json:"-"`

	// Max file size which will be considered. Files exceeding this size will be ignored and logged.
	MaxFileSize int64 `yaml:"-" json:"-"`

	// Whether or not to use the filename or the full filepath as key in the map
	IndexByFullFilepath bool `yaml:"-" json:"-"`

	// RefreshInterval specifies how often the repository shall be pulled to check for new changes.
	RefreshInterval time.Duration `yaml:"refreshInterval" json:"refreshInterval"`

	// Repository that contains markdown files that document a Kafka topic.
	Repository RepositoryConfig `yaml:"repository" json:"repository"`

	// Authentication Configs
	BasicAuth BasicAuthConfig `yaml:"basicAuth" json:"basicAuth"`
	SSH       SSHConfig       `yaml:"ssh" json:"ssh"`
}

// RegisterFlagsWithPrefix for all (sub)configs
func (c *Config) RegisterFlagsWithPrefix(f *flag.FlagSet, prefix string) {
	c.BasicAuth.RegisterFlagsWithPrefix(f, prefix)
	c.SSH.RegisterFlagsWithPrefix(f, prefix)
}

// Validate all root and child config structs
func (c *Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.RefreshInterval == 0 {
		return fmt.Errorf("git config is enabled but refresh interval is set to 0 (disabled)")
	}

	return c.Repository.Validate()
}

// SetDefaults for all root and child config structs
func (c *Config) SetDefaults() {
	c.Repository.SetDefaults()

	c.RefreshInterval = time.Minute
	c.MaxFileSize = 500 * 1000 // 500KB
	c.IndexByFullFilepath = false
}
