package nptrace

type Encoder interface {
	Encode(t *Task) []byte
}
