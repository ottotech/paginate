package paginate

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
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
	Valid bool // Valid is true if Int is not NULL
}

func (nb *NullBool) Scan(value interface{}) error {
	if value == nil {
		nb.Bool, nb.Valid = false, false
		return nil
	}
	boolVal, ok := value.(bool)
	if !ok {
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
