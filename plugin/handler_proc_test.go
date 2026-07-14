//go:build tests
// +build tests

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
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const listInstancesXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:ListInstancesResponse>
      <aInstances>
        <item>
          <mSid>HBD</mSid>
          <mHostname>hanapm</mHostname>
          <mSystemNumber>00</mSystemNumber>
          <mSapVersionInfo>753, patch 1314</mSapVersionInfo>
        </item>
        <item>
          <mSid>S4H</mSid>
          <mHostname>sapapp</mHostname>
          <mSystemNumber>01</mSystemNumber>
          <mSapVersionInfo>773, patch 100</mSapVersionInfo>
        </item>
      </aInstances>
    </SAPHostControl:ListInstancesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

const hdbProcessListXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPControl="urn:SAPControl">
  <SOAP-ENV:Body>
    <SAPControl:GetProcessListResponse>
      <process>
        <item>
          <name>hdbdaemon</name>
          <description>HDB Daemon</description>
          <dispstatus>SAPControl-GREEN</dispstatus>
          <textstatus>Running</textstatus>
          <starttime>2026 07 09 09:08:22</starttime>
          <elapsedtime>1:21:58</elapsedtime>
          <pid>8850</pid>
        </item>
      </process>
    </SAPControl:GetProcessListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

const s4hProcessListXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPControl="urn:SAPControl">
  <SOAP-ENV:Body>
    <SAPControl:GetProcessListResponse>
      <process>
        <item>
          <name>disp+work</name>
          <description>Dispatcher</description>
          <dispstatus>SAPControl-GREEN</dispstatus>
          <textstatus>Running</textstatus>
          <starttime>2026 07 09 09:10:00</starttime>
          <elapsedtime>1:20:00</elapsedtime>
          <pid>9900</pid>
        </item>
      </process>
    </SAPControl:GetProcessListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_procDiscoveryHandler(t *testing.T) {
	// Mock primary client
	client := newMockClient(listInstancesXmlResponse)

	// Save original newClientFn and restore after test
	origNewClientFn := newClientFn
	defer func() { newClientFn = origNewClientFn }()

	// Mock newClientFn to return appropriate response based on instance port
	newClientFn = func(uri string, timeout time.Duration, skipVerify bool) *SOAPClient {
		var xmlResponse string
		if containsPortOrSocket(uri, "00") {
			xmlResponse = hdbProcessListXmlResponse
		} else if containsPortOrSocket(uri, "01") {
			xmlResponse = s4hProcessListXmlResponse
		} else {
			xmlResponse = `<?xml version="1.0" encoding="UTF-8"?><SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"><SOAP-ENV:Body/></SOAP-ENV:Envelope>`
		}

		transport := &mockTransport{
			response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(xmlResponse)),
				Header:     make(http.Header),
			},
		}

		return &SOAPClient{
			HTTPClient: &http.Client{
				Transport: transport,
				Timeout:   timeout,
			},
			URL: uri,
		}
	}

	params := map[string]string{
		uriParam: "http://localhost:1128",
	}

	got, err := procDiscoveryHandler(context.Background(), client, params)
	if err != nil {
		t.Fatalf("procDiscoveryHandler() error = %v", err)
	}

	want := `[{"{#INSTANCE.NAME}":"HBD","{#PROCESS.NAME}":"hdbdaemon","{#PROCESS.DESC}":"HDB Daemon"},{"{#INSTANCE.NAME}":"S4H","{#PROCESS.NAME}":"disp+work","{#PROCESS.DESC}":"Dispatcher"}]`
	if got != want {
		t.Errorf("procDiscoveryHandler() = %v, want %v", got, want)
	}
}

func Test_procGetHandler(t *testing.T) {
	client := newMockClient(listInstancesXmlResponse)

	origNewClientFn := newClientFn
	defer func() { newClientFn = origNewClientFn }()

	newClientFn = func(uri string, timeout time.Duration, skipVerify bool) *SOAPClient {
		var xmlResponse string
		if containsPortOrSocket(uri, "00") {
			xmlResponse = hdbProcessListXmlResponse
		} else if containsPortOrSocket(uri, "01") {
			xmlResponse = s4hProcessListXmlResponse
		} else {
			xmlResponse = `<?xml version="1.0" encoding="UTF-8"?><SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"><SOAP-ENV:Body/></SOAP-ENV:Envelope>`
		}

		transport := &mockTransport{
			response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(xmlResponse)),
				Header:     make(http.Header),
			},
		}

		return &SOAPClient{
			HTTPClient: &http.Client{
				Transport: transport,
				Timeout:   timeout,
			},
			URL: uri,
		}
	}

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{
			name: "valid hdbdaemon",
			params: map[string]string{
				uriParam:   "http://localhost:1128",
				"Instance": "HBD",
				"Process":  "hdbdaemon",
			},
			want:    `{"name":"hdbdaemon","description":"HDB Daemon","dispstatus":"SAPControl-GREEN","textstatus":"Running","starttime":"2026 07 09 09:08:22","elapsedtime":"1:21:58","pid":8850}`,
			wantErr: false,
		},
		{
			name: "valid dispatcher",
			params: map[string]string{
				uriParam:   "http://localhost:1128",
				"Instance": "S4H",
				"Process":  "disp+work",
			},
			want:    `{"name":"disp+work","description":"Dispatcher","dispstatus":"SAPControl-GREEN","textstatus":"Running","starttime":"2026 07 09 09:10:00","elapsedtime":"1:20:00","pid":9900}`,
			wantErr: false,
		},
		{
			name: "invalid instance",
			params: map[string]string{
				uriParam:   "http://localhost:1128",
				"Instance": "XYZ",
				"Process":  "hdbdaemon",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid process",
			params: map[string]string{
				uriParam:   "http://localhost:1128",
				"Instance": "HBD",
				"Process":  "nonexistent",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing parameters",
			params: map[string]string{
				uriParam: "http://localhost:1128",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := procGetHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("procGetHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("procGetHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BuildInstanceURL(t *testing.T) {
	tests := []struct {
		name         string
		origURL      string
		systemNumber string
		want         string
		wantErr      bool
	}{
		{
			name:         "HTTP url",
			origURL:      "http://localhost:1128",
			systemNumber: "00",
			want:         "http://localhost:50013",
			wantErr:      false,
		},
		{
			name:         "HTTPS url",
			origURL:      "https://sap-host:1129",
			systemNumber: "01",
			want:         "https://sap-host:50114",
			wantErr:      false,
		},
		{
			name:         "Unix socket path",
			origURL:      "unix:///tmp/.sapstream51128",
			systemNumber: "99",
			want:         "unix:///tmp/.sapstream59913",
			wantErr:      false,
		},
		{
			name:         "Unix socket path in subfolder",
			origURL:      "unix:///var/run/sap/.sapstream51128",
			systemNumber: "02",
			want:         "unix:///var/run/sap/.sapstream50213",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildInstanceURL(tt.origURL, tt.systemNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildInstanceURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildInstanceURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func containsPortOrSocket(uri string, systemNumber string) bool {
	// Checks if uri contains instance port (e.g. 50013 for 00) or socket name (e.g. .sapstream50013)
	portStr := "5" + systemNumber + "1" // e.g. "5001" or "5011"
	return strings.Contains(uri, portStr)
}
