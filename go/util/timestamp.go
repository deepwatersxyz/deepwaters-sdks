package util

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Now() uint64 {
	return uint64(time.Now().UnixMicro())
}

type Timestamp struct {
	Time             time.Time
	MicrosSinceEpoch uint64
	utc              bool
}

func NewTimestampFromMicros(micros *uint64, utc bool) *Timestamp {
	if micros == nil {
		return nil
	}
	ts := &Timestamp{utc: utc}
	ts.SetMicros(*micros)
	return ts
}

func NewTimestampFromTime(tm time.Time, utc bool) *Timestamp {
	ts := &Timestamp{utc: utc}
	ts.SetTime(tm)
	return ts
}

func (t *Timestamp) SetTime(tm time.Time) {
	t.SetMicros(uint64(tm.UnixMicro()))
}

func (t *Timestamp) SetMicros(micros uint64) {
	t.MicrosSinceEpoch = micros
	t.Time = time.Unix(int64(micros/1000000), int64(micros%1000000)*1000)
	if t.utc {
		t.Time = t.Time.UTC()
	}
}

func (t *Timestamp) Equals(ts Timestamp) bool {
	return t.MicrosSinceEpoch == ts.MicrosSinceEpoch
}

func (t *Timestamp) Before(ts Timestamp) bool {
	return t.MicrosSinceEpoch < ts.MicrosSinceEpoch
}

func (t *Timestamp) After(ts Timestamp) bool {
	return t.MicrosSinceEpoch > ts.MicrosSinceEpoch
}

func (t *Timestamp) GetTimeStr() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%+v", t.Time)
}

func (t *Timestamp) GetRFC3339NanoTimeStr() string {
	if t == nil {
		return "<nil>"
	}
	return t.Time.Format(time.RFC3339Nano)
}

func (t *Timestamp) GetMicrosStr(pretty bool) string {
	if t == nil {
		return "<nil>"
	}
	if pretty {
		return PrettyNum(t.MicrosSinceEpoch)
	}
	return fmt.Sprintf("%+v", t.MicrosSinceEpoch)
}

func (t *Timestamp) GetMicrosP() *uint64 {
	if t == nil {
		return nil
	}
	return &t.MicrosSinceEpoch
}

func (t *Timestamp) Clone() *Timestamp {
	return NewTimestampFromMicros(&t.MicrosSinceEpoch, t.utc)
}

func (t *Timestamp) Add(d time.Duration) *Timestamp {
	ts := t.Clone()
	ts.SetTime(ts.Time.Add(d))
	return ts
}

func (t *Timestamp) Sub(ts Timestamp) time.Duration {
	return t.Time.Sub(ts.Time)
}

// ConvertUint64Value follows a similar logic to driver.Int.ConvertValue(v)
func ConvertUint64Value(v interface{}) (*uint64, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := rv.Int()
		if i64 < 0 {
			return nil, fmt.Errorf("sql/driver: value %d is negative", v)
		}
		u64 := uint64(i64)
		return &u64, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := rv.Uint()
		return &u64, nil
	case reflect.String:
		i, err := strconv.ParseUint(rv.String(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("sql/driver: value %q can't be converted to uint64", v)
		}
		return &i, nil
	}
	return nil, fmt.Errorf("sql/driver: unsupported value %v (type %T) converting to uint64", v, v)
}

func (t *Timestamp) Scan(src interface{}) error {
	if src == nil {
		t.SetMicros(0)
		return nil
	}

	u64, err := ConvertUint64Value(src)
	if err != nil {
		return err
	}

	t.SetMicros(*u64)
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	return int64(t.MicrosSinceEpoch), nil
}
