package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type values struct {
	BIMG_FORMPOST_ENV                      string
	BIMG_FORMPOST_ENV_DIR                  string
	BIMG_FORMPOST_TCP_ADDRESS              string
	BIMG_FORMPOST_TEMP_DIR                 string
	BIMG_FORMPOST_MAX_REQUEST_SIZE         int64
	BIMG_FORMPOST_MEMORY_BUFFER_SIZE       int64
	BIMG_FORMPOST_ENABLE_PROMETHEUS        bool
	BIMG_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT bool
	BIMG_FORMPOST_SHUTDOWN_TIMEOUT_SECONDS time.Duration
}

var v values

func (v *values) print() {
	buf := new(strings.Builder)
	buf.WriteString("Environment has been resolved to:\n")
	val := reflect.Indirect(reflect.ValueOf(v))
	valType := val.Type()
	valNumField := val.NumField()
	for i := 0; i < valNumField; i++ {
		a := valType.Field(i).Name
		b := val.Field(i).Interface()
		buf.WriteString(fmt.Sprintf("%-40s= %v\n", a, b))
	}
	log.Print(buf.String())
}

func GetEnv() string {
	return v.BIMG_FORMPOST_ENV
}

func GetEnvDir() string {
	return v.BIMG_FORMPOST_ENV_DIR
}

func GetTCPAddress() string {
	return v.BIMG_FORMPOST_TCP_ADDRESS
}

func GetTempDir() string {
	return v.BIMG_FORMPOST_TEMP_DIR
}

func GetMaxRequestSize() int64 {
	return v.BIMG_FORMPOST_MAX_REQUEST_SIZE
}

func GetMemoryBufferSize() int64 {
	return v.BIMG_FORMPOST_MEMORY_BUFFER_SIZE
}

func IsEnablePrometheus() bool {
	return v.BIMG_FORMPOST_ENABLE_PROMETHEUS
}

func IsEnableShutdown() bool {
	return v.BIMG_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT
}

func GetShutdownTimeout() time.Duration {
	return v.BIMG_FORMPOST_SHUTDOWN_TIMEOUT_SECONDS
}
