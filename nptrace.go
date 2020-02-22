package nptrace

import (
	"io"
	"time"
)

type NPTrace struct {
	encoder Encoder
	writer  WriterSync
}

func NewTracer(enc Encoder, writer io.Writer) *NPTrace {
	return &NPTrace{
		encoder: enc,
		writer:  NewWriterSync(writer),
	}
}

type Task struct {
	id      string
	start   time.Time
	trace   *Trace
	current *Trace
}

func (npt *NPTrace) New(id string, name string, arguments ...interface{}) *Task {
	now := time.Now()
	tsk := &Task{
		id:    id,
		start: now,
		trace: &Trace{
			name:  name,
			start: now,
			args:  arguments,
		},
	}

	tsk.current = tsk.trace
	return tsk
}

func (npt *NPTrace) Close(t *Task) (bool, error) {
	if t == nil || t.trace == nil {
		return false, nil
	}
	t.trace.duration = time.Now().Sub(t.start)
	b := npt.encoder.Encode(t)
	n, err := npt.writer.Write(b)
	return n > 0, err
}

type Trace struct {
	name     string
	start    time.Time
	duration time.Duration
	children []*Trace
	parent   *Trace
	args     []interface{}
}

func (tsk *Task) Start(name string, arguments ...interface{}) *Trace {
	tr := &Trace{
		name:   name,
		start:  time.Now(),
		args:   arguments,
		parent: tsk.current,
	}

	if tsk.current.children == nil {
		tsk.current.children = []*Trace{tr}
	} else {
		tsk.current.children = append(tsk.current.children, tr)
	}

	tsk.current = tr
	return tr
}

func (tsk *Task) Stop(t *Trace) {
	t.duration = time.Now().Sub(t.start)
	tsk.current = t.parent
}

func (t *Trace) Point(name string, arguments ...interface{}) {
	p := &Trace{
		name:   name,
		args:   arguments,
		parent: t,
	}

	if t.children == nil {
		t.children = []*Trace{}
		p.start = t.start
	} else {
		last := t.children[len(t.children)-1]
		p.start = last.start.Add(last.duration)
	}

	p.duration = time.Now().Sub(p.start)

	t.children = append(t.children, p)
}
