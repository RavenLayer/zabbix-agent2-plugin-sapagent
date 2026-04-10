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
	"reflect"
	"testing"

	"golang.zabbix.com/sdk/plugin"
)

func TestPlugin_Configure(t *testing.T) {
	type fields struct {
		Base    plugin.Base
		options PluginOptions
	}
	type args struct {
		global  *plugin.GlobalOptions
		options interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   PluginOptions
	}{
		{
			"named sessions",
			fields{},
			args{
				&plugin.GlobalOptions{Timeout: 30},
				[]byte("Sessions.Session1.Uri=http://1.2.3.4:1128\nSessions.Session1.Username=user1\nSessions.Session1.Password=pass1"),
			},
			PluginOptions{
				Timeout:  30,
				CacheTTL: 60,
				Sessions: map[string]Session{
					"Session1": {
						Uri:      "http://1.2.3.4:1128",
						Username: "user1",
						Password: "pass1",
					},
				},
			},
		},
		{
			"default options",
			fields{},
			args{
				&plugin.GlobalOptions{Timeout: 10},
				[]byte("Default.Uri=http://localhost:1128\nDefault.Username=sapadm\nDefault.Password=secret"),
			},
			PluginOptions{
				Timeout:  10,
				CacheTTL: 60,
				Default: Session{
					Uri:      "http://localhost:1128",
					Username: "sapadm",
					Password: "secret",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				Base:    tt.fields.Base,
				options: tt.fields.options,
			}
			p.Init(Name)
			p.Configure(tt.args.global, tt.args.options)

			if !reflect.DeepEqual(p.options, tt.want) {
				t.Errorf("Plugin.Configure() = %v, want %v", p.options, tt.want)
			}
		})
	}
}

func TestPlugin_EvalParamsWithSessions(t *testing.T) {
	p := &Plugin{}
	p.Init(Name)
	p.Configure(&plugin.GlobalOptions{Timeout: 30}, []byte("Sessions.S1.Uri=http://1.1.1.1:1128\nSessions.S1.Username=u1\nSessions.S1.Password=p1"))

	// Test named session evaluation
	// Use sap.agent.ping as a simple metric key
	key := "sap.agent.ping"
	rawParams := []string{"S1"}

	params, _, _, err := metrics[key].EvalParams(rawParams, p.options.Sessions)
	if err != nil {
		t.Fatalf("EvalParams failed: %v", err)
	}

	if params["Uri"] != "http://1.1.1.1:1128" {
		t.Errorf("Expected Uri http://1.1.1.1:1128, got %s", params["Uri"])
	}
	if params["Username"] != "u1" {
		t.Errorf("Expected Username u1, got %s", params["Username"])
	}
	if params["Password"] != "p1" {
		t.Errorf("Expected Password p1, got %s", params["Password"])
	}
}
