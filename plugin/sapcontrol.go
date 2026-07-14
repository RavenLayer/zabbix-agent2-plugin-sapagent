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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type GetProcessListItem struct {
	Name        string `xml:"name" json:"name"`
	Description string `xml:"description" json:"description"`
	DispStatus  string `xml:"dispstatus" json:"dispstatus"`
	TextStatus  string `xml:"textstatus" json:"textstatus"`
	StartTime   string `xml:"starttime" json:"starttime"`
	ElapsedTime string `xml:"elapsedtime" json:"elapsedtime"`
	Pid         int    `xml:"pid" json:"pid"`
}

type GetProcessListRequest struct {
	XMLName xml.Name `xml:"urn:SAPControl GetProcessList"`
}

type GetProcessListResponse struct {
	XMLName   xml.Name             `xml:"urn:SAPControl GetProcessListResponse"`
	Processes []GetProcessListItem `xml:"process>item"`
}

// GetProcessList handles the GetProcessList SOAP call
func (c *SAPClient) GetProcessList() (map[string]MetricResult, error) {
	sessionKey := fmt.Sprintf("%s|%s", c.URL, c.Username)
	reqType := "GetProcessList"

	// 1. Check Cache
	if data, fresh := c.cache.GetCachedData(sessionKey, reqType); fresh {
		return data, nil
	}

	// 2. SOAP Call
	req := GetProcessListRequest{}
	res := GetProcessListResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
	defer cancel()

	err := c.Call(ctx, "GetProcessList", req, &res)
	if err != nil {
		return nil, err
	}

	newData := make(map[string]MetricResult)
	for _, proc := range res.Processes {
		bz, err := json.Marshal(proc)
		if err != nil {
			continue
		}
		newData[proc.Name] = MetricResult{
			Value: string(bz),
		}
	}

	// 3. Update Cache
	c.cache.UpdateSession(sessionKey, reqType, newData)
	return newData, nil
}

// BuildInstanceURL constructs the instance-specific URL for SAPControl based on systemNumber
func BuildInstanceURL(origURL string, systemNumber string) (string, error) {
	if strings.HasPrefix(origURL, "unix://") {
		socketPath := strings.TrimPrefix(origURL, "unix://")
		dir := filepath.Dir(socketPath)
		// Standard name format for SAPControl UNIX socket: .sapstream5<SystemNumber>13
		newSocketPath := filepath.Join(dir, fmt.Sprintf(".sapstream5%s13", systemNumber))
		return "unix://" + newSocketPath, nil
	}

	u, err := url.Parse(origURL)
	if err != nil {
		return "", err
	}

	// Determine the port based on scheme
	var port string
	if u.Scheme == "https" {
		port = fmt.Sprintf("5%s14", systemNumber)
	} else {
		port = fmt.Sprintf("5%s13", systemNumber)
	}

	// Rebuild host with new port
	host := u.Hostname()
	u.Host = fmt.Sprintf("%s:%s", host, port)

	return u.String(), nil
}
