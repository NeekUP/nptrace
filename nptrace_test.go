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

	rt := task.Trace

	if rt.Name != "root" {
		t.Error("first Trace not expected")
		return
	}

	if rt.parent != nil {
		t.Error("first Trace has parent")
		return
	}

	ft := rt.Children[0]
	if rt.Children == nil || len(rt.Children) != 1 || ft.Children[0].Name != "s" || ft.Children[1].Name != "f1" || ft.Children[2].Name != "f2" {
		t.Error("unexpected Children of first Trace")
		return
	}

	st := rt.Children[0].Children[0]
	if len(st.Children) != 2 || st.Children[0].Name != "s1" || st.Children[1].Name != "s2" {
		t.Error("unexpected Children of second Trace")
		return
	}
}

func TestTrace(t *testing.T) {
	npt := NewTracer(&fakeEncoder{}, &fakeWriter{})
	task := npt.New("test", "root", "1", "2", "3", "4")
	task.Stop(task.Trace)

	tr := task.Trace
	if tr.Name != "root" {
		t.Error("Unexpected Name")
	}

	if tr.Time.After(time.Now()) || tr.Time.Before(time.Now().Add(-10*time.Millisecond)) {
		t.Error("Start time not valid")
	}

	if tr.Duration == 0 || tr.Duration > 10*time.Millisecond {
		t.Error("Duration not valid")
	}

	if len(tr.Args) != 4 {
		t.Error("Unexpected arguments count ")
	}
}

func TestPoint(t *testing.T) {
	npt := NewTracer(&fakeEncoder{}, &fakeWriter{})
	task := npt.New("test", "root", "1", "2", "3", "4")

	tr := task.Trace
	tr.Point("point", "1", "2")

	task.Stop(task.Trace)

	point := tr.Children[0]
	if point.Name != "point" {
		t.Error("Unexpected Name")
	}

	if point.Time.After(time.Now()) || point.Time.Before(time.Now().Add(-10*time.Millisecond)) {
		t.Error("Start time not valid")
	}

	if point.Duration == 0 || point.Duration > 10*time.Millisecond {
		t.Error("Duration not valid")
	}

	if len(point.Args) != 2 {
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
