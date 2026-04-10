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
	"strconv"

	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/zbxerr"
)

// memoryHandler returns host memory information.
func memoryHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	memType := params["Type"]
	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch memType {
	case "phy":
		switch mode {
		case "available":
			sapMetric = "OS.AvailablePhysicalMemory"
		case "free":
			sapMetric = "OS.FreePhysicalMemory"
		case "size":
			sapMetric = "OS.TotalVisibleMemorySize"
		default:
			return nil, errs.WrapConst(fmt.Errorf("unsupported mode: %s", mode), zbxerr.ErrorInvalidParams)
		}
	case "virt":
		switch mode {
		case "size":
			sapMetric = "OS.TotalVirtualMemorySize"
		default:
			return nil, errs.WrapConst(fmt.Errorf("unsupported mode: %s", mode), zbxerr.ErrorInvalidParams)
		}
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported memory type: %s", memType), zbxerr.ErrorInvalidParams)
	}

	val, ok := res[sapMetric]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", sapMetric), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", sapMetric, val.Value), zbxerr.ErrorCannotFetchData)
	}

	return uint64(floatVal * 1024), nil
}

// swapSizeHandler returns swap size information.
func swapSizeHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch mode {
	case "free":
		sapMetric = "OS.FreeSpaceInPagingFiles"
	case "used":
		// Used swap = Total - Free
		totalVal, okTotal := res["OS.SizeStoredInPagingFiles"]
		freeVal, okFree := res["OS.FreeSpaceInPagingFiles"]
		if !okTotal || !okFree {
			return nil, errs.WrapConst(fmt.Errorf("required swap metrics not found"), zbxerr.ErrorCannotFetchData)
		}
		total, errTotal := strconv.ParseFloat(totalVal.Value, 64)
		free, errFree := strconv.ParseFloat(freeVal.Value, 64)
		if errTotal != nil || errFree != nil {
			return nil, errs.WrapConst(fmt.Errorf("failed to parse swap values"), zbxerr.ErrorCannotFetchData)
		}
		return uint64((total - free) * 1024), nil
	case "total":
		sapMetric = "OS.SizeStoredInPagingFiles"
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported mode: %s", mode), zbxerr.ErrorInvalidParams)
	}

	val, ok := res[sapMetric]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", sapMetric), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", sapMetric, val.Value), zbxerr.ErrorCannotFetchData)
	}

	return uint64(floatVal * 1024), nil
}

// swapInHandler returns swap-in rate information.
func swapInHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch mode {
	case "bytes":
		sapMetric = "OS.KBytesPagedInRate"
	case "pages":
		sapMetric = "OS.PagesPagedInRate"
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported mode: %s", mode), zbxerr.ErrorInvalidParams)
	}

	val, ok := res[sapMetric]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", sapMetric), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", sapMetric, val.Value), zbxerr.ErrorCannotFetchData)
	}

	if mode == "bytes" {
		return uint64(floatVal * 1024), nil
	}

	return floatVal, nil
}

// swapOutHandler returns swap-out rate information.
func swapOutHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch mode {
	case "bytes":
		sapMetric = "OS.KBytesPagedOutRate"
	case "pages":
		sapMetric = "OS.PagesPagedOutRate"
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported mode: %s", mode), zbxerr.ErrorInvalidParams)
	}

	val, ok := res[sapMetric]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", sapMetric), zbxerr.ErrorCannotFetchData)
	}

	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", sapMetric, val.Value), zbxerr.ErrorCannotFetchData)
	}

	if mode == "bytes" {
		return uint64(floatVal * 1024), nil
	}

	return floatVal, nil
}
