// Package config package config
package config

import (
	"reflect"
	"trade_agent/global"
	"trade_agent/pkg/log"
)

func checkConfigIsEmpty(body interface{}) {
	t := reflect.ValueOf(body)
	if t.Kind() == reflect.Ptr || t.Kind() != reflect.Struct {
		log.Get().Panic("checkEmpty only accept non-prt struct input")
	}

	checkStructValue(t, body)
}

func checkStructValue(t reflect.Value, body interface{}) {
	for i := 0; i < t.NumField(); i++ {
		var isBool bool
		f := t.Field(i)
		switch f.Kind() {
		case reflect.Struct:
			checkConfigIsEmpty(f.Interface())
			continue
		case reflect.Slice:
			checkSliceValue(f, t.Type().Name(), reflect.TypeOf(body).Field(i).Name)
			continue

		case reflect.String:
			if f.IsZero() {
				log.Get().Panicf("%v.%v is empty", t.Type().Name(), reflect.TypeOf(body).Field(i).Name)
			}
		case reflect.Float64:
			if f.IsZero() {
				log.Get().Panicf("%v.%v is empty", t.Type().Name(), reflect.TypeOf(body).Field(i).Name)
			}
		case reflect.Int64:
			if f.IsZero() {
				log.Get().Panicf("%v.%v is empty", t.Type().Name(), reflect.TypeOf(body).Field(i).Name)
			}
		case reflect.Int:
			if f.IsZero() {
				log.Get().Panicf("%v.%v is empty", t.Type().Name(), reflect.TypeOf(body).Field(i).Name)
			}
		case reflect.Bool:
			isBool = true
			log.Get().Warnf("%v is %v", reflect.TypeOf(body).Field(i).Name, f.Bool())

		default:
			log.Get().Panicf("Missing Check Type: %v", f.Kind())
		}

		if global.Get().GetIsDevelopment() && !isBool {
			log.Get().Infof("%v.%v: %v", t.Type().Name(), reflect.TypeOf(body).Field(i).Name, f)
		}
	}
}

func checkSliceValue(t reflect.Value, structName, fieleName string) {
	if t.Len() == 0 {
		log.Get().Panicf("%v.%v is empty", structName, fieleName)
	} else if global.Get().GetIsDevelopment() {
		for i := 0; i < t.Len(); i++ {
			log.Get().Infof("%v.%v: %v", structName, fieleName, t.Index(i))
		}
	}
}
