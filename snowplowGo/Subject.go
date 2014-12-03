package snowplowGo

import (
	"strconv"
)

const (
	DEFAULT_PLATFORM = "srv"
)

var TrackerSettings map[string]string

func InitSubject() {
	TrackerSettings["p"] = DEFAULT_PLATFORM
}

//Remember about getSubject as golang variable type access

func SetPlatform(platform string) {
	TrackerSettings["p"] = platform
}

func SetUserId(userId string) {
	TrackerSettings["uid"] = userId
}

func SetScreenResolution(width int, height int) {
	TrackerSettings["res"] = strconv.Itoa(width) + "x" + strconv.Itoa(height)
}

func SetViewPort(width int, height int) {
	TrackerSettings["vp"] = strconv.Itoa(width) + "x" + strconv.Itoa(height)
}

func SetColorDepth(depth int) {
	TrackerSettings["cd"] = strconv.Itoa(depth)
}

func SetTimeZone(timezone string) {
	TrackerSettings["tz"] = timezone
}

func SetLanguage(language string) {
	TrackerSettings["lang"] = language
}
