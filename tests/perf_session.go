package tests

type PerfSession struct {
}

func (self *PerfSession) CheckInt(expected int, actual int, message string) {
	if expected != actual {
		panic(message)
	}
}

func (self *PerfSession) Finish(before PVersion, after PVersion) {

}
