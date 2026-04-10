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

	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/zbxerr"
)

// osNameHandler returns host name information.
func osNameHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetComputerSystem()
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	val, ok := res["SAP_ITSAMOperatingSystem.CSName"]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("property SAP_ITSAMOperatingSystem.CSName not found"), zbxerr.ErrorCannotFetchData)
	}

	return val.Value, nil
}

// osVersionHandler returns host version information.
func osVersionHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	res, err := client.GetComputerSystem()
	if err != nil {
		return nil, errs.WrapConst(err, zbxerr.ErrorCannotFetchData)
	}

	val, ok := res["SAP_ITSAMOperatingSystem.Caption"]
	if !ok {
		return nil, errs.WrapConst(fmt.Errorf("property SAP_ITSAMOperatingSystem.Caption not found"), zbxerr.ErrorCannotFetchData)
	}

	return val.Value, nil
}
