package tests

type PerfValidator interface {
	CheckInt(expected int, actual int, message string)
	Finish(before PVersion, after PVersion)
}
