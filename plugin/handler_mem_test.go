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
	"testing"
	"time"
)

type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func newMockClient(xmlResponse string) *SAPClient {
	transport := &mockTransport{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(xmlResponse)),
			Header:     make(http.Header),
		},
	}
	soapClient := &SOAPClient{
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   5 * time.Second,
		},
		URL: "http://localhost:1128",
	}
	return &SAPClient{
		SOAPClient: soapClient,
		cache:      NewSessionCache(1 * time.Minute),
	}
}

const memXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:GetCIMObjectResponse>
      <result>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.AvailablePhysicalMemory.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.AvailablePhysicalMemory</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927018+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>12342740</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.FreePhysicalMemory.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.FreePhysicalMemory</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927018+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>11503984</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.TotalVisibleMemorySize.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.TotalVisibleMemorySize</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927018+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>16110496</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.TotalVirtualMemorySize.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.TotalVirtualMemorySize</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927018+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>16110496</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.FreeSpaceInPagingFiles.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.FreeSpaceInPagingFiles</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927018+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>1024</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.SizeStoredInPagingFiles.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.SizeStoredInPagingFiles</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927018+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>2048</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.KBytesPagedInRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.KBytesPagedInRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927117+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>1.23</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.KBytesPagedOutRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.KBytesPagedOutRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927117+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>4.56</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.PagesPagedInRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.PagesPagedInRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927117+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>7.89</mValue></item>
          </mProperties>
        </item>
        <item>
          <mProperties>
            <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_MetricValue</mValue></item>
            <item><mName>InstanceID</mName><mType>string</mType><mValue>OS.PagesPagedOutRate.60</mValue></item>
            <item><mName>MetricDefinitionId</mName><mType>string</mType><mValue>OS.PagesPagedOutRate</mValue></item>
            <item><mName>MeasuredElementName</mName><mType>string</mType><mValue>OperatingSystem</mValue></item>
            <item><mName>TimeStamp</mName><mType>datetime</mType><mValue>20260321000933.927117+060</mValue></item>
            <item><mName>MetricValue</mName><mType>string</mType><mValue>0.12</mValue></item>
          </mProperties>
        </item>
      </result>
    </SAPHostControl:GetCIMObjectResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_memoryHandler(t *testing.T) {
	client := newMockClient(memXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"phy available", map[string]string{"Type": "phy", "Mode": "available"}, uint64(12342740 * 1024), false},
		{"phy free", map[string]string{"Type": "phy", "Mode": "free"}, uint64(11503984 * 1024), false},
		{"phy size", map[string]string{"Type": "phy", "Mode": "size"}, uint64(16110496 * 1024), false},
		{"virt size", map[string]string{"Type": "virt", "Mode": "size"}, uint64(16110496 * 1024), false},
		{"invalid type", map[string]string{"Type": "invalid", "Mode": "size"}, nil, true},
		{"invalid mode", map[string]string{"Type": "phy", "Mode": "invalid"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := memoryHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("memoryHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_swapSizeHandler(t *testing.T) {
	client := newMockClient(memXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"swap free", map[string]string{"Mode": "free"}, uint64(1024 * 1024), false},
		{"swap total", map[string]string{"Mode": "total"}, uint64(2048 * 1024), false},
		{"swap used", map[string]string{"Mode": "used"}, uint64(1024 * 1024), false},
		{"invalid mode", map[string]string{"Mode": "invalid"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := swapSizeHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("swapSizeHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("swapSizeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_swapInHandler(t *testing.T) {
	client := newMockClient(memXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"swap in bytes", map[string]string{"Mode": "bytes"}, uint64(1259), false},
		{"swap in pages", map[string]string{"Mode": "pages"}, 7.89, false},
		{"invalid mode", map[string]string{"Mode": "invalid"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := swapInHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("swapInHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("swapInHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_swapOutHandler(t *testing.T) {
	client := newMockClient(memXmlResponse)

	tests := []struct {
		name    string
		params  map[string]string
		want    interface{}
		wantErr bool
	}{
		{"swap out bytes", map[string]string{"Mode": "bytes"}, uint64(4669), false},
		{"swap out pages", map[string]string{"Mode": "pages"}, 0.12, false},
		{"invalid mode", map[string]string{"Mode": "invalid"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := swapOutHandler(context.Background(), client, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("swapOutHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("swapOutHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
