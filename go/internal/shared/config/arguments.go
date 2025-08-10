package config

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func LoadFlags(cfg any) error {
	rv := reflect.ValueOf(cfg)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be pointer to struct")
	}
	rv = rv.Elem()
	rt := rv.Type()

	type holder struct {
		idx        int
		name       string
		kind       reflect.Kind
		sval       *string
		ival       *int
		bval       *bool
		hasDefault bool
	}

	var hs []holder

	for i := range rt.NumField() {
		f := rt.Field(i)
		if !rv.Field(i).CanSet() {
			continue
		}
		name := strings.ToLower(f.Name)
		def := f.Tag.Get("env_default")

		h := holder{
			idx:        i,
			name:       name,
			kind:       f.Type.Kind(),
			hasDefault: def != "",
		}

		switch f.Type.Kind() {
		case reflect.String:
			v := def
			flag.StringVar(&v, name, v, f.Name)
			h.sval = &v
		case reflect.Int:
			var d int
			if def != "" {
				d, _ = strconv.Atoi(def)
			}
			flag.IntVar(&d, name, d, f.Name)
			h.ival = &d
		case reflect.Bool:
			var b bool
			if def != "" {
				b, _ = strconv.ParseBool(def)
			}
			flag.BoolVar(&b, name, b, f.Name)
			h.bval = &b
		default:
			return fmt.Errorf("unsupported kind %s for field %s", f.Type.Kind(), f.Name)
		}

		hs = append(hs, h)
	}

	flag.Parse()

	// Which flags were actually set by the user?
	setFlags := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { setFlags[f.Name] = true })

	// Validate + assign back to struct
	for _, h := range hs {
		if !setFlags[h.name] && !h.hasDefault {
			return fmt.Errorf("missing required flag: -%s", h.name)
		}
		switch h.kind {
		case reflect.String:
			rv.Field(h.idx).SetString(*h.sval)
		case reflect.Int:
			rv.Field(h.idx).SetInt(int64(*h.ival))
		case reflect.Bool:
			rv.Field(h.idx).SetBool(*h.bval)
		}
	}

	return nil
}
