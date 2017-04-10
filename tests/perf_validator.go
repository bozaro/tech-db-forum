package tests

type PerfValidator interface {
	CheckBetween(min int, val int, max int, message string)
}
