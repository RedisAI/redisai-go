package redisai_test

import (
	_ "github.com/RedisAI/redisai-go/redisai"
	"os"
)

func getTLSdetails() (tlsready bool, tls_cert string, tls_key string, tls_cacert string) {
	tlsready = false
	value, exists := os.LookupEnv("TLS_CERT")
	if exists && value != "" {
		tls_cert = value
	} else {
		return
	}
	value, exists = os.LookupEnv("TLS_KEY")
	if exists && value != "" {
		tls_key = value
	} else {
		return
	}
	value, exists = os.LookupEnv("TLS_CACERT")
	if exists && value != "" {
		tls_cacert = value
	} else {
		return
	}
	tlsready = true
	return
}