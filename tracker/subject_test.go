//
// Copyright (c) 2016-2018 Snowplow Analytics Ltd. All rights reserved.
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
  "testing"
  "github.com/stretchr/testify/assert"
)

// TestSubjectSetFunctions asserts behaviour of all of the Subject setter functions.
func TestSubjectSetFunctions(t *testing.T) {
  assert := assert.New(t)
  subject := *InitSubject()

  subject.SetUserId("new-user-id")
  subject.SetScreenResolution(1920, 1080)
  subject.SetViewPort(1080, 1080)
  subject.SetColorDepth(1080)
  subject.SetTimeZone("ACT")
  subject.SetLanguage("EN")
  subject.SetIpAddress("127.0.0.1")
  subject.SetUseragent("useragent-string")
  subject.SetDomainUserId("domain-user-id")
  subject.SetNetworkUserId("network-user-id")

  subjectMap := subject.Get()
  assert.Equal("new-user-id", subjectMap[UID])
  assert.Equal("1920x1080", subjectMap[RESOLUTION])
  assert.Equal("1080x1080", subjectMap[VIEWPORT])
  assert.Equal("1080", subjectMap[COLOR_DEPTH])
  assert.Equal("ACT", subjectMap[TIMEZONE])
  assert.Equal("EN", subjectMap[LANGUAGE])
  assert.Equal("127.0.0.1", subjectMap[IP_ADDRESS])
  assert.Equal("useragent-string", subjectMap[USERAGENT])
  assert.Equal("domain-user-id", subjectMap[DOMAIN_UID])
  assert.Equal("network-user-id", subjectMap[NETWORK_UID])
}
