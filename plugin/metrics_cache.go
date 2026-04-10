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
	"sync"
	"time"
)

// MetricResult stores the actual data extracted from XML
type MetricResult struct {
	Value     string
	Timestamp string
}

// HostCache holds multiple data buckets for a single SAP Host
type HostCache struct {
	// Buckets: map[RequestType]map[MetricKey]MetricResult
	Buckets    map[string]map[string]MetricResult
	LastUpdate map[string]time.Time
}

// SessionCache manages caches for all different SAP hosts
type SessionCache struct {
	mu       sync.Mutex
	ttl      time.Duration
	sessions map[string]*HostCache
}

func NewSessionCache(defaultTTL time.Duration) *SessionCache {
	return &SessionCache{
		ttl:      defaultTTL,
		sessions: make(map[string]*HostCache),
	}
}

// GetCachedData returns the data bucket for a specific host and request type if not expired
func (sc *SessionCache) GetCachedData(sessionKey, requestType string) (map[string]MetricResult, bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	session, exists := sc.sessions[sessionKey]
	if !exists {
		return nil, false
	}

	updateTime, ok := session.LastUpdate[requestType]
	if !ok || time.Since(updateTime) > sc.ttl {
		return nil, false
	}

	return session.Buckets[requestType], true
}

// UpdateSession saves a refreshed data bucket
func (sc *SessionCache) UpdateSession(sessionKey, requestType string, data map[string]MetricResult) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if _, exists := sc.sessions[sessionKey]; !exists {
		sc.sessions[sessionKey] = &HostCache{
			Buckets:    make(map[string]map[string]MetricResult),
			LastUpdate: make(map[string]time.Time),
		}
	}

	sc.sessions[sessionKey].Buckets[requestType] = data
	sc.sessions[sessionKey].LastUpdate[requestType] = time.Now()
}

// Cleanup removes stale entries from the cache
func (sc *SessionCache) Cleanup() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for sessionKey, session := range sc.sessions {
		for reqType, updateTime := range session.LastUpdate {
			if time.Since(updateTime) > sc.ttl {
				delete(session.Buckets, reqType)
				delete(session.LastUpdate, reqType)
			}
		}

		if len(session.Buckets) == 0 {
			delete(sc.sessions, sessionKey)
		}
	}
}
