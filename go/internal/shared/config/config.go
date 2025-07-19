package config

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/joho/godotenv"
)

func Load(cfg any, envFiles ...string) error {
	err := godotenv.Load(envFiles...)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return loadEnv(cfg, "")
}

func loadEnv(cfg any, prefix string) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			newPrefix := prefix + getEnvName(fieldType.Name+"_")
			if err := loadEnv(field.Addr().Interface(), newPrefix); err != nil {
				return err
			}
			continue
		}

		envVariableName := prefix + getEnvName(fieldType.Name)
		envValue, exists := os.LookupEnv(envVariableName)
		if !exists {
			envValue = fieldType.Tag.Get("env_default")
		}

		// Set the value based on the field type
		switch field.Type() {
		case reflect.TypeOf(""):
			field.SetString(envValue)
		case reflect.TypeOf(0):
			intValue, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		case reflect.TypeOf(0.0):
			floatValue, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return err
			}
			field.SetFloat(floatValue)
		case reflect.TypeOf(true):
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}
			field.SetBool(boolValue)
		case reflect.TypeOf(time.Duration(0)):
			durationValue, err := time.ParseDuration(envValue)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(durationValue))
		}
	}
	return nil
}

func getEnvName(input string) string {
	if len(input) <= 1 {
		return input
	}

	var result strings.Builder
	result.WriteByte(input[0])

	for i := 1; i < len(input); i++ {
		current := input[i]
		previous := input[i-1]

		// Add underscore if current char is uppercase and previous char is lowercase
		if unicode.IsUpper(rune(current)) && unicode.IsLower(rune(previous)) {
			result.WriteByte('_')
		}

		result.WriteByte(current)
	}

	return strings.ToUpper(result.String())
}
