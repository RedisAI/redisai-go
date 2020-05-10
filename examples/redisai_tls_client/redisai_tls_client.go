package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/RedisAI/redisai-go/redisai"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"os"
)

var (
	tlsCertFile   = flag.String("tls-cert-file", "redis.crt", "A a X.509 certificate to use for authenticating the  server to connected clients, masters or cluster peers. The file should be PEM formatted.")
	tlsKeyFile    = flag.String("tls-key-file", "redis.key", "A a X.509 privat ekey to use for authenticating the  server to connected clients, masters or cluster peers. The file should be PEM formatted.")
	tlsCaCertFile = flag.String("tls-ca-cert-file", "ca.crt", "A PEM encoded CA's certificate file.")
	host          = flag.String("host", "127.0.0.1:6379", "Redis host.")
	password      = flag.String("password", "", "Redis password.")
)

func exists(filename string) (exists bool) {
	exists = false
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || info.IsDir() {
		return
	}
	exists = true
	return
}

/*
 * Example of how to establish an SSL connection from your app to the RedisAI Server
 */
func main() {
	flag.Parse()
	// Quickly check if the files exist
	if !exists(*tlsCertFile) || !exists(*tlsKeyFile) || !exists(*tlsCaCertFile) {
		fmt.Println("Some of the required files does not exist. Leaving example...")
		return
	}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(*tlsCertFile, *tlsKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(*tlsCaCertFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	// If InsecureSkipVerify is true, TLS accepts any certificate
	// presented by the server and any host name in that certificate.
	// In this mode, TLS is susceptible to man-in-the-middle attacks.
	// This should be used only for testing.
	clientTLSConfig.InsecureSkipVerify = true

	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", *host,
			redis.DialPassword(*password),
			redis.DialTLSConfig(clientTLSConfig),
			redis.DialUseTLS(true),
			redis.DialTLSSkipVerify(true),
		)
	}}

	// create a connection from Pool
	client := redisai.Connect("", pool)

	// Set a tensor
	// AI.TENSORSET foo FLOAT 2 2 VALUES 1.1 2.2 3.3 4.4
	_ = client.TensorSet("foo", redisai.TypeFloat, []int{2, 2}, []float32{1.1, 2.2, 3.3, 4.4})

	// Get a tensor content as a slice of values
	// dt DataType, shape []int, data interface{}, err error
	// AI.TENSORGET foo VALUES
	_, _, fooTensorValues, err := client.TensorGetValues("foo")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fooTensorValues)
	//Output: [1.1 2.2 3.3 4.4]
}
