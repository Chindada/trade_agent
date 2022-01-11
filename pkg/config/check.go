// Package config package config
package config

import (
	"reflect"
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
		f := t.Field(i)
		switch f.Kind() {
		case reflect.Struct:
			checkConfigIsEmpty(f.Interface())
		case reflect.Slice:
			checkSliceValue(f, reflect.TypeOf(body).Field(i).Name)

		case reflect.String:
			if f.IsZero() {
				log.Get().Panicf("%v is empty", reflect.TypeOf(body).Field(i).Name)
			}
		case reflect.Float64:
			if f.IsZero() {
				log.Get().Panicf("%v is empty", reflect.TypeOf(body).Field(i).Name)
			}
		case reflect.Int64:
			if f.IsZero() {
				log.Get().Panicf("%v is empty", reflect.TypeOf(body).Field(i).Name)
			}
		case reflect.Int:
			if f.IsZero() {
				log.Get().Panicf("%v is empty", reflect.TypeOf(body).Field(i).Name)
			}
		case reflect.Bool:
			log.Get().Warnf("Bool %v is %v", reflect.TypeOf(body).Field(i).Name, f.Bool())

		default:
			log.Get().Panicf("Missing check type: %v", f.Kind())
		}
	}
}

func checkSliceValue(t reflect.Value, sliceType interface{}) {
	if t.Len() == 0 {
		log.Get().Panicf("%v is empty", sliceType)
	}
}
