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

const pingSuccessXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Body>
    <ExecuteOperationResponse>
      <result>
        <mOperationResults>
          <item>
            <mMessageValue>SUCCESS</mMessageValue>
          </item>
        </mOperationResults>
      </result>
    </ExecuteOperationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

const pingFailureXmlResponse = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Body>
    <ExecuteOperationResponse>
      <result>
        <mOperationResults>
          <item>
            <mMessageValue>FAILURE</mMessageValue>
          </item>
        </mOperationResults>
      </result>
    </ExecuteOperationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func Test_pingHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(pingSuccessXmlResponse)
		got, err := pingHandler(context.Background(), client, nil)
		if err != nil {
			t.Errorf("pingHandler() error = %v", err)
			return
		}
		if got != pingOk {
			t.Errorf("pingHandler() = %v, want %v", got, pingOk)
		}
	})

	t.Run("failure", func(t *testing.T) {
		client := newMockClient(pingFailureXmlResponse)
		got, err := pingHandler(context.Background(), client, nil)
		if err != nil {
			t.Errorf("pingHandler() error = %v", err)
			return
		}
		if got != pingError {
			t.Errorf("pingHandler() = %v, want %v", got, pingError)
		}
	})
}
