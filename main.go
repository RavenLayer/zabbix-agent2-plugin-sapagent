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

package main

import (
	"errors"
	"fmt"
	"os"

	"golang.zabbix.com/sdk/plugin/container"
	"golang.zabbix.com/sdk/plugin/flag"
	"golang.zabbix.com/sdk/zbxerr"
	"ravenlayer.com/plugin/sapagent/plugin"
)

const COPYRIGHT_MESSAGE = `Copyright (c) 2026, RavenLayer <admin@ravenlayer.com>. All rights reserved.`

func main() {
	err := flag.HandleFlags(
		plugin.Name,
		os.Args[0],
		COPYRIGHT_MESSAGE,
		plugin.PLUGIN_VERSION_RC,
		plugin.PLUGIN_VERSION_MAJOR,
		plugin.PLUGIN_VERSION_MINOR,
		plugin.PLUGIN_VERSION_PATCH,
	)
	if err != nil {
		if !errors.Is(err, zbxerr.ErrorOSExitZero) {
			panic(fmt.Sprintf("failed to handle flags %s", err.Error()))
		}

		return
	}

	h, err := container.NewHandler(plugin.Impl.Name())
	if err != nil {
		panic(fmt.Sprintf("failed to create plugin handler %s", err.Error()))
	}
	plugin.Impl.Logger = &h

	err = h.Execute()
	if err != nil {
		panic(fmt.Sprintf("failed to execute plugin handler %s", err.Error()))
	}
}
