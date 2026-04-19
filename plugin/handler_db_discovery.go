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
	"sort"
	"strings"

	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/zbxerr"
)

// databaseDiscoveryHandler returns JSON discovery data of databases.
func databaseDiscoveryHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.ListDatabaseSystems()
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	type dbDiscoveryItem struct {
		Name     string `json:"{#DB.NAME}"`
		Protocol string `json:"{#DB.PROTOCOL}"`
		Host     string `json:"{#DB.HOST}"`
		Port     string `json:"{#DB.PORT}"`
	}

	var discovery []dbDiscoveryItem
	for _, db := range res {
		item := dbDiscoveryItem{}
		for _, prop := range db.Database {
			if prop.Key == "Database/Name" {
				// Parse "<database>@<instance>" -> "<database>"
				item.Name = strings.Split(prop.Value, "@")[0]
			}
		}

		for _, prop := range db.Properties {
			if prop.Key == "ConnectAddress" {
				// Parse "Protocol=<proto>;Host=<hostname>;Port=<port>"
				parts := strings.Split(prop.Value, ";")
				for _, part := range parts {
					kv := strings.SplitN(part, "=", 2)
					if len(kv) != 2 {
						continue
					}
					switch strings.ToLower(kv[0]) {
					case "protocol":
						item.Protocol = kv[1]
					case "host":
						item.Host = kv[1]
					case "port":
						item.Port = kv[1]
					}
				}
			}
		}

		if item.Name != "" {
			discovery = append(discovery, item)
		}
	}

	sort.Slice(discovery, func(i, j int) bool {
		return discovery[i].Name < discovery[j].Name
	})

	jsonRes, err := json.Marshal(discovery)
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	return string(jsonRes), nil
}
