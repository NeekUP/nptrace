package nptrace

import (
	"testing"
	"time"
)

func TestStructure(t *testing.T) {
	npt := NewTracer(&fakeEncoder{}, &fakeWriter{})
	task := npt.New("test", "root")
	f := task.Start("f")
	s := task.Start("s")
	s.Point("s1")
	f.Point("f1")
	s.Point("s2")
	f.Point("f2")
	task.Stop(s)
	task.Stop(f)
	npt.Close(task)

	rt := task.trace

	if rt.name != "root" {
		t.Error("first trace not expected")
		return
	}

	if rt.parent != nil {
		t.Error("first trace has parent")
		return
	}

	ft := rt.children[0]
	if rt.children == nil || len(rt.children) != 1 || ft.children[0].name != "s" || ft.children[1].name != "f1" || ft.children[2].name != "f2" {
		t.Error("unexpected children of first trace")
		return
	}

	st := rt.children[0].children[0]
	if len(st.children) != 2 || st.children[0].name != "s1" || st.children[1].name != "s2" {
		t.Error("unexpected children of second trace")
		return
	}
}

func TestTrace(t *testing.T) {
	npt := NewTracer(&fakeEncoder{}, &fakeWriter{})
	task := npt.New("test", "root", "1", "2", "3", "4")
	task.Stop(task.trace)

	tr := task.trace
	if tr.name != "root" {
		t.Error("Unexpected name")
	}

	if tr.start.After(time.Now()) || tr.start.Before(time.Now().Add(-10*time.Millisecond)) {
		t.Error("Start time not valid")
	}

	if tr.duration == 0 || tr.duration > 10*time.Millisecond {
		t.Error("duration not valid")
	}

	if len(tr.args) != 4 {
		t.Error("Unexpected arguments count ")
	}
}

func TestPoint(t *testing.T) {
	npt := NewTracer(&fakeEncoder{}, &fakeWriter{})
	task := npt.New("test", "root", "1", "2", "3", "4")

	tr := task.trace
	tr.Point("point", "1", "2")

	task.Stop(task.trace)

	point := tr.children[0]
	if point.name != "point" {
		t.Error("Unexpected name")
	}

	if point.start.After(time.Now()) || point.start.Before(time.Now().Add(-10*time.Millisecond)) {
		t.Error("Start time not valid")
	}

	if point.duration == 0 || point.duration > 10*time.Millisecond {
		t.Error("duration not valid")
	}

	if len(point.args) != 2 {
		t.Error("Unexpected arguments count ")
	}
}

type fakeEncoder struct {
}

func (enc fakeEncoder) Encode(t *Task) []byte {
	return []byte{}
}

type fakeWriter struct {
}

func (wr fakeWriter) Write(p []byte) (n int, err error) {
	return 1, nil
}

func (wr fakeWriter) Lock() {
}

func (wr fakeWriter) Unock() {
}
