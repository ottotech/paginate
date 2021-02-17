package paginate

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type NullInt struct {
	Int   int
	Valid bool // Valid is true if Int is not NULL
}

func (ni *NullInt) Scan(value interface{}) error {
	if value == nil {
		ni.Int, ni.Valid = 0, false
		return nil
	}
	intVal := 0
	switch t := value.(type) {
	case int:
		intVal = t
	case int8:
		intVal = int(t)
	case int16:
		intVal = int(t)
	case int32:
		intVal = int(t)
	case int64:
		intVal = int(t)
	case uint:
		intVal = int(t)
	case uint8:
		intVal = int(t)
	case uint16:
		intVal = int(t)
	case uint32:
		intVal = int(t)
	case uint64:
		intVal = int(t)
	// As an special case when passing custom types that implement the Scanner interface to the
	// driver "go-sql-driver/mysql", the driver will return []uint8. Therefore, we need to try to parse that returned
	// value to string and finally to int in this case if possible.
	// Check issue: https://github.com/go-sql-driver/mysql/issues/441
	case []uint8:
		_int, err := strconv.Atoi(string(t))
		if err != nil {
			return err
		}
		intVal = _int
	default:
		return fmt.Errorf("column is not int")
	}
	ni.Valid = true
	ni.Int = intVal
	return nil
}

func (ni NullInt) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int, nil
}

func (ni NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int)
}

func (ni *NullInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Int, ni.Valid = 0, false
		return nil
	}
	intVal, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	ni.Int, ni.Valid = intVal, true
	return nil
}

type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

func (nb *NullBool) Scan(value interface{}) error {
	if value == nil {
		nb.Bool, nb.Valid = false, false
		return nil
	}

	boolVal := false
	switch t := value.(type) {
	case bool:
		boolVal = t
	// As an special case when passing custom types that implement the Scanner interface to the
	// driver "go-sql-driver/mysql", the driver will return []uint8. Therefore, we need to try to parse that returned
	// value to string and finally to bool in this case if possible.
	// Check issue: https://github.com/go-sql-driver/mysql/issues/441
	case []uint8:
		_bool, err := strconv.ParseBool(string(t))
		if err != nil {
			return err
		}
		boolVal = _bool
	case int64:
		_bool, err := strconv.ParseBool(strconv.FormatInt(t, 10))
		if err != nil {
			return err
		}
		boolVal = _bool
	default:
		return fmt.Errorf("column is not boolean")
	}
	nb.Valid = true
	nb.Bool = boolVal
	return nil
}

func (nb NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

func (nb *NullBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nb.Bool, nb.Valid = false, false
		return nil
	}
	boolVal, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	nb.Bool, nb.Valid = boolVal, true
	return nil
}

type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	stringVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("column is not string")
	}
	ns.Valid = true
	ns.String = stringVal
	return nil
}

func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.String, ns.Valid = string(data), true
	return nil
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}
	timeVal, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("column is not timestamp")
	}
	nt.Valid = true
	nt.Time = timeVal
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}

	// Here we unmarshall to json the time object as defined in
	// https://golang.org/pkg/time/#Time.UnmarshalJSON.
	t, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	if err != nil {
		return err
	}

	nt.Time, nt.Valid = t, true
	return nil
}

type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

func (n *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		n.Float64, n.Valid = 0, false
		return nil
	}
	var float64Val float64
	switch t := value.(type) {
	case float32:
		float64Val = float64(t)
	case float64:
		float64Val = t
	default:
		return fmt.Errorf("column is not float64")
	}
	n.Valid = true
	n.Float64 = float64Val
	return nil
}

func (n NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Float64)
}

func (n *NullFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Float64, n.Valid = 0, false
		return nil
	}

	float64Val, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}

	n.Float64, n.Valid = float64Val, true
	return nil
}
