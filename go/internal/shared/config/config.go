package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/joho/godotenv"
)

func Load(cfg any, files ...string) error {
	err := godotenv.Load(files...)
	if err != nil {
		return err
	}
	return loadEnv(cfg, "")
}

func loadEnv(cfg any, prefix string) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("cfg must be a non-nil pointer to struct")
	}
	v = v.Elem()
	t := v.Type()

	durType := reflect.TypeOf(time.Duration(0))

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		ft := t.Field(i)

		isPtr := field.Kind() == reflect.Pointer
		baseType := field.Type()
		if isPtr {
			baseType = baseType.Elem()
		}

		// ----- Structs (and *Struct) -----
		if baseType.Kind() == reflect.Struct && baseType != durType {
			newPrefix := prefix + getEnvName(ft.Name+"_")
			if !isPtr {
				if err := loadEnv(field.Addr().Interface(), newPrefix); err != nil {
					return err
				}
				continue
			}
			inst := reflect.New(baseType)
			if err := loadEnv(inst.Interface(), newPrefix); err != nil {
				field.Set(reflect.Zero(field.Type())) // -> nil
				continue
			}
			field.Set(inst)
			continue
		}

		envName := prefix + getEnvName(ft.Name)
		val, ok := os.LookupEnv(envName)
		if !ok {
			val = ft.Tag.Get("env_default")
		}

		// Required unless it's a pointer
		if val == "" {
			if isPtr {
				field.Set(reflect.Zero(field.Type())) // -> nil
				continue
			}
			return fmt.Errorf("missing required env var %s for field %s", envName, ft.Name)
		}

		target := field
		if isPtr {
			if field.IsNil() {
				field.Set(reflect.New(baseType))
			}
			target = field.Elem()
		}

		// ----- Slices (and *[]T) -----
		if baseType.Kind() == reflect.Slice {
			sliceVal, err := parseSlice(val, baseType.Elem(), durType)
			if err != nil {
				if isPtr {
					field.Set(reflect.Zero(field.Type())) // -> nil
					continue
				}
				return fmt.Errorf("%s: %w", envName, err)
			}
			target.Set(sliceVal)
			continue
		}

		// ----- Scalars -----
		switch baseType {
		case reflect.TypeOf(""):
			target.SetString(val)

		case reflect.TypeOf(0):
			i64, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				if isPtr {
					field.Set(reflect.Zero(field.Type()))
					continue
				}
				return fmt.Errorf("%s: %w", envName, err)
			}
			target.SetInt(i64)

		case reflect.TypeOf(0.0):
			f64, err := strconv.ParseFloat(val, 64)
			if err != nil {
				if isPtr {
					field.Set(reflect.Zero(field.Type()))
					continue
				}
				return fmt.Errorf("%s: %w", envName, err)
			}
			target.SetFloat(f64)

		case reflect.TypeOf(true):
			b, err := strconv.ParseBool(val)
			if err != nil {
				if isPtr {
					field.Set(reflect.Zero(field.Type()))
					continue
				}
				return fmt.Errorf("%s: %w", envName, err)
			}
			target.SetBool(b)

		case durType:
			d, err := time.ParseDuration(val)
			if err != nil {
				if isPtr {
					field.Set(reflect.Zero(field.Type()))
					continue
				}
				return fmt.Errorf("%s: %w", envName, err)
			}
			target.SetInt(int64(d))
		}
	}
	return nil
}

// --- Helpers ---

func parseSlice(val string, elemType, durType reflect.Type) (reflect.Value, error) {
	val = strings.TrimSpace(val)

	// JSON array support
	if strings.HasPrefix(val, "[") {
		dst := reflect.New(reflect.SliceOf(elemType)).Interface()
		if err := json.Unmarshal([]byte(val), dst); err == nil {
			return reflect.ValueOf(dst).Elem(), nil
		}
		// fall through to CSV parsing on JSON failure
	}

	parts := splitCSV(val)
	out := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			// treat empty element as parse error
			return reflect.Value{}, fmt.Errorf("empty element in list")
		}
		ev, err := parseElem(p, elemType, durType)
		if err != nil {
			return reflect.Value{}, err
		}
		out = reflect.Append(out, ev)
	}
	return out, nil
}

func parseElem(s string, typ, durType reflect.Type) (reflect.Value, error) {
	switch typ {
	case reflect.TypeOf(""):
		return reflect.ValueOf(s), nil
	case reflect.TypeOf(0):
		i64, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		v := reflect.New(typ).Elem()
		v.SetInt(i64)
		return v, nil
	case reflect.TypeOf(0.0):
		f64, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		v := reflect.New(typ).Elem()
		v.SetFloat(f64)
		return v, nil
	case reflect.TypeOf(true):
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, err
		}
		v := reflect.New(typ).Elem()
		v.SetBool(b)
		return v, nil
	case durType:
		d, err := time.ParseDuration(s)
		if err != nil {
			return reflect.Value{}, err
		}
		v := reflect.New(typ).Elem()
		v.SetInt(int64(d))
		return v, nil
	default:
		// Attempt JSON decode of a single element into the target type
		dst := reflect.New(typ).Interface()
		if err := json.Unmarshal([]byte(s), dst); err != nil {
			return reflect.Value{}, fmt.Errorf("unsupported slice elem type %s", typ)
		}
		return reflect.ValueOf(dst).Elem(), nil
	}
}

func splitCSV(s string) []string {
	// simple split; if you need quoted CSV, swap for a CSV reader
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func getEnvName(input string) string {
	if len(input) <= 1 {
		return input
	}
	var result strings.Builder
	result.WriteByte(input[0])
	for i := 1; i < len(input); i++ {
		c := input[i]
		p := input[i-1]
		if unicode.IsUpper(rune(c)) && unicode.IsLower(rune(p)) {
			result.WriteByte('_')
		}
		result.WriteByte(c)
	}
	return strings.ToUpper(result.String())
}
