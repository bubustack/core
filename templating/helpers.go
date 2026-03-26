package templating

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const (
	typeNameArray  = "array"
	typeNameBool   = "bool"
	typeNameNull   = "null"
	typeNameNumber = "number"
	typeNameObject = "object"
	typeNameString = "string"
)

func lenValue(value any) (int, error) {
	if value == nil {
		return 0, nil
	}
	switch v := value.(type) {
	case string:
		return len(v), nil
	case []byte:
		return len(v), nil
	case []any:
		return len(v), nil
	case map[string]any:
		if ref, path, ok := storageSelectorFromStringMap(v); ok {
			return 0, &ErrOffloadedDataUsage{Reason: formatOffloadedReason("len", ref, path)}
		}
		return len(v), nil
	case map[any]any:
		if ref, path, ok := storageSelectorFromAnyMap(v); ok {
			return 0, &ErrOffloadedDataUsage{Reason: formatOffloadedReason("len", ref, path)}
		}
		return len(v), nil
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return 0, nil
	}
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return 0, nil
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return rv.Len(), nil
	default:
		return 0, fmt.Errorf("len: unsupported argument type %T", value)
	}
}

func hashOfValue(value any) (string, error) {
	if value == nil {
		return "", fmt.Errorf("hash_of: value is null")
	}
	if m, ok := value.(map[string]any); ok {
		if ref, path, ok := storageSelectorFromStringMap(m); ok {
			return "", &ErrOffloadedDataUsage{Reason: formatOffloadedReason("hash_of", ref, path)}
		}
	}
	if m, ok := value.(map[any]any); ok {
		if ref, path, ok := storageSelectorFromAnyMap(m); ok {
			return "", &ErrOffloadedDataUsage{Reason: formatOffloadedReason("hash_of", ref, path)}
		}
	}
	switch v := value.(type) {
	case string:
		return hashString(v), nil
	case []byte:
		return hashBytes(v), nil
	default:
		return "", fmt.Errorf("hash_of: unsupported argument type %T", value)
	}
}

func hashString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum)
}

func hashBytes(value []byte) string {
	sum := sha256.Sum256(value)
	return fmt.Sprintf("%x", sum)
}

func typeOfValue(value any) string {
	if value == nil {
		return typeNameNull
	}
	if offloaded, ok := offloadedTypeName(value); ok {
		return offloaded
	}
	if scalar := scalarTypeName(value); scalar != "" {
		return scalar
	}

	rv, ok := dereferenceTypeValue(value)
	if !ok {
		return typeNameNull
	}
	return reflectTypeName(rv)
}

func offloadedTypeName(value any) (string, bool) {
	switch typed := value.(type) {
	case map[string]any:
		if ref, path, ok := storageSelectorFromStringMap(typed); ok {
			return fmt.Sprintf("offloaded(%s%s)", ref, formatOffloadedPath(path)), true
		}
	case map[any]any:
		if ref, path, ok := storageSelectorFromAnyMap(typed); ok {
			return fmt.Sprintf("offloaded(%s%s)", ref, formatOffloadedPath(path)), true
		}
	}
	return "", false
}

func scalarTypeName(value any) string {
	switch value.(type) {
	case map[string]any, map[any]any:
		return typeNameObject
	case []any:
		return typeNameArray
	case string:
		return typeNameString
	case bool:
		return typeNameBool
	case float32, float64, int, int8, int16, int32, int64:
		return typeNameNumber
	case uint, uint8, uint16, uint32, uint64:
		return typeNameNumber
	}
	return ""
}

func dereferenceTypeValue(value any) (reflect.Value, bool) {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return reflect.Value{}, false
	}
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return reflect.Value{}, false
		}
		rv = rv.Elem()
	}
	return rv, true
}

func reflectTypeName(rv reflect.Value) string {
	switch rv.Kind() {
	case reflect.Map, reflect.Struct:
		return typeNameObject
	case reflect.Slice, reflect.Array:
		return typeNameArray
	case reflect.String:
		return typeNameString
	case reflect.Bool:
		return typeNameBool
	default:
		if isNumericKind(rv.Kind()) {
			return typeNameNumber
		}
		return rv.Kind().String()
	}
}

func isNumericKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func sampleValue(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	if m, ok := value.(map[string]any); ok {
		if ref, path, ok := storageSelectorFromStringMap(m); ok {
			return nil, &ErrOffloadedDataUsage{Reason: formatOffloadedReason("sample", ref, path)}
		}
	}
	if m, ok := value.(map[any]any); ok {
		if ref, path, ok := storageSelectorFromAnyMap(m); ok {
			return nil, &ErrOffloadedDataUsage{Reason: formatOffloadedReason("sample", ref, path)}
		}
	}
	return value, nil
}

func storageSelectorFromStringMap(m map[string]any) (string, string, bool) {
	if m == nil {
		return "", "", false
	}
	refRaw, ok := m[StorageRefKey]
	if !ok {
		return "", "", false
	}
	ref, ok := refRaw.(string)
	if !ok || strings.TrimSpace(ref) == "" {
		return "", "", false
	}
	path, _ := m[StoragePathKey].(string)
	return ref, path, true
}

func storageSelectorFromAnyMap(m map[any]any) (string, string, bool) {
	if m == nil {
		return "", "", false
	}
	refRaw, ok := m[StorageRefKey]
	if !ok {
		return "", "", false
	}
	ref, ok := refRaw.(string)
	if !ok || strings.TrimSpace(ref) == "" {
		return "", "", false
	}
	path, _ := m[StoragePathKey].(string)
	return ref, path, true
}

func formatOffloadedReason(op, ref, path string) string {
	if ref == "" {
		return fmt.Sprintf("%s requires hydrated data", op)
	}
	if path == "" {
		return fmt.Sprintf("%s requires hydrated data for %s", op, ref)
	}
	return fmt.Sprintf("%s requires hydrated data for %s%s", op, ref, formatOffloadedPath(path))
}

func formatOffloadedPath(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	if strings.HasPrefix(path, ".") || strings.HasPrefix(path, "[") {
		return path
	}
	return "." + path
}

func parseJSONValue(raw string) (any, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, false
	}
	var out any
	if err := json.Unmarshal([]byte(trimmed), &out); err != nil {
		return nil, false
	}
	return out, true
}
