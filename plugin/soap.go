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
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// SOAP structures
type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *Header  `xml:",omitempty"`
	Body    Body
}

type Header struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Content []interface{}
}

type Body struct {
	XMLName xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Payload interface{} `xml:",any"`
	Fault   *SOAPFault  `xml:",omitempty"`
}

// SOAPFault implements the error interface
type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`
	Code    string   `xml:"faultcode"`
	String  string   `xml:"faultstring"`
	Actor   string   `xml:"faultactor"`
	Detail  string   `xml:"detail"`
}

func (f *SOAPFault) Error() string {
	return fmt.Sprintf("SOAP Fault: [%s] %s", f.Code, f.String)
}

type SOAPClient struct {
	HTTPClient *http.Client
	URL        string
	Username   string
	Password   string
	Headers    map[string]string
}

// NewClient is a factory that decides between HTTP and Unix
func NewClient(uri string, timeout time.Duration, skipVerify bool) *SOAPClient {
	var transport *http.Transport

	if strings.HasPrefix(uri, "unix://") {
		// Example: unix:///tmp/.sapstream51128
		socketPath := strings.TrimPrefix(uri, "unix://")
		transport = &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
			},
		}
		// When using Unix sockets, the "host" in the URL is usually ignored
		// by the dialer but required by the http.Request.
		uri = "http://localhost"
	} else {
		// Standard HTTP/HTTPS
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
		}
	}

	return &SOAPClient{
		HTTPClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		URL:     uri,
		Headers: make(map[string]string),
	}
}

func (s *SOAPClient) SetBasicAuth(username, password string) {
	s.Username = username
	s.Password = password
}

func (s *SOAPClient) SetHeader(key, value string) {
	if s.Headers == nil {
		s.Headers = make(map[string]string)
	}
	s.Headers[key] = value
}

func (s *SOAPClient) SetTLSConfig(caFile, serverName string, skipVerify bool) error {
	if s.HTTPClient == nil || s.HTTPClient.Transport == nil {
		return fmt.Errorf("HTTPClient or Transport is nil")
	}

	transport, ok := s.HTTPClient.Transport.(*http.Transport)
	if !ok {
		return fmt.Errorf("transport is not *http.Transport")
	}

	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}

	transport.TLSClientConfig.InsecureSkipVerify = skipVerify
	transport.TLSClientConfig.ServerName = serverName

	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		transport.TLSClientConfig.RootCAs = caCertPool
	}

	return nil
}

func (s *SOAPClient) Call(ctx context.Context, action string, request interface{}, response interface{}) error {
	envelope := Envelope{Body: Body{Payload: request}}
	xmlBuffer := new(bytes.Buffer)
	if err := xml.NewEncoder(xmlBuffer).Encode(envelope); err != nil {
		return err
	}

	Impl.Debugf("SOAP Request[%s]: %s", action, xmlBuffer.String())

	req, err := http.NewRequestWithContext(ctx, "POST", s.URL, xmlBuffer)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", action)
	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}
	if s.Username != "" || s.Password != "" {
		req.SetBasicAuth(s.Username, s.Password)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	Impl.Debugf("SOAP Response[%s]: %d %s", action, resp.StatusCode, string(respBytes))

	// Unmarshal to a temporary envelope to check for Faults and Data
	respEnvelope := &Envelope{
		Body: Body{Payload: response},
	}

	if err := xml.Unmarshal(respBytes, respEnvelope); err != nil {
		return err
	}

	// Check if the Body contains a Fault
	if respEnvelope.Body.Fault != nil {
		return respEnvelope.Body.Fault
	}

	return nil
}
