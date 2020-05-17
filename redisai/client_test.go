package redisai

import (
	"github.com/gomodule/redigo/redis"
	"os"
	"reflect"
	"testing"
	"time"
)

func getConnectionDetails() (host string, password string) {
	value, exists := os.LookupEnv("REDISAI_TEST_HOST")
	host = "redis://127.0.0.1:6379"
	password = ""
	valuePassword, existsPassword := os.LookupEnv("REDISAI_TEST_PASSWORD")
	if exists && value != "" {
		host = value
	}
	if existsPassword && valuePassword != "" {
		password = valuePassword
	}
	return
}

func createPool() *redis.Pool {
	host, _ := getConnectionDetails()
	cpool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.DialURL(host) },
	}
	return cpool
}

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

func createTestClient() *Client {
	host, _ := getConnectionDetails()
	return Connect(host, nil)
}

func TestConnect(t *testing.T) {
	host, _ := getConnectionDetails()
	cpool1 := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.DialURL(host) },
	}

	type args struct {
		url             string
		Pool            *redis.Pool
		PipelineActive  bool
		PipelineMaxSize uint32
		PipelinePos     uint32
		ActiveConn      redis.Conn
	}
	tests := []struct {
		name        string
		args        args
		pool        *redis.Pool
		comparePool bool
	}{
		{"test:Connect:WithPool:1", args{host, nil, false, 0, 0, nil}, cpool1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotC := Connect(tt.args.url, tt.args.Pool)
			if tt.comparePool == true && !reflect.DeepEqual(gotC.Pool, tt.pool) {
				t.Errorf("Connect() = %v, want %v", gotC.Pool, tt.pool)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	key1 := "test:Close:1:ActiveConnNil"
	key2 := "test:Close:2"
	type fields struct {
		Pool            *redis.Pool
		PipelineActive  bool
		PipelineMaxSize uint32
		PipelinePos     uint32
		ActiveConn      redis.Conn
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		createConn bool
	}{
		{key1, fields{createPool(), false, 0, 0, nil}, false, false},
		{key2, fields{createPool(), false, 0, 0, nil}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Pool:                  tt.fields.Pool,
				PipelineActive:        tt.fields.PipelineActive,
				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
				PipelinePos:           tt.fields.PipelinePos,
				ActiveConn:            tt.fields.ActiveConn,
			}
			if tt.createConn == true {
				c.ActiveConnNX()
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Pipeline(t *testing.T) {
	key1 := "test:Pipeline:1"
	type fields struct {
		Pool            *redis.Pool
		PipelineActive  bool
		PipelineMaxSize uint32
		PipelinePos     uint32
		ActiveConn      redis.Conn
	}
	type args struct {
		PipelineMaxSize uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{key1, fields{createPool(), true, 3, 0, nil}, args{3}},
		{key1, fields{createPool(), false, 3, 0, nil}, args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Pool:                  tt.fields.Pool,
				PipelineActive:        tt.fields.PipelineActive,
				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
				PipelinePos:           tt.fields.PipelinePos,
				ActiveConn:            tt.fields.ActiveConn,
			}
			if tt.fields.PipelineActive == false {
				c.Pipeline(tt.args.PipelineMaxSize)
				if c.PipelineActive != true {
					t.Errorf("c.PipelineActive was incorrect, got: %t, want: %t.", c.PipelineActive, true)
				}
				if c.PipelineAutoFlushSize != tt.args.PipelineMaxSize {
					t.Errorf("c.PipelineAutoFlushSize was incorrect, got: %d, want: %d.", c.PipelineAutoFlushSize, tt.args.PipelineMaxSize)
				}
			}
			c.ActiveConnNX()
			c.Flush()
			for i := uint32(0); i < tt.fields.PipelineMaxSize; i++ {
				var oldPos = c.PipelinePos
				_, err := c.TensorGet("test:Pool:1", TensorContentTypeMeta)
				if err != nil {
					t.Errorf("while working on TestClient_Pipeline, TensorGet() returned error = %v", err)
				}
				if oldPos+1 != c.PipelinePos && c.PipelinePos != 0 {
					t.Errorf("PipelinePos was incorrect, got: %d, want: %d.", c.PipelinePos, oldPos+1)
				}
			}
			if 0 != c.PipelinePos {
				t.Errorf("PipelinePos was incorrect, got: %d, want: %d.", c.PipelinePos, 0)
			}
			c.Close()
		})
	}
}

func TestClient_Receive(t *testing.T) {
	type fields struct {
		Pool            *redis.Pool
		PipelineActive  bool
		PipelineMaxSize uint32
		PipelinePos     uint32
		ActiveConn      redis.Conn
	}
	tests := []struct {
		name      string
		fields    fields
		wantReply interface{}
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Pool:                  tt.fields.Pool,
				PipelineActive:        tt.fields.PipelineActive,
				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
				PipelinePos:           tt.fields.PipelinePos,
				ActiveConn:            tt.fields.ActiveConn,
			}
			gotReply, err := c.Receive()
			if (err != nil) != tt.wantErr {
				t.Errorf("Receive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReply, tt.wantReply) {
				t.Errorf("Receive() gotReply = %v, want %v", gotReply, tt.wantReply)
			}
		})
	}
}

func TestClient_DisablePipeline(t *testing.T) {
	// Create a client.
	client := createTestClient()

	// Enable pipeline of commands on the client, autoFlushing at 3 commands
	client.Pipeline(3)

	// Set a tensor
	// AI.TENSORSET foo FLOAT 2 2 VALUES 1.1 2.2 3.3 4.4
	err := client.TensorSet("foo1", TypeFloat, []int64{2, 2}, []float32{1.1, 2.2, 3.3, 4.4})
	if err != nil {
		t.Errorf("TensorSet() error = %v", err)
	}
	// AI.TENSORSET foo2 FLOAT 1" 1 VALUES 1.1
	err = client.TensorSet("foo2", TypeFloat, []int64{1, 1}, []float32{1.1})
	if err != nil {
		t.Errorf("TensorSet() error = %v", err)
	}
	// AI.TENSORGET foo2 META
	_, err = client.TensorGet("foo2", TensorContentTypeMeta)
	if err != nil {
		t.Errorf("TensorGet() error = %v", err)
	}
	// Ignore the AI.TENSORSET Reply
	_, err = client.Receive()
	if err != nil {
		t.Errorf("Receive() error = %v", err)
	}
	// Ignore the AI.TENSORSET Reply
	_, err = client.Receive()
	if err != nil {
		t.Errorf("Receive() error = %v", err)
	}
	err, _, _, _ = ProcessTensorGetReply(client.Receive())
	if err != nil {
		t.Errorf("ProcessTensorGetReply() error = %v", err)
	}

	err = client.DisablePipeline()
	if err != nil {
		t.Errorf("DisablePipeline() error = %v", err)
	}

}
