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

const dbDiscoveryXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SAPHostControl="urn:SAPHostControl">
  <SOAP-ENV:Body>
    <SAPHostControl:ListDatabaseSystemsResponse>
      <result>
        <item>
          <mStatus>SAPHostControl-DB-RUNNING</mStatus>
          <mDatabase>
            <item><mKey>Database/Name</mKey><mValue>SYSTEMDB@HDB</mValue></item>
            <item><mKey>Database/Vendor</mKey><mValue>HDB</mValue></item>
            <item><mKey>Database/Type</mKey><mValue>hdb</mValue></item>
          </mDatabase>
          <mProperties>
             <item><mKey>ConnectAddress</mKey><mValue>Protocol=tcp;Host=sapel9.ravenlayer.loc;Port=30013</mValue></item>
          </mProperties>
        </item>
        <item>
          <mStatus>SAPHostControl-DB-WARNING</mStatus>
          <mDatabase>
            <item><mKey>Database/Name</mKey><mValue>HDB@HDB</mValue></item>
            <item><mKey>Database/Vendor</mKey><mValue>HDB</mValue></item>
            <item><mKey>Database/Type</mKey><mValue>hdb</mValue></item>
          </mDatabase>
          <mProperties>
             <item><mKey>ConnectAddress</mKey><mValue>Protocol=tcp;Host=sapel9.ravenlayer.loc;Port=30015</mValue></item>
          </mProperties>
        </item>
      </result>
    </SAPHostControl:ListDatabaseSystemsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_databaseDiscoveryHandler(t *testing.T) {
	client := newMockClient(dbDiscoveryXmlResponse)

	got, err := databaseDiscoveryHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("databaseDiscoveryHandler() error = %v", err)
	}

	want := `[{"{#DB.NAME}":"HDB","{#DB.PROTOCOL}":"tcp","{#DB.HOST}":"sapel9.ravenlayer.loc","{#DB.PORT}":"30015"},{"{#DB.NAME}":"SYSTEMDB","{#DB.PROTOCOL}":"tcp","{#DB.HOST}":"sapel9.ravenlayer.loc","{#DB.PORT}":"30013"}]`
	if got != want {
		t.Errorf("databaseDiscoveryHandler() = %v, want %v", got, want)
	}
}
