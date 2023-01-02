//
// Copyright (c) 2016-2023 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package tracker

import (
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
)

type Subject struct {
	payload payload.Payload
}

// InitSubject returns a new subject object.
func InitSubject() *Subject {
	return &Subject{payload: *payload.Init()}
}

// Get returns the key-value store as a map[string]string.
func (s Subject) Get() map[string]string {
	return s.payload.Get()
}

// SetUserId adds a user id to teh key-value store.
func (s Subject) SetUserId(userId string) {
	s.payload.Add(UID, common.NewString(userId))
}

// SetScreenResolution adds the screen-resolution mesaurement to the key-value store.
func (s Subject) SetScreenResolution(width int, height int) {
	s.payload.Add(RESOLUTION, common.NewString(common.IntToString(width)+"x"+common.IntToString(height)))
}

// SetViewPort adds the view-port measurement to the key-value store.
func (s Subject) SetViewPort(width int, height int) {
	s.payload.Add(VIEWPORT, common.NewString(common.IntToString(width)+"x"+common.IntToString(height)))
}

// SetColorDepth adds the color-depth measurement to the key-value store.
func (s Subject) SetColorDepth(depth int) {
	s.payload.Add(COLOR_DEPTH, common.NewString(common.IntToString(depth)))
}

// SetTimeZone adds a timezone to the key-value store.
func (s Subject) SetTimeZone(timezone string) {
	s.payload.Add(TIMEZONE, common.NewString(timezone))
}

// SetLanguage adds a language to the key-value store.
func (s Subject) SetLanguage(language string) {
	s.payload.Add(LANGUAGE, common.NewString(language))
}

// SetIpAddress adds an ip address to the key-value store.
func (s Subject) SetIpAddress(ipAddress string) {
	s.payload.Add(IP_ADDRESS, common.NewString(ipAddress))
}

// SetUseragent adds a useragent to the key-value store.
func (s Subject) SetUseragent(useragent string) {
	s.payload.Add(USERAGENT, common.NewString(useragent))
}

// SetDomainUserId adds a domain user id to the key-value store.
func (s Subject) SetDomainUserId(domainUserId string) {
	s.payload.Add(DOMAIN_UID, common.NewString(domainUserId))
}

// SetNetworkUserId adds a network user id to the key-value store.
func (s Subject) SetNetworkUserId(networkUserId string) {
	s.payload.Add(NETWORK_UID, common.NewString(networkUserId))
}
