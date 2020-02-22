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
	Id      string
	Time    time.Time
	Trace   *Trace
	current *Trace
}

func (npt *NPTrace) New(id string, name string, arguments ...interface{}) *Task {
	now := time.Now()
	tsk := &Task{
		Id:   id,
		Time: now,
		Trace: &Trace{
			Name: name,
			Time: now,
			Args: arguments,
		},
	}

	tsk.current = tsk.Trace
	return tsk
}

func (npt *NPTrace) Close(t *Task) (bool, error) {
	if t == nil || t.Trace == nil {
		return false, nil
	}
	t.Trace.Duration = time.Now().Sub(t.Time)
	b := npt.encoder.Encode(t)
	n, err := npt.writer.Write(b)
	return n > 0, err
}

type Trace struct {
	Name     string
	Time     time.Time
	Duration time.Duration
	Children []*Trace
	parent   *Trace
	Args     []interface{}
}

func (tsk *Task) Start(name string, arguments ...interface{}) *Trace {
	tr := &Trace{
		Name:   name,
		Time:   time.Now(),
		Args:   arguments,
		parent: tsk.current,
	}

	if tsk.current.Children == nil {
		tsk.current.Children = []*Trace{tr}
	} else {
		tsk.current.Children = append(tsk.current.Children, tr)
	}

	tsk.current = tr
	return tr
}

func (tsk *Task) Stop(t *Trace) {
	t.Duration = time.Now().Sub(t.Time)
	tsk.current = t.parent
}

func (t *Trace) Point(name string, arguments ...interface{}) {
	p := &Trace{
		Name:   name,
		Args:   arguments,
		parent: t,
	}

	if t.Children == nil {
		t.Children = []*Trace{}
		p.Time = t.Time
	} else {
		last := t.Children[len(t.Children)-1]
		p.Time = last.Time.Add(last.Duration)
	}

	p.Duration = time.Now().Sub(p.Time)

	t.Children = append(t.Children, p)
}
