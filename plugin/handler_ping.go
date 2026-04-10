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
)

const (
	pingError = 0
	pingOk    = 1
)

// pingHandler checks if we can connect to the SAP Agent using the dedicated ping method.
func pingHandler(ctx context.Context,
	client *SAPClient,
	_ map[string]string, _ ...string) (interface{}, error) {

	err := client.ExecutePing()
	if err != nil {
		return pingError, nil
	}

	return pingOk, nil
}
