package nptrace

import (
	"testing"
	"time"
)

func TestJsonEncoder_Encode(t *testing.T) {
	cfg := &JsonEncoderConfig{TimeFormat: time.RFC3339}
	encoder := NewJsonEncoder(cfg)

	npt := NewTracer(nil, nil)
	task := npt.New("test", "root")
	task.start = time.Date(2010, 01, 01, 02, 02, 03, 3, time.UTC)
	task.trace.start = time.Date(2010, 01, 01, 02, 02, 03, 4, time.UTC)

	s := &Trace{
		name:     "f 1",
		start:    time.Date(2010, 01, 01, 02, 02, 03, 4, time.UTC),
		duration: 1001 * time.Nanosecond,
		args:     []interface{}{"1", "ds", 1, time.Date(2010, 01, 01, 02, 02, 03, 5, time.UTC), true, complex(0.1, 0.2), 1 << 35},
	}
	task.trace.children = []*Trace{}
	task.trace.children = append(task.trace.children, s)

	expected := []byte(`{"id":"test","time":"2010-01-01T02:02:03Z","trace":{"name":"root","duration":0,"args":[],"traces":[{"name":"f 1","duration":1001,"args":["1","ds",1,"2010-01-01T02:02:03Z",true,"0.1+0.2i",34359738368],"traces":[]}]}}`)
	result := encoder.Encode(task)
	for i := 0; i < len(result); i++ {
		if result[i] != expected[i] {
			t.Errorf("Position %d. result:%s, expected:%s", i, string(result[i]), string(expected[i]))
		}
	}
}
