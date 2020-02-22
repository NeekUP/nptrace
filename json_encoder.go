package nptrace

import (
	"fmt"
	"strconv"
	"time"
)

func NewJsonEncoderConfig(timeFormat string, durationFormat func(d time.Duration) []byte) *JsonEncoderConfig {
	return &JsonEncoderConfig{
		TimeFormat:     timeFormat,
		DurationFormat: durationFormat,
	}
}

type JsonEncoderConfig struct {
	TimeFormat     string
	DurationFormat func(d time.Duration) []byte
}

func NewJsonEncoder(cfg *JsonEncoderConfig) Encoder {
	return &jsonEncoder{
		cfg: cfg,
		buf: []byte{},
	}
}

type jsonEncoder struct {
	buf []byte
	cfg *JsonEncoderConfig
}

func (enc *jsonEncoder) Encode(npt *Task) []byte {
	if npt == nil {
		return []byte{}
	}
	enc.append('{')
	enc.addToken("id", npt.Id)
	enc.append(',')
	enc.addToken("time", npt.Time)
	enc.append(',')
	enc.addToken("trace", npt.Trace)
	enc.append('}')
	enc.append('\n')

	return enc.buf
}

func (enc *jsonEncoder) encodeTrace(trace *Trace) {
	if trace == nil {
		return
	}
	enc.append('{')
	enc.addToken("name", trace.Name)
	enc.append(',')
	enc.addToken("duration", trace.Duration)
	enc.append(',')
	enc.addToken("args", trace.Args)
	enc.append(',')
	enc.addToken("traces", trace.Children)
	enc.append('}')
}

func (enc *jsonEncoder) addToken(name string, value interface{}) {
	enc.addString(name)
	enc.append(':')
	enc.addValue(value)
}

func (enc *jsonEncoder) addValue(value interface{}) {
	if value == nil {
		enc.addString("")
		return
	}

	switch v := value.(type) {
	case bool:
		enc.buf = strconv.AppendBool(enc.buf, v)
	case int:
		enc.buf = strconv.AppendInt(enc.buf, int64(v), 10)
	case int8:
		enc.buf = strconv.AppendInt(enc.buf, int64(v), 10)
	case int16:
		enc.buf = strconv.AppendInt(enc.buf, int64(v), 10)
	case int32:
		enc.buf = strconv.AppendInt(enc.buf, int64(v), 10)
	case int64:
		enc.buf = strconv.AppendInt(enc.buf, v, 10)
	case uint:
		enc.buf = strconv.AppendUint(enc.buf, uint64(v), 10)
	case uintptr:
		enc.buf = strconv.AppendUint(enc.buf, uint64(v), 10)
	case uint8:
		enc.buf = strconv.AppendUint(enc.buf, uint64(v), 10)
	case uint16:
		enc.buf = strconv.AppendUint(enc.buf, uint64(v), 10)
	case uint32:
		enc.buf = strconv.AppendUint(enc.buf, uint64(v), 10)
	case uint64:
		enc.buf = strconv.AppendUint(enc.buf, v, 10)
	case float32:
		enc.buf = strconv.AppendFloat(enc.buf, float64(v), 'f', -1, 32)
	case float64:
		enc.buf = strconv.AppendFloat(enc.buf, float64(v), 'f', -1, 64)
	case complex64:
		enc.appendComplex(complex128(v))
	case complex128:
		enc.appendComplex(v)
	case string:
		enc.addString(v)
	case time.Time:
		enc.addTime(v)
	case time.Duration:
		enc.appendBytes(enc.cfg.DurationFormat(v))
	case *Trace:
		enc.encodeTrace(v)
	case []*Trace:
		enc.append('[')
		for i := 0; i < len(v); i++ {
			enc.encodeTrace(v[i])
			if i != len(v)-1 {
				enc.append(',')
			}
		}
		enc.append(']')
	case []interface{}:
		enc.append('[')
		for i := 0; i < len(v); i++ {
			enc.addValue(v[i])
			if i != len(v)-1 {
				enc.append(',')
			}
		}
		enc.append(']')
	default:
		enc.addString(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (enc *jsonEncoder) appendComplex(v complex128) {
	r, i := float64(real(v)), float64(imag(v))
	enc.append('"')
	enc.buf = strconv.AppendFloat(enc.buf, float64(r), 'f', -1, 64)
	enc.append('+')
	enc.buf = strconv.AppendFloat(enc.buf, float64(i), 'f', -1, 64)
	enc.append('i')
	enc.append('"')
}

func (enc *jsonEncoder) addString(v string) {
	enc.append('"')
	enc.appendBytes([]byte(v))
	enc.append('"')
}

func (enc *jsonEncoder) addTime(v time.Time) {
	enc.append('"')
	enc.buf = v.AppendFormat(enc.buf, enc.cfg.TimeFormat)
	enc.append('"')
}

func (enc *jsonEncoder) append(v byte) {
	enc.buf = append(enc.buf, v)
}

func (enc *jsonEncoder) appendBytes(v []byte) {
	enc.buf = append(enc.buf, v...)
}
