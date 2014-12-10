/*
Copyright (c) 2014-2015 Snowplow Analytics Ltd. All rights reserved.

This program is licensed to you under the Apache License Version 2.0,
and you may not use this file except in compliance with the Apache License
Version 2.0. You may obtain a copy of the Apache License Version 2.0 at

    http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing,
software distributed under the Apache License Version 2.0 is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the Apache License Version 2.0 for the specific
language governing permissions and limitations there under.
*/
package snowplowGo

import (
	"strconv"
)

const (
	DEFAULT_PLATFORM = "srv"
)

// TODO(alexanderdean): read this on interfaces: http://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go
// TODO(alexanderdean): turn the settings into a variable inside the Subject interface
var TrackerSettings map[string]string

// Initializes the Subject
func InitSubject() {
	TrackerSettings["p"] = DEFAULT_PLATFORM
}

// Sets the platform from which the event is fired
// TODO(AALEKH): remember about getSubject as golang variable type access
func SetPlatform(platform string) {
	TrackerSettings["p"] = platform
}

// Sets a custom user identification for the event
func SetUserId(userId string) {
	TrackerSettings["uid"] = userId
}

// Sets the screen resolution
func SetScreenResolution(width int, height int) {
	TrackerSettings["res"] = strconv.Itoa(width) + "x" + strconv.Itoa(height)
}

// Sets the view port resolution
func SetViewPort(width int, height int) {
	TrackerSettings["vp"] = strconv.Itoa(width) + "x" + strconv.Itoa(height)
}

// Sets the colour depth
func SetColorDepth(depth int) {
	TrackerSettings["cd"] = strconv.Itoa(depth)
}

// Sets the event timezone
func SetTimeZone(timezone string) {
	TrackerSettings["tz"] = timezone
}

// Sets the language used
func SetLanguage(language string) {
	TrackerSettings["lang"] = language
}
