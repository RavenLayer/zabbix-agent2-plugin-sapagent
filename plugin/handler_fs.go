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
	metricFSSize           = "FS.Size"
	metricFSAvailebleSpace = "FS.AvailableSpace"
)

// filesystemDiscoveryHandler returns JSON discovery data of mounted filesystem.
func filesystemDiscoveryHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	fsMap := make(map[string]struct{})
	for key := range res {
		if strings.HasPrefix(key, "FS.") {
			// Key format is MetricDefinitionId[MeasuredElementName]
			start := strings.Index(key, "[")
			end := strings.Index(key, "]")
			if start != -1 && end != -1 && end > start {
				fsName := key[start+1 : end]
				fsMap[fsName] = struct{}{}
			}
		}
	}

	type fsDiscoveryItem struct {
		FSName string `json:"{#FSNAME}"`
	}
	var discovery []fsDiscoveryItem
	for name := range fsMap {
		discovery = append(discovery, fsDiscoveryItem{FSName: name})
	}

	sort.Slice(discovery, func(i, j int) bool {
		return discovery[i].FSName < discovery[j].FSName
	})

	jsonRes, err := json.Marshal(discovery)
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	return string(jsonRes), nil
}

// filesystemSizeHandler returns host filesystem information.
func filesystemSizeHandler(ctx context.Context,
	client *SAPClient,
	params map[string]string, _ ...string) (interface{}, error) {

	fs := params["FS"]
	mode := params["Mode"]

	res, err := client.GetCIMObject("SAP_MetricValue")
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	var sapMetric string
	switch mode {
	case "total":
		sapMetric = metricFSSize
	case "free":
		sapMetric = metricFSAvailebleSpace
	default:
		return nil, errs.WrapConst(fmt.Errorf("unsupported size mode: %s", mode), zbxerr.ErrorInvalidParams)
	}

	key := fmt.Sprintf("%s[%s]", sapMetric, fs)
	val, ok := res[key]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("metric %s not found", key), zbxerr.ErrorCannotFetchData)
	}

	uintVal, err := strconv.ParseUint(val.Value, 10, 64)
	if err != nil {
		return nil, errs.WrapConst(fmt.Errorf("failed to parse metric %s value: %s", key, val.Value), zbxerr.ErrorCannotFetchData)
	}

	return uintVal, nil
}
