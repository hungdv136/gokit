package types

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

// Map is a map from string to interface
type Map map[string]interface{}

// CreateMapFromStruct creates map from struct
func CreateMapFromStruct(obj interface{}) (Map, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return CreateMapFromReader(bytes.NewReader(b))
}

// CreateMapFromReader creates map from a JSON reader
func CreateMapFromReader(reader io.Reader) (Map, error) {
	m := Map{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&m); err != nil {
		if err == io.EOF {
			return m, nil
		}

		return nil, err
	}

	return m, nil
}

// ForceString returns value of key as a string, ignores error
func (m Map) ForceString(key string) string {
	val, ok := m.GetString(key)
	if !ok {
		log.Println("value of", key, "is not a string")
	}
	return val
}

// GetString returns value of a key as string
func (m Map) GetString(key string) (string, bool) {
	s, ok := m[key].(string)
	return s, ok
}

// JSON converts map to JSON string
func (m Map) JSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ForceJSON converts map to JSON string
func (m Map) ForceJSON() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(b)
}

// GetMap returns value of key as map
func (m Map) GetMap(key string) (Map, bool) {
	if v, ok := m[key].(map[string]interface{}); ok {
		return Map(v), true
	}

	if u, ok := m[key].(Map); ok {
		return u, true
	}

	return Map{}, false
}

// ForceMap returns value of key as a map, ignores error
func (m Map) ForceMap(key string) Map {
	val, ok := m.GetMap(key)
	if !ok {
		log.Println("value of", key, "is not a map")
	}
	return val
}
