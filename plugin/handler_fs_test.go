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
	"context"
	"testing"
)

const fsXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:GetCIMObjectResponse>
      <result>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>/.FS.AvailableSpace.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>FS.AvailableSpace</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>/</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.937680+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>33661591552</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>/.FS.Size.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>FS.Size</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>/</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.937680+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>46381764608</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>/boot.FS.AvailableSpace.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>FS.AvailableSpace</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>/boot</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.933324+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>589307904</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>/boot.FS.Size.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>FS.Size</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>/boot</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.933324+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>1020702720</mValue></item>
          </mProperties>
        </item>
      </result>
    </SAPHostControl:GetCIMObjectResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_filesystemDiscoveryHandler(t *testing.T) {
	client := newMockClient(fsXmlResponse)

	got, err := filesystemDiscoveryHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("filesystemDiscoveryHandler() error = %v", err)
	}

	want := `[{"{#FSNAME}":"/"},{"{#FSNAME}":"/boot"}]`
	if got != want {
		t.Errorf("filesystemDiscoveryHandler() = %v, want %v", got, want)
	}
}

func Test_filesystemSizeHandler(t *testing.T) {
	client := newMockClient(fsXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"root total", map[string]string{"FS": "/", "Mode": "total"}, uint64(46381764608), false},
		{"root free", map[string]string{"FS": "/", "Mode": "free"}, uint64(33661591552), false},
		{"boot total", map[string]string{"FS": "/boot", "Mode": "total"}, uint64(1020702720), false},
		{"boot free", map[string]string{"FS": "/boot", "Mode": "free"}, uint64(589307904), false},
		{"invalid mode", map[string]string{"FS": "/", "Mode": "invalid"}, nil, true},
		{"not found", map[string]string{"FS": "/nonexistent", "Mode": "total"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filesystemSizeHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("filesystemSizeHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("filesystemSizeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
