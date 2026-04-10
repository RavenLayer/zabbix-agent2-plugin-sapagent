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
	"strings"

	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/zbxerr"
)

// networkDiscoveryHandler returns JSON discovery data of network interfaces.
func networkDiscoveryHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	ifaceMap := make(map[string]struct{})
	for key := range res {
		if strings.HasPrefix(key, "Net.") {
			// Key format is MetricDefinitionId[MeasuredElementName]
			start := strings.Index(key, "[")
			end := strings.Index(key, "]")
			if start != -1 && end != -1 && end > start {
				ifaceName := key[start+1 : end]
				ifaceMap[ifaceName] = struct{}{}
			}
		}
	}

	type netDiscoveryItem struct {
		IfaceName string `json:"{#IFNAME}"`
	}
	var discovery []netDiscoveryItem
	for name := range ifaceMap {
		discovery = append(discovery, netDiscoveryItem{IfaceName: name})
	}

	sort.Slice(discovery, func(i, j int) bool {
		return discovery[i].IfaceName < discovery[j].IfaceName
	})

	jsonRes, err := json.Marshal(discovery)
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	return string(jsonRes), nil
}

// networkInHandler returns network incoming traffic information.
func networkInHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	iface := params["Iface"]
	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch mode {
	case "bytes":
		sapMetric = "Net.KByteReceiveRate"
	case "packets":
		sapMetric = "Net.PacketReceiveRate"
	case "errors":
		sapMetric = "Net.ReceiveErrorRate"
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported traffic mode: %s", mode), zbxerr.ErrorInvalidParams)
	}

	key := fmt.Sprintf("%s[%s]", sapMetric, iface)
	val, ok := res[key]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", key), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", key, val.Value), zbxerr.ErrorCannotFetchData)
	}

	if mode == "bytes" {
		return uint64(floatVal * 1024), nil
	}

	return floatVal, nil
}

// networkOutHandler returns network outgoing traffic information.
func networkOutHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	iface := params["Iface"]
	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch mode {
	case "bytes":
		sapMetric = "Net.KByteTransmitRate"
	case "packets":
		sapMetric = "Net.PacketTransmitRate"
	case "errors":
		sapMetric = "Net.TransmitErrorRate"
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported traffic mode: %s", mode), zbxerr.ErrorInvalidParams)
	}

	key := fmt.Sprintf("%s[%s]", sapMetric, iface)
	val, ok := res[key]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", key), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", key, val.Value), zbxerr.ErrorCannotFetchData)
	}

	if mode == "bytes" {
		return uint64(floatVal * 1024), nil
	}

	return floatVal, nil
}
