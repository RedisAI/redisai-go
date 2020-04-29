package redisai

import (
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

func createPool() *redis.Pool {
	value, exists := os.LookupEnv("REDISAI_TEST_HOST")
	host := "redis://localhost:6379"
	if exists && value != "" {
		host = value
	}
	cpool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.DialURL(host) },
	}
	return cpool
}

func createTestClient() *Client{
	value, exists := os.LookupEnv("REDISAI_TEST_HOST")
	host := "redis://localhost:6379"
	if exists && value != "" {
		host = value
	}
	return Connect(host, nil)
}

//func TestClient_LoadBackend(t *testing.T) {
//	keyTest1 := "test:LoadBackend:1:Unexistent"
//	keyTest2 := "test:LoadBackend:2:Unexistent:Pipelined"
//
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		backend_identifier BackendType
//		location           string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{keyTest1, fields{createPool(), false, 0, 0, nil}, args{BackendTF, "unexistant"}, true},
//		{keyTest2, fields{createPool(), true, 1, 0, nil}, args{BackendTF, "unexistant"}, true},
//
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			err := c.LoadBackend(tt.args.backend_identifier, tt.args.location)
//			if tt.fields.PipelineActive {
//				c.Flush()
//				_, err = c.Receive()
//			}
//			if (err != nil) != tt.wantErr {
//				t.Errorf("LoadBackend() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestConnect(t *testing.T) {
//	value, exists := os.LookupEnv("REDISAI_TEST_HOST")
//	urlTest1 := "redis://localhost:6379"
//	if exists && value != "" {
//		urlTest1 = value
//	}
//	cpool1 := &redis.Pool{
//		MaxIdle:     3,
//		IdleTimeout: 240 * time.Second,
//		Dial:        func() (redis.Conn, error) { return redis.DialURL(urlTest1) },
//	}
//
//	type args struct {
//		url             string
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	tests := []struct {
//		name        string
//		args        args
//		pool        *redis.Pool
//		comparePool bool
//	}{
//		{"test:Connect:WithPool:1", args{urlTest1, nil, false, 0, 0, nil}, cpool1, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			gotC := Connect(tt.args.url, tt.args.Pool)
//			if tt.comparePool == true && !reflect.DeepEqual(gotC.Pool, tt.pool) {
//				t.Errorf("Connect() = %v, want %v", gotC.Pool, tt.pool)
//			}
//		})
//	}
//}
//
//func TestClient_Close(t *testing.T) {
//	key1 := "test:Close:1:ActiveConnNil"
//	key2 := "test:Close:2"
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	tests := []struct {
//		name       string
//		fields     fields
//		wantErr    bool
//		createConn bool
//	}{
//		{key1, fields{createPool(), false, 0, 0, nil}, false, false},
//		{key2, fields{createPool(), false, 0, 0, nil}, false, true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			if tt.createConn == true {
//				c.ActiveConnNX()
//			}
//			if err := c.Close(); (err != nil) != tt.wantErr {
//				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestClient_Pipeline(t *testing.T) {
//	key1 := "test:Pipeline:1"
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		PipelineMaxSize uint32
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		{key1, fields{createPool(), true, 3, 0, nil}, args{3}},
//		{key1, fields{createPool(), false, 3, 0, nil}, args{3}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			if tt.fields.PipelineActive == false {
//				c.Pipeline(tt.args.PipelineMaxSize)
//				if c.PipelineActive != true {
//					t.Errorf("c.PipelineActive was incorrect, got: %t, want: %t.", c.PipelineActive, true)
//				}
//				if c.PipelineAutoFlushSize != tt.args.PipelineMaxSize {
//					t.Errorf("c.PipelineAutoFlushSize was incorrect, got: %d, want: %d.", c.PipelineAutoFlushSize, tt.args.PipelineMaxSize)
//				}
//			}
//			c.ActiveConnNX()
//			c.Flush()
//			for i := uint32(0); i < tt.fields.PipelineMaxSize; i++ {
//				var oldPos = c.PipelinePos
//				_, err := c.TensorGet("test:Pool:1", TensorContentTypeMeta)
//				if err != nil {
//					t.Errorf("while working on TestClient_Pipeline, TensorGet() returned error = %v", err)
//				}
//				if oldPos+1 != c.PipelinePos && c.PipelinePos != 0 {
//					t.Errorf("PipelinePos was incorrect, got: %d, want: %d.", c.PipelinePos, oldPos+1)
//				}
//			}
//			if 0 != c.PipelinePos {
//				t.Errorf("PipelinePos was incorrect, got: %d, want: %d.", c.PipelinePos, 0)
//			}
//			c.Close()
//		})
//	}
//}
//
//func TestClient_Receive(t *testing.T) {
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	tests := []struct {
//		name      string
//		fields    fields
//		wantReply interface{}
//		wantErr   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			gotReply, err := c.Receive()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Receive() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(gotReply, tt.wantReply) {
//				t.Errorf("Receive() gotReply = %v, want %v", gotReply, tt.wantReply)
//			}
//		})
//	}
//}
