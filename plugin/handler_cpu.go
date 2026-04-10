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

const (
	metricProcPrefix = "Proc."
	metricOSPrefix   = "OS."
)

// cpuDiscoveryHandler returns JSON discovery data of logical CPUs.
func cpuDiscoveryHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	cpuMap := make(map[string]struct{})
	for key := range res {
		if strings.HasPrefix(key, metricProcPrefix) {
			if strings.HasSuffix(key, "]") {
				start := strings.Index(key, "[")
				if start != -1 {
					cpuMap[key[start+1:len(key)-1]] = struct{}{}
				}
			}
		}
	}

	type cpuDiscoveryItem struct {
		ID string `json:"{#CPU.ID}"`
	}
	var discovery []cpuDiscoveryItem
	for id := range cpuMap {
		discovery = append(discovery, cpuDiscoveryItem{ID: id})
	}

	sort.Slice(discovery, func(i, j int) bool {
		return discovery[i].ID < discovery[j].ID
	})

	jsonRes, err := json.Marshal(discovery)
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	return string(jsonRes), nil
}

// cpuUtilHandler returns CPU utilization percentage.
func cpuUtilHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	cpu := params["CPU"]
	utilType := params["Type"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch utilType {
	case "user":
		sapMetric = "UserTimePercentage"
	case "system":
		sapMetric = "SystemTimePercentage"
	case "idle":
		sapMetric = "IdleTimePercentage"
	case "iowait":
		sapMetric = "WaitTimePercentage"
	case "steal":
		if cpu != "all" {
			return nil, errs.WrapConst(fmt.Errorf("steal time is only available for aggregate (all)"), zbxerr.ErrorInvalidParams)
		}
		sapMetric = "StealTimePercentage"
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported utilization type: %s", utilType), zbxerr.ErrorInvalidParams)
	}

	var key string
	if cpu == "all" {
		key = "OS." + sapMetric
	} else {
		key = fmt.Sprintf("Proc.%s[%s]", sapMetric, cpu)
	}

	val, ok := res[key]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", key), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", key, val.Value), zbxerr.ErrorCannotFetchData)
	}

	return floatVal, nil
}

// cpuLoadHandler returns CPU load average.
func cpuLoadHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch mode {
	case "avg1":
		sapMetric = "OS.LoadAverage1Min"
	case "avg5":
		sapMetric = "OS.LoadAverage5Min"
	case "avg15":
		sapMetric = "OS.LoadAverage15Min"
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported load average mode: %s", mode), zbxerr.ErrorInvalidParams)
	}

	val, ok := res[sapMetric]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", sapMetric), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", sapMetric, val.Value), zbxerr.ErrorCannotFetchData)
	}

	return floatVal, nil
}

// cpuNumHandler returns total number of logical CPUs.
func cpuNumHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	val, ok := res["OS.CPUCount"]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric OS.CPUCount not found"), zbxerr.ErrorCannotFetchData)
	}

	uintVal, err := strconv.ParseUint(val.Value, 10, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric OS.CPUCount value: %s", val.Value), zbxerr.ErrorCannotFetchData)
	}

	return uintVal, nil
}

// cpuIntrHandler returns interrupt rate per second.
func cpuIntrHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	val, ok := res["OS.InterruptRate"]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric OS.InterruptRate not found"), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric OS.InterruptRate value: %s", val.Value), zbxerr.ErrorCannotFetchData)
	}

	return floatVal, nil
}

// cpuSwitchesHandler returns context switch rate per second.
func cpuSwitchesHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	val, ok := res["OS.ContextSwitchRate"]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric OS.ContextSwitchRate not found"), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric OS.ContextSwitchRate value: %s", val.Value), zbxerr.ErrorCannotFetchData)
	}

	return floatVal, nil
}
