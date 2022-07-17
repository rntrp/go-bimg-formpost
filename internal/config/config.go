package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Load() {
	// https://github.com/joho/godotenv#precedence--conventions
	loadDotEnv()
	tryLoad(v.BIMG_FORMPOST_ENV_DIR, ".env."+v.BIMG_FORMPOST_ENV+".local")
	if v.BIMG_FORMPOST_ENV != "test" {
		tryLoad(v.BIMG_FORMPOST_ENV_DIR, ".env.local")
	}
	tryLoad(v.BIMG_FORMPOST_ENV_DIR, ".env."+v.BIMG_FORMPOST_ENV)
	tryLoad(v.BIMG_FORMPOST_ENV_DIR, ".env")
	loadEnv()
	v.print()
}

func tryLoad(path, file string) {
	f := filepath.Join(path, file)
	if err := godotenv.Load(f); err == nil {
		log.Println("Loaded config from " + f)
	}
}

func loadDotEnv() {
	v.BIMG_FORMPOST_ENV = os.Getenv("BIMG_FORMPOST_ENV")
	if len(v.BIMG_FORMPOST_ENV) == 0 {
		v.BIMG_FORMPOST_ENV = "development"
	}
	v.BIMG_FORMPOST_ENV_DIR = os.Getenv("BIMG_FORMPOST_ENV_DIR")
}

func loadEnv() {
	v.BIMG_FORMPOST_TCP_ADDRESS = parseString("BIMG_FORMPOST_TCP_ADDRESS", ":8080")
	v.BIMG_FORMPOST_TEMP_DIR = parseString("BIMG_FORMPOST_TEMP_DIR", os.TempDir())
	v.BIMG_FORMPOST_MAX_REQUEST_SIZE = parseInt64("BIMG_FORMPOST_MAX_REQUEST_SIZE", -1)
	v.BIMG_FORMPOST_MEMORY_BUFFER_SIZE = parseInt64("BIMG_FORMPOST_MEMORY_BUFFER_SIZE", 10<<20)
	v.BIMG_FORMPOST_ENABLE_PROMETHEUS = parseBool("BIMG_FORMPOST_ENABLE_PROMETHEUS", false)
	v.BIMG_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT = parseBool("BIMG_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT", false)
	v.BIMG_FORMPOST_SHUTDOWN_TIMEOUT_SECONDS = parseDuration("BIMG_FORMPOST_SHUTDOWN_TIMEOUT_SECONDS", 0, time.Second)
}

func parseBool(env string, def bool) bool {
	if b, err := strconv.ParseBool(os.Getenv(env)); err == nil {
		return b
	}
	return def
}

func parseString(env, def string) string {
	if s, ok := os.LookupEnv(env); ok {
		return s
	}
	return def
}

func parseInt64(env string, def int64) int64 {
	if i, err := strconv.ParseInt(os.Getenv(env), 10, 64); err == nil {
		return i
	}
	return def
}

func parseDuration(env string, def int64, unit time.Duration) time.Duration {
	return time.Duration(parseInt64(os.Getenv(env), def)) * unit
}
