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

// deviceDiscoveryHandler returns JSON discovery data of block devices.
func deviceDiscoveryHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	devMap := make(map[string]struct{})
	for key := range res {
		if strings.HasPrefix(key, "Disk.") {
			// Key format is MetricDefinitionId[MeasuredElementName]
			start := strings.Index(key, "[")
			end := strings.Index(key, "]")
			if start != -1 && end != -1 && end > start {
				devName := key[start+1 : end]
				devMap[devName] = struct{}{}
			}
		}
	}

	type devDiscoveryItem struct {
		DevName string `json:"{#DEVNAME}"`
	}
	var discovery []devDiscoveryItem
	for name := range devMap {
		discovery = append(discovery, devDiscoveryItem{DevName: name})
	}

	sort.Slice(discovery, func(i, j int) bool {
		return discovery[i].DevName < discovery[j].DevName
	})

	jsonRes, err := json.Marshal(discovery)
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	return string(jsonRes), nil
}

// deviceHandler returns host block devices information.
func deviceHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	dev := params["Dev"]
	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch mode {
	case "queue":
		sapMetric = "Disk.AverageQueueLength"
	case "svc":
		sapMetric = "Disk.AverageServiceTime"
	case "wait":
		sapMetric = "Disk.AverageWaitTime"
	case "io":
		sapMetric = "Disk.IORate"
	case "tput":
		sapMetric = "Disk.TotalThroughput"
	case "util":
		sapMetric = "Disk.Util"
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported mode: %s", mode), zbxerr.ErrorInvalidParams)
	}

	key := fmt.Sprintf("%s[%s]", sapMetric, dev)
	val, ok := res[key]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", key), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", key, val.Value), zbxerr.ErrorCannotFetchData)
	}

	if mode == "tput" {
		return uint64(floatVal * 1024), nil
	}

	return floatVal, nil
}
