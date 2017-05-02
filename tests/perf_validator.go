package tests

import "github.com/go-openapi/strfmt"

type PerfValidator interface {
	CheckVersion(before PVersion, after PVersion) bool
	CheckInt(expected int, actual int, message string)
	CheckInt32(expected int32, actual int32, message string)
	CheckInt64(expected int64, actual int64, message string)
	CheckStr(expected string, actual string, message string)
	CheckBool(expected bool, actual bool, message string)
	CheckHash(expected PHash, actual string, message string)
	CheckDate(expected *strfmt.DateTime, actual *strfmt.DateTime, message string)
	Finish(before PVersion, after PVersion)
}
