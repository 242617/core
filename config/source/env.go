package source

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Env fills struct fields from environment variables using `env` tags.
func Env() ConfigSource {
	return &env{}
}

type env struct{}

func (e *env) Scan(p interface{}) error {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("unexpected kind: %q", v.Kind())
	}
	return e.describe(v.Elem())
}

func (e *env) describe(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		tf := v.Type().Field(i)
		tag := tf.Tag.Get("env")

		if vf.Kind() == reflect.Struct {
			if err := e.describe(vf); err != nil {
				return err
			}
			continue
		}

		val := os.Getenv(tag)
		if val == "" {
			continue
		}

		switch vf.Kind() {
		case reflect.String:
			vf.SetString(val)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if vf.Kind() == reflect.Int64 && vf.Type() == reflect.TypeOf(time.Nanosecond) {
				dur, err := time.ParseDuration(val)
				if err != nil {
					return err
				}
				vf.Set(reflect.ValueOf(dur))
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
