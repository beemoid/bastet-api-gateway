package models

import (
	"database/sql"
	"encoding/json"
)

// NullString is a custom type that handles NULL database values
// and marshals to JSON as either a string or null
type NullString struct {
	sql.NullString
}

// MarshalJSON implements the json.Marshaler interface
// Returns null if invalid, otherwise returns the string value
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements the json.Unmarshaler interface
// Handles both null and string values from JSON
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

// NullTime is a custom type that handles NULL database time values
// and marshals to JSON as either a time string or null
type NullTime struct {
	sql.NullTime
}

// MarshalJSON implements the json.Marshaler interface
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var t *string
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		nt.Valid = true
		// Parse the time string - adjust format as needed
		// This is a simplified version
	} else {
		nt.Valid = false
	}
	return nil
}
