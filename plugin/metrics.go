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
	"context"
	"fmt"
	"strings"

	"golang.zabbix.com/sdk/metric"
	"golang.zabbix.com/sdk/plugin"
	"golang.zabbix.com/sdk/uri"
)

const (
	keyPing              = "sap.agent.ping"
	keyVersion           = "sap.agent.version"
	keyCPUDiscovery      = "sap.agent.cpu.discovery"
	keyCPUUtil           = "sap.agent.cpu.util"
	keyCPULoad           = "sap.agent.cpu.load"
	keyCPUNum            = "sap.agent.cpu.num"
	keyCPUIntr           = "sap.agent.cpu.intr"
	keyCPUSwitches       = "sap.agent.cpu.switches"
	keyMem               = "sap.agent.mem"
	keySwapSize          = "sap.agent.swap.size"
	keySwapIn            = "sap.agent.swap.in"
	keySwapOut           = "sap.agent.swap.out"
	keyNetDiscovery      = "sap.agent.net.discovery"
	keyNetIn             = "sap.agent.net.in"
	keyNetOut            = "sap.agent.net.out"
	keyFsDiscovery       = "sap.agent.fs.discovery"
	keyFsSize            = "sap.agent.fs.size"
	keyDevDiscovery      = "sap.agent.dev.discovery"
	keyDev               = "sap.agent.dev"
	keyOsName            = "sap.agent.os.name"
	keyOsVersion         = "sap.agent.os.version"
	keyDatabaseDiscovery = "sap.agent.db.discovery"
	keyProcDiscovery     = "sap.agent.proc.discovery"
	keyProcGet           = "sap.agent.proc.get"

	uriParam                   = "Uri"
	usernameParam              = "Username"
	passwordParam              = "Password"
	tlsRootCAFileParam         = "TLSRootCAFile"
	tlsServerNameParam         = "TLSServerName"
	tlsInsecureSkipVerifyParam = "TLSInsecureSkipVerify"
)

var uriDefaults = &uri.Defaults{Scheme: "http", Port: "1128"}

var (
	maxPasswordLen = 512
)

var (
	// ConnParams
	paramURI = metric.NewConnParam(uriParam, "URI to connect or session name.").
			WithDefault(uriDefaults.Scheme + "://localhost:" + uriDefaults.Port).WithSession().
			WithValidator(AgentURIValidator{
			Defaults:       uriDefaults,
			AllowedSchemes: []string{"http", "https", "unix"},
		})
	paramUsername = metric.NewConnParam(usernameParam, "SAP Agent username.").WithDefault("sapadm")
	paramPassword = metric.NewConnParam(passwordParam, "SAP Agent password.").
			WithDefault("").
			WithValidator(metric.LenValidator{Max: &maxPasswordLen})
	// SessionOnlyParams
	paramTLSCaFile             = metric.NewSessionOnlyParam(tlsRootCAFileParam, "TLS CA file path.").WithDefault("")
	paramTLSServerName         = metric.NewSessionOnlyParam(tlsServerNameParam, "TLS ServerName to verify the hostname (SNI).").WithDefault("")
	paramTLSInsecureSkipVerify = metric.NewSessionOnlyParam(tlsInsecureSkipVerifyParam, "Controls whether a client verifies the server's certificate chain and host name.").WithDefault("")
)

var metrics = metric.MetricSet{
	keyPing: metric.New(
		"Returns whether the connection to the SAP Agent can be established.",
		getParameters(nil), true,
	),
	keyVersion: metric.New(
		"Returns SAP Agent version information.",
		getParameters(nil), false,
	),
	keyCPUDiscovery: metric.New(
		"Returns JSON discovery data of logical CPUs.",
		getParameters(nil), false,
	),
	keyCPUUtil: metric.New(
		"Returns CPU utilization percentage.",
		getParameters(
			&AdditionalParam{metric.NewParam("CPU", "CPU number or all.").WithDefault("all"), 3},
			&AdditionalParam{metric.NewParam("Type", "Utilization type (user, system, idle, iowait, steal).").WithValidator(StringEnumValidator{AllowedValues: []string{"user", "system", "idle", "iowait", "steal"}}).SetRequired(), 4},
		), false,
	),
	keyCPULoad: metric.New(
		"Returns CPU load average.",
		getParameters(
			&AdditionalParam{metric.NewParam("Mode", "Load average mode (avg1, avg5, avg15).").WithDefault("avg1"), 3},
		), false,
	),
	keyCPUNum: metric.New(
		"Returns total number of logical CPUs.",
		getParameters(nil), false,
	),
	keyCPUIntr: metric.New(
		"Returns interrupt rate per second.",
		getParameters(nil), false,
	),
	keyCPUSwitches: metric.New(
		"Returns context switch rate per second.",
		getParameters(nil), false,
	),
	keyMem: metric.New(
		"Returns host memory information.",
		getParameters(
			&AdditionalParam{metric.NewParam("Type", "Memory type (phy, virt).").WithValidator(StringEnumValidator{AllowedValues: []string{"phy", "virt"}}).SetRequired(), 3},
			&AdditionalParam{metric.NewParam("Mode", "Memory mode (available, free, size).").WithValidator(StringEnumValidator{AllowedValues: []string{"available", "free", "size"}}).SetRequired(), 4},
		), false,
	),
	keySwapSize: metric.New(
		"Returns swap size information.",
		getParameters(
			&AdditionalParam{metric.NewParam("Mode", "Swap mode (free, used, total).").WithValidator(StringEnumValidator{AllowedValues: []string{"free", "used", "total"}}).SetRequired(), 3},
		), false,
	),
	keySwapIn: metric.New(
		"Returns swap-in rate information.",
		getParameters(
			&AdditionalParam{metric.NewParam("Mode", "Swap mode (bytes, pages).").WithValidator(StringEnumValidator{AllowedValues: []string{"bytes", "pages"}}).SetRequired(), 3},
		), false,
	),
	keySwapOut: metric.New(
		"Returns swap-out rate information.",
		getParameters(
			&AdditionalParam{metric.NewParam("Mode", "Swap mode (bytes, pages).").WithValidator(StringEnumValidator{AllowedValues: []string{"bytes", "pages"}}).SetRequired(), 3},
		), false,
	),
	keyNetDiscovery: metric.New(
		"Returns JSON discovery data of network interfaces.",
		getParameters(nil), false,
	),
	keyNetIn: metric.New(
		"Returns network incoming traffic information.",
		getParameters(
			&AdditionalParam{metric.NewParam("Iface", "Network interface name.").SetRequired(), 3},
			&AdditionalParam{metric.NewParam("Mode", "Traffic mode (bytes, packets, errors).").WithValidator(StringEnumValidator{AllowedValues: []string{"bytes", "packets", "errors"}}).SetRequired(), 4},
		), false,
	),
	keyNetOut: metric.New(
		"Returns network outgoing traffic information.",
		getParameters(
			&AdditionalParam{metric.NewParam("Iface", "Network interface name.").SetRequired(), 3},
			&AdditionalParam{metric.NewParam("Mode", "Traffic mode (bytes, packets, errors).").WithValidator(StringEnumValidator{AllowedValues: []string{"bytes", "packets", "errors"}}).SetRequired(), 4},
		), false,
	),
	keyFsDiscovery: metric.New(
		"Returns JSON discovery data of mounted filesystem.",
		getParameters(nil), false,
	),
	keyFsSize: metric.New(
		"Returns host filesystem information.",
		getParameters(
			&AdditionalParam{metric.NewParam("FS", "Filesystem name (mount point).").SetRequired(), 3},
			&AdditionalParam{metric.NewParam("Mode", "Size mode (total, free).").WithValidator(StringEnumValidator{AllowedValues: []string{"total", "free"}}).SetRequired(), 4},
		), false,
	),
	keyDevDiscovery: metric.New(
		"Returns JSON discovery data of block devices.",
		getParameters(nil), false,
	),
	keyDev: metric.New(
		"Returns host block devices information.",
		getParameters(
			&AdditionalParam{metric.NewParam("Dev", "Block device name.").SetRequired(), 3},
			&AdditionalParam{metric.NewParam("Mode", "Device mode (queue, svc, wait, io, tput, util).").WithValidator(StringEnumValidator{AllowedValues: []string{"queue", "svc", "wait", "io", "tput", "util"}}).SetRequired(), 4},
		), false,
	),
	keyOsName: metric.New(
		"Returns host name information.",
		getParameters(nil), false,
	),
	keyOsVersion: metric.New(
		"Returns host version information.",
		getParameters(nil), false,
	),
	keyDatabaseDiscovery: metric.New(
		"Returns JSON discovery data of database instances.",
		getParameters(nil), false,
	),
	keyProcDiscovery: metric.New(
		"Returns JSON discovery data of instance processes.",
		getParameters(nil), false,
	),
	keyProcGet: metric.New(
		"Returns JSON data of a specific instance process.",
		getParameters(
			&AdditionalParam{metric.NewParam("Instance", "Instance name.").SetRequired(), 3},
			&AdditionalParam{metric.NewParam("Process", "Process name.").SetRequired(), 4},
		), false,
	),
}

type AgentURIValidator struct {
	Defaults       *uri.Defaults
	AllowedSchemes []string
}

type StringEnumValidator struct {
	AllowedValues []string
}

// handlerFunc defines an interface that must be implemented by handlers.
type handlerFunc func(ctx context.Context,
	client *SAPClient,
	params map[string]string,
	extraParams ...string) (res interface{}, err error)

type AdditionalParam struct {
	param    *metric.Param
	position int
}

func init() {
	// Register with the internal Agent 2 manager
	err := plugin.RegisterMetrics(&Impl, Name, metrics.List()...)
	if err != nil {
		panic(err)
	}
}

// getHandlerFunc returns a handlerFunc related to a given key.
func getHandlerFunc(key string) handlerFunc {
	switch key {
	case keyVersion:
		return versionHandler
	case keyPing:
		return pingHandler
	case keyCPUDiscovery:
		return cpuDiscoveryHandler
	case keyCPUUtil:
		return cpuUtilHandler
	case keyCPULoad:
		return cpuLoadHandler
	case keyCPUNum:
		return cpuNumHandler
	case keyCPUIntr:
		return cpuIntrHandler
	case keyCPUSwitches:
		return cpuSwitchesHandler
	case keyMem:
		return memoryHandler
	case keySwapSize:
		return swapSizeHandler
	case keySwapIn:
		return swapInHandler
	case keySwapOut:
		return swapOutHandler
	case keyNetDiscovery:
		return networkDiscoveryHandler
	case keyNetIn:
		return networkInHandler
	case keyNetOut:
		return networkOutHandler
	case keyFsDiscovery:
		return filesystemDiscoveryHandler
	case keyFsSize:
		return filesystemSizeHandler
	case keyDevDiscovery:
		return deviceDiscoveryHandler
	case keyDev:
		return deviceHandler
	case keyOsName:
		return osNameHandler
	case keyOsVersion:
		return osVersionHandler
	case keyDatabaseDiscovery:
		return databaseDiscoveryHandler
	case keyProcDiscovery:
		return procDiscoveryHandler
	case keyProcGet:
		return procGetHandler
	default:
		return nil
	}
}

func (v AgentURIValidator) Validate(value *string) error {
	if value == nil {
		return nil
	}

	u, err := uri.New(*value, v.Defaults)
	if err != nil {
		return err
	}

	isValidScheme := false
	if v.AllowedSchemes != nil {
		for _, s := range v.AllowedSchemes {
			if u.Scheme() == s {
				isValidScheme = true
				break
			}
		}

		if !isValidScheme {
			return fmt.Errorf("allowed schemes: %s", strings.Join(v.AllowedSchemes, ", "))
		}
	}

	return nil
}

func (v StringEnumValidator) Validate(value *string) error {
	isValidValue := false

	if v.AllowedValues != nil {
		for _, s := range v.AllowedValues {
			if *value == s {
				isValidValue = true
				break
			}
		}

		if !isValidValue {
			return fmt.Errorf("allowed values: %s", strings.Join(v.AllowedValues, ", "))
		}
	}

	return nil
}

func getParameters(adds ...*AdditionalParam) []*metric.Param {
	m := []*metric.Param{
		paramURI,
		paramUsername,
		paramPassword,
		paramTLSCaFile,
		paramTLSServerName,
		paramTLSInsecureSkipVerify,
	}

	for _, add := range adds {
		if add != nil && add.param != nil {
			m = append(m[:add.position+1], m[add.position:]...)
			m[add.position] = add.param
		}
	}

	return m
}
