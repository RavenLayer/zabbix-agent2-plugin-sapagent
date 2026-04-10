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

const versionXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Body>
    <ExecuteOperationResponse>
      <result>
        <mOperationResults>
          <item>
            <mMessageValue>kernel release 753</mMessageValue>
          </item>
          <item>
            <mMessageValue>patch number 1100</mMessageValue>
          </item>
          <item>
            <mMessageValue>latest change number 2123456</mMessageValue>
          </item>
        </mOperationResults>
      </result>
    </ExecuteOperationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_versionHandler(t *testing.T) {
	client := newMockClient(versionXmlResponse)

	got, err := versionHandler(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("versionHandler() error = %v", err)
	}

	want := "753 PL 1100 (2123456)"
	if got != want {
		t.Errorf("versionHandler() = %v, want %v", got, want)
	}
}
