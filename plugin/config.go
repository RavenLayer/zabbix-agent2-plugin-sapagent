/*
 * Copyright (c) 2026 RavenLayer <sasha@ravenlayer.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugin

import (
	"golang.zabbix.com/sdk/conf"
	"golang.zabbix.com/sdk/plugin"
)

// Session struct holds individual options for SAP Agent connection for each session.
type Session struct {
	// Uri is a connection string consisting of a network scheme, a host address and a port or a path to a Unix-socket.
	Uri string `conf:"optional"`

	// Username of SAP Agent host.
	Username string `conf:"optional"`

	// Password of SAP Agent host.
	Password string `conf:"optional"`

	// TLSRootCAFile is the path to the TLS CA file.
	TLSRootCAFile string `conf:"optional"`

	// TLSServerName is the TLS ServerName to verify the hostname (SNI).
	TLSServerName string `conf:"optional"`

	// TLSInsecureSkipVerify controls whether a client verifies the server's certificate chain and host name.
	TLSInsecureSkipVerify string `conf:"optional"`
}

// PluginOptions are options for SAP Agent connection.
type PluginOptions struct {
	plugin.SystemOptions `conf:"optional,name=System"`

	// Timeout is the maximum time in seconds for waiting when a HTTP connection has to be established.
	// Default value equals to the global agent timeout.
	Timeout int `conf:"optional,range=1:30"`

	// CacheTTL is the maximum time in seconds for keeping cached metric values in memory.
	CacheTTL int `conf:"optional"`

	// Sessions stores pre-defined named sets of connections settings.
	Sessions map[string]Session `conf:"optional"`

	// Default stores default connection parameter values from configuration file
	Default Session `conf:"optional"`
}

var DefaultOptions = PluginOptions{}

// Configure implements the Configurator interface.
// Initializes configuration structures.
func (p *Plugin) Configure(global *plugin.GlobalOptions, options interface{}) {
	// Start with defaults
	p.options = DefaultOptions

	if err := conf.Unmarshal(options, &p.options); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err)
	}

	if p.options.Timeout == 0 {
		p.options.Timeout = global.Timeout
	}

	if p.options.CacheTTL == 0 {
		p.options.CacheTTL = 60
	}
}

// Validate implements the Configurator interface.
// Returns an error if validation of a plugin's configuration is failed.
func (p *Plugin) Validate(options interface{}) error {
	var opts PluginOptions

	return conf.Unmarshal(options, &opts)
}
