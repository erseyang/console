// Copyright 2022 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package connect

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
	Enabled        bool            `yaml:"enabled" json:"enabled"`
	Clusters       []ConfigCluster `yaml:"clusters" json:"clusters"`
	ConnectTimeout time.Duration   `yaml:"connectTimeout" json:"connectTimeout"` // used for connectivity test
	ReadTimeout    time.Duration   `yaml:"readTimeout" json:"readTimeout"`       // overall REST/HTTP read timeout
	RequestTimeout time.Duration   `yaml:"requestTimeout" json:"requestTimeout"` // timeout for REST requests to Kafka Connect
}

func (c *Config) SetDefaults() {
	for _, cluster := range c.Clusters {
		cluster.SetDefaults()
	}
	c.ConnectTimeout = 15 * time.Second
	c.ReadTimeout = 60 * time.Second
	c.RequestTimeout = 6 * time.Second
}

// RegisterFlags registers all nested config flags.
func (c *Config) RegisterFlags(f *flag.FlagSet) {
	for i, cluster := range c.Clusters {
		flagNamePrefix := fmt.Sprintf("connect.clusters.%d.", i)
		cluster.RegisterFlagsWithPrefix(f, flagNamePrefix)
	}
}

func (c *Config) Validate() error {
	for i, cluster := range c.Clusters {
		err := cluster.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate cluster at index '%d' (name: '%v'): %w", i, cluster.Name, err)
		}
	}
	return nil
}
