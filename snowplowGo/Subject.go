/*
Subject.go
Copyright (c) 2014 Snowplow Analytics Ltd. All rights reserved.
This program is licensed to you under the Apache License Version 2.0,
and you may not use this file except in compliance with the Apache License
Version 2.0. You may obtain a copy of the Apache License Version 2.0 at
http://www.apache.org/licenses/LICENSE-2.0.
Unless required by applicable law or agreed to in writing,
software distributed under the Apache License Version 2.0 is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the Apache License Version 2.0 for the specific
language governing permissions and limitations there under.
Authors: Aalekh Nigam
Copyright: Copyright (c) 2014 Snowplow Analytics Ltd
License: Apache License Version 2.0
*/
package snowplowGo

import (
	"strconv"
)

const (
	DEFAULT_PLATFORM = "srv"
)

var TrackerSettings map[string]string

//Initialize Subject.go
func InitSubject() {
	TrackerSettings["p"] = DEFAULT_PLATFORM
}

//Remember about getSubject as golang variable type access
 /**
* Sets the platform from which the event is fired
*
* @param string platform
*/
func SetPlatform(platform string) {
	TrackerSettings["p"] = platform
}

 /**
* Sets a custom user identification for the event
*
* @param string userId
*/

func SetUserId(userId string) {
	TrackerSettings["uid"] = userId
}

/**
* Sets the screen resolution
*
* @param int width
* @param int height
*/
func SetScreenResolution(width int, height int) {
	TrackerSettings["res"] = strconv.Itoa(width) + "x" + strconv.Itoa(height)
}

/**
* Sets the view port resolution
*
* @param int width
* @param int height
*/
func SetViewPort(width int, height int) {
	TrackerSettings["vp"] = strconv.Itoa(width) + "x" + strconv.Itoa(height)
}

 /**
* Sets the colour depth
*
* @param int depth
*/
func SetColorDepth(depth int) {
	TrackerSettings["cd"] = strconv.Itoa(depth)
}

/**
* Sets the event timezone
*
* @param string timezone
*/
func SetTimeZone(timezone string) {
	TrackerSettings["tz"] = timezone
}

 /**
* Sets the language used
*
* @param string language
*/
func SetLanguage(language string) {
	TrackerSettings["lang"] = language
}
