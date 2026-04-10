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

const devXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:GetCIMObjectResponse>
      <result>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>redhat-root.Disk.AverageQueueLength.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Disk.AverageQueueLength</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>redhat-root</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927451+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>0.04</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>redhat-root.Disk.Util.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Disk.Util</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>redhat-root</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927451+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>2.15</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>vda.Disk.AverageQueueLength.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Disk.AverageQueueLength</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>vda</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927292+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>0.04</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>vda.Disk.Util.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Disk.Util</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>vda</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927292+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>2.13</mValue></item>
          </mProperties>
        </item>
      </result>
    </SAPHostControl:GetCIMObjectResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_deviceDiscoveryHandler(t *testing.T) {
	client := newMockClient(devXmlResponse)

	got, err := deviceDiscoveryHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("deviceDiscoveryHandler() error = %v", err)
	}

	want := `[{"{#DEVNAME}":"redhat-root"},{"{#DEVNAME}":"vda"}]`
	if got != want {
		t.Errorf("deviceDiscoveryHandler() = %v, want %v", got, want)
	}
}

func Test_deviceHandler(t *testing.T) {
	// Add tput to mock response for testing
	xmlWithTput := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:GetCIMObjectResponse>
      <result>
        <item>
          <mProperties>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Disk.AverageQueueLength</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>vda</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>0.04</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Disk.Util</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>vda</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>2.13</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Disk.TotalThroughput</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>vda</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>123.45</mValue></item>
          </mProperties>
        </item>
      </result>
    </SAPHostControl:GetCIMObjectResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
	clientWithTput := newMockClient(xmlWithTput)

	tests := []struct {
		name    string
		client  *SAPClient
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"vda queue", clientWithTput, map[string]string{"Dev": "vda", "Mode": "queue"}, 0.04, false},
		{"vda util", clientWithTput, map[string]string{"Dev": "vda", "Mode": "util"}, 2.13, false},
		{"vda tput", clientWithTput, map[string]string{"Dev": "vda", "Mode": "tput"}, uint64(126412), false},
		{"invalid mode", clientWithTput, map[string]string{"Dev": "vda", "Mode": "invalid"}, nil, true},
		{"not found", clientWithTput, map[string]string{"Dev": "nonexistent", "Mode": "util"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deviceHandler(context.Background(), tt.client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("deviceHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("deviceHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
