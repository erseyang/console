// Copyright 2022 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package kafka

import "flag"

// TLSConfig to connect to Kafka via TLS
type TLSConfig struct {
	Enabled               bool   `yaml:"enabled" json:"enabled"`
	CaFilepath            string `yaml:"caFilepath" json:"caFilepath"`
	CertFilepath          string `yaml:"certFilepath" json:"certFilepath"`
	KeyFilepath           string `yaml:"keyFilepath" json:"keyFilepath"`
	Passphrase            string `yaml:"passphrase" json:"passphrase"`
	InsecureSkipTLSVerify bool   `yaml:"insecureSkipTlsVerify" json:"insecureSkipTLSVerify"`
}

// RegisterFlags for all sensitive Kafka TLS configs
func (c *TLSConfig) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&c.Passphrase, "kafka.tls.passphrase", "", "Passphrase to optionally decrypt the private key")
}
