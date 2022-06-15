// Copyright 2022 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package kafka

import (
	"flag"
	"fmt"

	"github.com/redpanda-data/console/backend/pkg/msgpack"
	"github.com/redpanda-data/console/backend/pkg/proto"
	"github.com/redpanda-data/console/backend/pkg/schema"
)

// Config required for opening a connection to Kafka
type Config struct {
	// General
	Brokers  []string `yaml:"brokers" json:"brokers"`
	ClientID string   `yaml:"clientId" json:"clientID"`
	RackID   string   `yaml:"rackId" json:"rackID"`

	// Schema Registry
	Schema      schema.Config  `yaml:"schemaRegistry" json:"schemaRegistry"`
	Protobuf    proto.Config   `yaml:"protobuf" json:"protobuf"`
	MessagePack msgpack.Config `yaml:"messagePack" json:"messagePack"`

	TLS  TLSConfig  `yaml:"tls" json:"tls"`
	SASL SASLConfig `yaml:"sasl" json:"sasl"`
}

// RegisterFlags registers all nested config flags.
func (c *Config) RegisterFlags(f *flag.FlagSet) {
	c.TLS.RegisterFlags(f)
	c.SASL.RegisterFlags(f)
	c.Protobuf.RegisterFlags(f)
	c.Schema.RegisterFlags(f)
}

// Validate the Kafka config
func (c *Config) Validate() error {
	if len(c.Brokers) == 0 {
		return fmt.Errorf("you must specify at least one broker to connect to")
	}

	err := c.Schema.Validate()
	if err != nil {
		return err
	}

	err = c.Protobuf.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate protobuf config: %w", err)
	}

	err = c.SASL.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate sasl config: %w", err)
	}

	err = c.MessagePack.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate msgpack config: %w", err)
	}

	return nil
}

// SetDefaults for Kafka config
func (c *Config) SetDefaults() {
	c.ClientID = "redpanda-console"

	c.SASL.SetDefaults()
	c.Protobuf.SetDefaults()
	c.MessagePack.SetDefaults()
}

// RedactedConfig returns a copy of the config object which redacts sensitive fields. This is useful if you
// want to log the entire config object but without sensitive fields.
func (c Config) RedactedConfig() Config {
	c.TLS.Passphrase = redactString(c.TLS.Passphrase)
	c.SASL.Password = redactString(c.SASL.Password)
	c.SASL.GSSAPIConfig.Password = redactString(c.SASL.GSSAPIConfig.Password)
	c.SASL.OAUth.Token = redactString(c.SASL.OAUth.Token)
	c.SASL.AWSMskIam.SecretKey = redactString(c.SASL.AWSMskIam.SecretKey)
	c.SASL.AWSMskIam.SessionToken = redactString(c.SASL.AWSMskIam.SessionToken)

	return c
}

func redactString(in string) string {
	if in == "" {
		return in
	}
	return "<redacted>"
}
