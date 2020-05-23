package source

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Default creates config source that fills config with default values
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
			err := d.describe(vf)
			if err != nil {
				return err
			}
			continue
		}

		val := tag
		if val == "" {
			continue
		}

		switch vf.Kind() {

		case reflect.String:
			vf.SetString(val)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if vf.Kind() == reflect.Int64 && vf.Type() == reflect.TypeOf(time.Nanosecond) {
				v, err := time.ParseDuration(val)
				if err != nil {
					return err
				}
				vf.Set(reflect.ValueOf(v))
				continue
			}

			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			vf.SetInt(i)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return err
			}
			vf.SetUint(u)

		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			vf.SetFloat(f)

		case reflect.Bool:
			vf.SetBool(strings.ToLower(val) == "true")

		default:
			return fmt.Errorf("unsupported type: %q", vf.Kind())
		}

	}

	return nil
}
