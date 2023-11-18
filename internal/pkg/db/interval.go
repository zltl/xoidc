package db

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// PgDuration wraps a time.Duration to provide implementations of
// sql.Scanner and driver.Valuer for reading/writing from/to a DB.
type PgDuration time.Duration

// Value converts the PgDuration into a string.
func (d PgDuration) Value() (driver.Value, error) {
	return time.Duration(d).String(), nil
}

// Scan converts the received string in the format hh:mm:ss into a PgDuration.
func (d *PgDuration) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		// Convert format of hh:mm:ss into format parseable by time.ParseDuration()
		v = strings.Replace(v, ":", "h", 1)
		v = strings.Replace(v, ":", "m", 1)
		v += "s"
		dur, err := time.ParseDuration(v)
		if err != nil {
			return err
		}
		*d = PgDuration(dur)
		return nil
	default:
		return fmt.Errorf("cannot sql.Scan() PgDuration from: %#v", v)
	}
}
