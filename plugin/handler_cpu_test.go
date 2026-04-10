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

const cpuXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:GetCIMObjectResponse>
      <result>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.CPUCount.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.CPUCount</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.926878+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>4</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.LoadAverage1Min.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.LoadAverage1Min</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.926920+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>0.15</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.UserTimePercentage.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.UserTimePercentage</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.926596+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>12.5</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Proc.UserTimePercentage</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>0.Proc.UserTimePercentage.60</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>0</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.926668+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>10.0</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>Proc.UserTimePercentage</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>1.Proc.UserTimePercentage.60</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>1</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.926713+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>15.0</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.InterruptRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.InterruptRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.926852+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>1500</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.ContextSwitchRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.ContextSwitchRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.926864+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>2500</mValue></item>
          </mProperties>
        </item>
      </result>
    </SAPHostControl:GetCIMObjectResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_cpuDiscoveryHandler(t *testing.T) {
	client := newMockClient(cpuXmlResponse)

	got, err := cpuDiscoveryHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("cpuDiscoveryHandler() error = %v", err)
	}

	want := `[{"{#CPU.ID}":"0"},{"{#CPU.ID}":"1"}]`
	if got != want {
		t.Errorf("cpuDiscoveryHandler() = %v, want %v", got, want)
	}
}

func Test_cpuUtilHandler(t *testing.T) {
	client := newMockClient(cpuXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"all user", map[string]string{"CPU": "all", "Type": "user"}, 12.5, false},
		{"cpu0 user", map[string]string{"CPU": "0", "Type": "user"}, 10.0, false},
		{"cpu1 user", map[string]string{"CPU": "1", "Type": "user"}, 15.0, false},
		{"invalid type", map[string]string{"CPU": "all", "Type": "invalid"}, nil, true},
		{"not found", map[string]string{"CPU": "2", "Type": "user"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cpuUtilHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("cpuUtilHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cpuUtilHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cpuLoadHandler(t *testing.T) {
	client := newMockClient(cpuXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"avg1", map[string]string{"Mode": "avg1"}, 0.15, false},
		{"invalid mode", map[string]string{"Mode": "invalid"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cpuLoadHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("cpuLoadHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cpuLoadHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cpuNumHandler(t *testing.T) {
	client := newMockClient(cpuXmlResponse)

	got, err := cpuNumHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("cpuNumHandler() error = %v", err)
	}

	if got != uint64(4) {
		t.Errorf("cpuNumHandler() = %v, want %v", got, uint64(4))
	}
}

func Test_cpuIntrHandler(t *testing.T) {
	client := newMockClient(cpuXmlResponse)

	got, err := cpuIntrHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("cpuIntrHandler() error = %v", err)
	}

	if got != 1500.0 {
		t.Errorf("cpuIntrHandler() = %v, want %v", got, 1500.0)
	}
}

func Test_cpuSwitchesHandler(t *testing.T) {
	client := newMockClient(cpuXmlResponse)

	got, err := cpuSwitchesHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("cpuSwitchesHandler() error = %v", err)
	}

	if got != 2500.0 {
		t.Errorf("cpuSwitchesHandler() = %v, want %v", got, 2500.0)
	}
}
