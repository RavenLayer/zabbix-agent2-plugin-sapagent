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

// We don't redefine mockTransport and newMockClient if they are already in the package,
// but to be safe and allow running this test file alone with 'go test handler_os_test.go handler_os.go ...'
// we will assume they are available from handler_mem_test.go since they are in the same package.

const osXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:GetComputerSystemResponse>
      <result>
        <mProperties>
          <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_ITSAMComputerSystem</mValue></item>
          <item><mName>Name</mName><mType>string</mType><mValue>sapel9</mValue></item>
        </mProperties>
        <mMembers>
          <item>
            <mProperties>
              <item><mName>CreationClassName</mName><mType>string</mType><mValue>SAP_ITSAMOperatingSystem</mValue></item>
              <item><mName>CSName</mName><mType>string</mType><mValue>sapel9</mValue></item>
              <item><mName>Caption</mName><mType>string</mType><mValue>Red Hat Enterprise Linux 9.7 (Plow)</mValue></item>
            </mProperties>
          </item>
        </mMembers>
      </result>
    </SAPHostControl:GetComputerSystemResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_osHandlers(t *testing.T) {
	client := newMockClient(osXmlResponse)

	t.Run("osNameHandler", func(t *testing.T) {
		got, err := osNameHandler(context.Background(), client, nil)
		if err != nil {
			t.Errorf("osNameHandler() error = %v", err)
			return
		}
		if got != "sapel9" {
			t.Errorf("osNameHandler() = %v, want %v", got, "sapel9")
		}
	})

	t.Run("osVersionHandler", func(t *testing.T) {
		got, err := osVersionHandler(context.Background(), client, nil)
		if err != nil {
			t.Errorf("osVersionHandler() error = %v", err)
			return
		}
		want := "Red Hat Enterprise Linux 9.7 (Plow)"
		if got != want {
			t.Errorf("osVersionHandler() = %v, want %v", got, want)
		}
	})
}
