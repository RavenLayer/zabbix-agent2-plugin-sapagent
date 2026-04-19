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
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ArgumentItem struct {
	Key   string `xml:"mKey"`
	Value string `xml:"mValue"`
}

type GetCIMObjectRequest struct {
	XMLName   xml.Name       `xml:"urn:SAPHostControl GetCIMObject"`
	Arguments []ArgumentItem `xml:"aArguments>item"`
}

type PropertyItem struct {
	Name  string `xml:"mName"`
	Type  string `xml:"mType"`
	Value string `xml:"mValue"`
}

type CIMObject struct {
	Properties []PropertyItem `xml:"mProperties>item"`
}

func (c CIMObject) GetProperty(name string) string {
	for _, p := range c.Properties {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}

type GetCIMObjectResponse struct {
	XMLName xml.Name    `xml:"urn:SAPHostControl GetCIMObjectResponse"`
	Objects []CIMObject `xml:"result>item"`
}

type GetComputerSystemRequest struct {
	XMLName   xml.Name       `xml:"urn:SAPHostControl GetComputerSystem"`
	Arguments []ArgumentItem `xml:"aArguments>item"`
}

type GetComputerSystemResponse struct {
	XMLName xml.Name `xml:"urn:SAPHostControl GetComputerSystemResponse"`
	Result  struct {
		Properties []PropertyItem `xml:"mProperties>item"`
		Members    []CIMObject    `xml:"mMembers>item"`
	} `xml:"result"`
}

type ListDatabaseSystemsRequest struct {
	XMLName xml.Name `xml:"urn:SAPHostControl ListDatabaseSystems"`
}

type DatabaseItem struct {
	Status     string         `xml:"mStatus"`
	Database   []ArgumentItem `xml:"mDatabase>item"`
	Properties []ArgumentItem `xml:"mProperties>item"`
}

type ListDatabaseSystemsResponse struct {
	XMLName xml.Name       `xml:"ListDatabaseSystemsResponse"`
	Result  []DatabaseItem `xml:"result>item"`
}

type ExecuteOperationRequest struct {
	XMLName   xml.Name       `xml:"urn:SAPHostControl ExecuteOperation"`
	Operation string         `xml:"aOperation"`
	Arguments []ArgumentItem `xml:"aArguments>item,omitempty"`
	Options   ExecuteOptions `xml:"aOptions"`
}

type ExecuteOptions struct {
	Timeout     int    `xml:"mTimeout"`
	Options     string `xml:"mOptions"`
	SoftTimeout int    `xml:"mSoftTimeout"`
	Properties  string `xml:"mProperties"`
}

type ExecuteOperationResponse struct {
	XMLName xml.Name              `xml:"ExecuteOperationResponse"`
	Results []OperationResultItem `xml:"result>mOperationResults>item"`
}

type OperationResultItem struct {
	Key   string `xml:"mMessageKey"`
	Value string `xml:"mMessageValue"`
}

// SAPClient extends SOAPClient with SAP-specific logic and cache access
type SAPClient struct {
	*SOAPClient
	cache *SessionCache
}

// SAPService acts as the "Middle Man" between Zabbix and SOAP
type SAPService struct {
	mu      sync.Mutex
	cache   *SessionCache
	clients map[string]*SAPClient
}

var (
	reKernel = regexp.MustCompile(`kernel release\s+([\d]+)`)
	rePatch  = regexp.MustCompile(`patch number\s+([\d]+)`)
	reChange = regexp.MustCompile(`latest change number\s+([\d]+)`)
)

func NewSAPService(ttl time.Duration) *SAPService {
	return &SAPService{
		cache:   NewSessionCache(ttl),
		clients: make(map[string]*SAPClient),
	}
}

func (s *SAPService) GetClient(params map[string]string, timeout int) (*SAPClient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	uri := params[uriParam]
	user := params[usernameParam]
	pass := params[passwordParam]
	tlsCaFile := params[tlsRootCAFileParam]
	tlsServerName := params[tlsServerNameParam]
	tlsInsecureSkipVerify, _ := strconv.ParseBool(params[tlsInsecureSkipVerifyParam])

	// Create a unique key for the connection pool
	clientKey := fmt.Sprintf("%s|%s|%s|%s|%t", uri, user, tlsCaFile, tlsServerName, tlsInsecureSkipVerify)

	if client, ok := s.clients[clientKey]; ok {
		return client, nil
	}

	soapClient := NewClient(uri, time.Duration(timeout)*time.Second, tlsInsecureSkipVerify)
	soapClient.SetBasicAuth(user, pass)
	_ = soapClient.SetTLSConfig(tlsCaFile, tlsServerName, tlsInsecureSkipVerify)

	client := &SAPClient{
		SOAPClient: soapClient,
		cache:      s.cache,
	}

	s.clients[clientKey] = client
	return client, nil
}

// GetCIMObject handles the GetCIMObjectRequest SOAP call
func (c *SAPClient) GetCIMObject(reqType string) (map[string]MetricResult, error) {
	sessionKey := fmt.Sprintf("%s|%s", c.URL, c.Username)

	// 1. Check Cache
	if data, fresh := c.cache.GetCachedData(sessionKey, reqType); fresh {
		return data, nil
	}

	// 2. SOAP Call
	req := GetCIMObjectRequest{
		Arguments: []ArgumentItem{{Key: "EnumerateInstances", Value: reqType}},
	}

	res := GetCIMObjectResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
	defer cancel()

	err := c.Call(ctx, "GetCIMObject", req, &res)
	if err != nil {
		return nil, err
	}

	// 3. Transform the raw SOAP response into our MetricResult map
	newData := make(map[string]MetricResult)
	for _, obj := range res.Objects {
		id := obj.GetProperty("MetricDefinitionId")
		if id == "" {
			id = obj.GetProperty("Id")
		}
		val := obj.GetProperty("MetricValue")
		if val == "" {
			val = obj.GetProperty("Value")
		}

		if id == "" {
			continue
		}

		inst := obj.GetProperty("InstanceID")
		cleanInst := sanitizeInstanceID(inst, id)
		measuredElement := obj.GetProperty("MeasuredElementName")

		// Create a unique key for the map: "MetricName" or "MetricName[Instance]"
		var key string
		if measuredElement == "OperatingSystem" {
			key = id
		} else if measuredElement != "" {
			key = fmt.Sprintf("%s[%s]", id, measuredElement)
		} else if inst != "" && cleanInst != id {
			key = fmt.Sprintf("%s[%s]", id, cleanInst)
		}

		newData[key] = MetricResult{
			Value:     val,
			Timestamp: obj.GetProperty("TimeStamp"),
		}
	}

	// 4. Update Cache
	c.cache.UpdateSession(sessionKey, reqType, newData)
	return newData, nil
}

// GetComputerSystem handles the GetComputerSystemRequest SOAP call
func (c *SAPClient) GetComputerSystem() (map[string]MetricResult, error) {
	sessionKey := fmt.Sprintf("%s|%s", c.URL, c.Username)
	reqType := "GetComputerSystem"

	// 1. Check Cache
	if data, fresh := c.cache.GetCachedData(sessionKey, reqType); fresh {
		return data, nil
	}

	// 2. SOAP Call
	req := GetComputerSystemRequest{
		Arguments: []ArgumentItem{{Key: "provider", Value: "sapitsam"}},
	}

	res := GetComputerSystemResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
	defer cancel()

	err := c.Call(ctx, "GetComputerSystem", req, &res)
	if err != nil {
		return nil, err
	}

	// 3. Transform the raw SOAP response into our MetricResult map
	newData := make(map[string]MetricResult)
	// Top-level properties
	className := ""
	for _, p := range res.Result.Properties {
		if p.Name == "CreationClassName" {
			className = p.Value
		}
	}
	for _, p := range res.Result.Properties {
		key := p.Name
		if className != "" {
			key = fmt.Sprintf("%s.%s", className, p.Name)
		}
		newData[key] = MetricResult{
			Value: p.Value,
		}
	}
	// Members
	for _, member := range res.Result.Members {
		mClassName := member.GetProperty("CreationClassName")
		if mClassName == "" {
			continue
		}
		for _, p := range member.Properties {
			key := fmt.Sprintf("%s.%s", mClassName, p.Name)
			newData[key] = MetricResult{
				Value: p.Value,
			}
		}
	}

	// 4. Update Cache
	c.cache.UpdateSession(sessionKey, reqType, newData)
	return newData, nil
}

// Utility to normalize metric instance identifiers (e.g. "ens3.Net.KByteReceiveRate.60" -> "ens3")
func sanitizeInstanceID(instID, defID string) string {
	if instID == "" {
		return "default"
	}

	// 1. Remove aggregation interval (e.g., .60)
	parts := strings.Split(instID, ".")
	if len(parts) > 1 {
		last := parts[len(parts)-1]
		isNum := true
		for _, r := range last {
			if r < '0' || r > '9' {
				isNum = false
				break
			}
		}
		if isNum {
			parts = parts[:len(parts)-1]
		}
	}

	res := strings.Join(parts, ".")

	// 2. Remove MetricDefinitionId if it's part of the instance name
	res = strings.TrimSuffix(res, "."+defID)

	// 3. Clean up leading dots/slashes
	res = strings.TrimPrefix(res, ".")
	res = strings.TrimPrefix(res, "/")

	if res == "" {
		return "root"
	}

	return res
}

// ExecutePing handles the ExecuteOperationRequest SOAP call without caching
func (c *SAPClient) ExecutePing() error {
	req := ExecuteOperationRequest{
		Operation: "SAPHostAgent.Ping",
		Options: ExecuteOptions{
			Timeout:     -1,
			SoftTimeout: 0,
		},
	}

	res := ExecuteOperationResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
	defer cancel()

	err := c.Call(ctx, "ExecuteOperation", req, &res)
	if err != nil {
		return err
	}

	for _, item := range res.Results {
		if strings.Contains(item.Value, "SUCCESS") {
			return nil
		}
	}

	return fmt.Errorf("ping failed: SUCCESS not found in results")
}

// ExecuteVersionInfo handles the ExecuteOperationRequest SOAP call
func (c *SAPClient) ExecuteVersionInfo() (string, error) {
	sessionKey := fmt.Sprintf("%s|%s", c.URL, c.Username)
	reqType := "AgentVersion"

	// 1. Check Cache
	if data, fresh := c.cache.GetCachedData(sessionKey, reqType); fresh {
		return data["version"].Value, nil
	}

	// 2. SOAP Call
	req := ExecuteOperationRequest{
		Operation: "versioninfo",
		Options: ExecuteOptions{
			Timeout:     -1,
			SoftTimeout: 0,
		},
	}

	res := ExecuteOperationResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
	defer cancel()

	err := c.Call(ctx, "ExecuteOperation", req, &res)
	if err != nil {
		return "", err
	}

	// 3. Parse the results from the corrected hierarchy
	if len(res.Results) == 0 {
		return "", fmt.Errorf("XML decoder found 0 items via path result>mOperationResults>item")
	}

	var kernel, patch, change string
	for _, item := range res.Results {
		text := strings.TrimSpace(item.Value)
		if text == "" {
			continue
		}

		if kernel == "" {
			if match := reKernel.FindStringSubmatch(text); len(match) > 1 {
				kernel = match[1]
			}
		}
		if patch == "" {
			if match := rePatch.FindStringSubmatch(text); len(match) > 1 {
				patch = match[1]
			}
		}
		if change == "" {
			if match := reChange.FindStringSubmatch(text); len(match) > 1 {
				change = match[1]
			}
		}
	}

	if kernel == "" || patch == "" || change == "" {
		return "", fmt.Errorf("regex failed to find version in %d items", len(res.Results))
	}

	fullVersion := fmt.Sprintf("%s PL %s (%s)", kernel, patch, change)

	// 4. Update Cache
	c.cache.UpdateSession(sessionKey, reqType, map[string]MetricResult{
		"version": {Value: fullVersion},
	})

	return fullVersion, nil
}

// ListDatabaseSystems handles the ListDatabaseSystemsRequest SOAP call without caching
func (c *SAPClient) ListDatabaseSystems() ([]DatabaseItem, error) {
	req := ListDatabaseSystemsRequest{}
	res := ListDatabaseSystemsResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
	defer cancel()

	err := c.Call(ctx, "ListDatabaseSystems", req, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}
