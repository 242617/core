package source

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Default fills struct fields with values from `default` tags.
func Default() ConfigSource {
	return &def{}
}

type def struct{}

func (d *def) Scan(p interface{}) error {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("unexpected kind: %q", v.Kind())
	}
	return d.describe(v.Elem())
}

func (d *def) describe(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		tf := v.Type().Field(i)
		tag := tf.Tag.Get("default")

		if vf.Kind() == reflect.Struct {
			if err := d.describe(vf); err != nil {
				return err
			}
			continue
		}

		if tag == "" {
			continue
		}

		switch vf.Kind() {
		case reflect.String:
			vf.SetString(tag)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if vf.Kind() == reflect.Int64 && vf.Type() == reflect.TypeOf(time.Nanosecond) {
				dur, err := time.ParseDuration(tag)
				if err != nil {
					return err
				}
				vf.Set(reflect.ValueOf(dur))
				continue
			}
			i, err := strconv.ParseInt(tag, 10, 64)
			if err != nil {
				return err
			}
			vf.SetInt(i)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u, err := strconv.ParseUint(tag, 10, 64)
			if err != nil {
				return err
			}
			vf.SetUint(u)

		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(tag, 64)
			if err != nil {
				return err
			}
			vf.SetFloat(f)

		case reflect.Bool:
			vf.SetBool(strings.ToLower(tag) == "true")

		default:
			return fmt.Errorf("unsupported type: %q", vf.Kind())
		}
	}

	return nil
}
