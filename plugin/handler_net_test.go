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

const netXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:GetCIMObjectResponse>
      <result>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>ens3.Net.KByteReceiveRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Net.KByteReceiveRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>ens3</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927651+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>23277.99</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>ens3.Net.KByteTransmitRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Net.KByteTransmitRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>ens3</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927651+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>425.93</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>ens3.Net.PacketReceiveRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Net.PacketReceiveRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>ens3</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927651+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>848.78</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>ens3.Net.PacketTransmitRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Net.PacketTransmitRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>ens3</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927651+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>763.98</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>lo.Net.KByteReceiveRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Net.KByteReceiveRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>lo</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927651+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>97.87</mValue></item>
          </mProperties>
        </item>
      </result>
    </SAPHostControl:GetCIMObjectResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_networkDiscoveryHandler(t *testing.T) {
	client := newMockClient(netXmlResponse)

	got, err := networkDiscoveryHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("networkDiscoveryHandler() error = %v", err)
	}

	want := `[{"{#IFNAME}":"ens3"},{"{#IFNAME}":"lo"}]`
	if got != want {
		t.Errorf("networkDiscoveryHandler() = %v, want %v", got, want)
	}
}

func Test_networkInHandler(t *testing.T) {
	client := newMockClient(netXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"ens3 bytes", map[string]string{"Iface": "ens3", "Mode": "bytes"}, uint64(23836661), false},
		{"ens3 packets", map[string]string{"Iface": "ens3", "Mode": "packets"}, 848.78, false},
		{"lo bytes", map[string]string{"Iface": "lo", "Mode": "bytes"}, uint64(100218), false},
		{"invalid mode", map[string]string{"Iface": "ens3", "Mode": "invalid"}, nil, true},
		{"not found", map[string]string{"Iface": "eth1", "Mode": "bytes"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := networkInHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("networkInHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("networkInHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_networkOutHandler(t *testing.T) {
	client := newMockClient(netXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"ens3 bytes", map[string]string{"Iface": "ens3", "Mode": "bytes"}, uint64(436152), false},
		{"ens3 packets", map[string]string{"Iface": "ens3", "Mode": "packets"}, 763.98, false},
		{"invalid mode", map[string]string{"Iface": "ens3", "Mode": "invalid"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := networkOutHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("networkOutHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("networkOutHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
