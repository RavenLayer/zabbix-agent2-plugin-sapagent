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
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/zbxerr"
)

// newClientFn allows mocking SOAPClient creation in unit tests
var newClientFn = NewClient

// procDiscoveryHandler returns JSON discovery data of SAP instance processes.
func procDiscoveryHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	instances, err := client.ListInstances()
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	tlsCaFile := params[tlsRootCAFileParam]
	tlsServerName := params[tlsServerNameParam]
	tlsInsecureSkipVerify, _ := strconv.ParseBool(params[tlsInsecureSkipVerifyParam])

	type procDiscoveryItem struct {
		InstanceName string `json:"{#INSTANCE.NAME}"`
		ProcessName  string `json:"{#PROCESS.NAME}"`
		ProcessDesc  string `json:"{#PROCESS.DESC}"`
	}

	var discovery []procDiscoveryItem

	for _, inst := range instances {
		instanceURL, err := BuildInstanceURL(params[uriParam], inst.MSystemNumber)
		if err != nil {
			Impl.Debugf("failed to build instance URL for system %s: %s", inst.MSystemNumber, err.Error())
			continue
		}

		instClient := newClientFn(instanceURL, client.HTTPClient.Timeout, tlsInsecureSkipVerify)
		instClient.SetBasicAuth(client.Username, client.Password)
		_ = instClient.SetTLSConfig(tlsCaFile, tlsServerName, tlsInsecureSkipVerify)

		instSAPClient := &SAPClient{
			SOAPClient: instClient,
			cache:      client.cache,
		}

		cachedProcesses, err := instSAPClient.GetProcessList()
		if err != nil {
			Impl.Debugf("failed to fetch process list for instance %s: %s", inst.MSid, err.Error())
			continue
		}

		for _, procVal := range cachedProcesses {
			var proc GetProcessListItem
			if err := json.Unmarshal([]byte(procVal.Value), &proc); err == nil {
				discovery = append(discovery, procDiscoveryItem{
					InstanceName: inst.MSid,
					ProcessName:  proc.Name,
					ProcessDesc:  proc.Description,
				})
			}
		}
	}

	sort.Slice(discovery, func(i, j int) bool {
		if discovery[i].InstanceName != discovery[j].InstanceName {
			return discovery[i].InstanceName < discovery[j].InstanceName
		}
		return discovery[i].ProcessName < discovery[j].ProcessName
	})

	jsonRes, err := json.Marshal(discovery)
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	return string(jsonRes), nil
}

// procGetHandler returns the JSON-serialized metrics for a specific process of a specific instance.
func procGetHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	instanceName := params["Instance"]
	processName := params["Process"]

	if instanceName == "" || processName == "" {
		return nil, errs.WrapConst(fmt.Errorf("instance name and process name are required"), zbxerr.ErrorInvalidParams)
	}

	instances, err := client.ListInstances()
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var targetInstance *SapInstance
	for _, inst := range instances {
		if inst.MSid == instanceName {
			targetInstance = &inst
			break
		}
	}

	if targetInstance == nil {
		return nil, errs.WrapConst(fmt.Errorf("instance %s not found", instanceName), zbxerr.ErrorCannotFetchData)
	}

	tlsCaFile := params[tlsRootCAFileParam]
	tlsServerName := params[tlsServerNameParam]
	tlsInsecureSkipVerify, _ := strconv.ParseBool(params[tlsInsecureSkipVerifyParam])

	instanceURL, err := BuildInstanceURL(params[uriParam], targetInstance.MSystemNumber)
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	instClient := newClientFn(instanceURL, client.HTTPClient.Timeout, tlsInsecureSkipVerify)
	instClient.SetBasicAuth(client.Username, client.Password)
	_ = instClient.SetTLSConfig(tlsCaFile, tlsServerName, tlsInsecureSkipVerify)

	instSAPClient := &SAPClient{
		SOAPClient: instClient,
		cache:      client.cache,
	}

	cachedProcesses, err := instSAPClient.GetProcessList()
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	procVal, ok := cachedProcesses[processName]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("process %s not found in instance %s", processName, instanceName), zbxerr.ErrorCannotFetchData)
	}

	return procVal.Value, nil
}
