package redisai

import (
	"github.com/gomodule/redigo/redis"
	"sync/atomic"
	"time"
)

const (
	// BackendTF represents a TensorFlow backend
	BackendTF = string("TF")
	// BackendTorch represents a Torch backend
	BackendTorch = string("TORCH")
	// BackendONNX represents an ONNX backend
	BackendONNX = string("ORT")

	// DeviceCPU represents a CPU device
	DeviceCPU = string("CPU")
	// DeviceGPU represents a GPU device
	DeviceGPU = string("GPU")

	// TypeFloat represents a float type
	TypeFloat = string("FLOAT")
	// TypeDouble represents a double type
	TypeDouble = string("DOUBLE")
	// TypeInt8 represents a int8 type
	TypeInt8 = string("INT8")
	// TypeInt16 represents a int16 type
	TypeInt16 = string("INT16")
	// TypeInt32 represents a int32 type
	TypeInt32 = string("INT32")
	// TypeInt64 represents a int64 type
	TypeInt64 = string("INT64")
	// TypeUint8 represents a uint8 type
	TypeUint8 = string("UINT8")
	// TypeUint16 represents a uint16 type
	TypeUint16 = string("UINT16")
	// TypeFloat32 is an alias for float
	TypeFloat32 = string("FLOAT")
	// TypeFloat64 is an alias for double
	TypeFloat64 = string("DOUBLE")

	// TensorContentTypeBLOB is an alias for BLOB tensor content
	TensorContentTypeBlob = string("BLOB")

	// TensorContentTypeBLOB is an alias for BLOB tensor content
	TensorContentTypeValues = string("VALUES")

	// TensorContentTypeBLOB is an alias for BLOB tensor content
	TensorContentTypeMeta = string("META")
)

type AiClient interface {
	// Close ensures that no connection is kept alive and prior to that we flush all db commands
	Close() error
	DoOrSend(string, redis.Args, error) (interface{}, error)
}

// Client is a RedisAI client
type Client struct {
	Pool                  *redis.Pool
	PipelineActive        bool
	PipelineAutoFlushSize uint32
	PipelinePos           uint32
	ActiveConn            redis.Conn
}

// Connect establish an connection to the RedisAI Server.
//
//If a pool `*redis.Pool` is passed then it will be used to connect to the server.
//
// See the examples on how to connect with/without pool and on how to establish a secure SSL connection.
func Connect(url string, pool *redis.Pool) (c *Client) {
	var cpool *redis.Pool = nil
	if pool == nil {
		cpool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.DialURL(url) },
		}
	} else {
		cpool = pool
	}

	c = &Client{
		Pool:                  cpool,
		PipelineActive:        false,
		PipelineAutoFlushSize: 0,
		PipelinePos:           0,
		ActiveConn:            nil,
	}

	return c
}

// Close ensures that no connection is kept alive and prior to that we flush all db commands
func (c *Client) Close() (err error) {
	if c.ActiveConn != nil {
		err = c.ActiveConn.Flush()
		if err != nil {
			return
		}
		err = c.ActiveConn.Close()
		if err != nil {
			return
		}
	}
	return
}

func (c *Client) ActiveConnNX() {
	if c.ActiveConn == nil {
		c.ActiveConn = c.Pool.Get()
	}
}

func (c *Client) Pipeline(PipelineAutoFlushAtSize uint32) {
	c.PipelineActive = true
	c.PipelinePos = 0
	c.PipelineAutoFlushSize = PipelineAutoFlushAtSize
}

func (c *Client) DisablePipeline() (err error) {
	err = c.Flush()
	c.PipelineActive = false
	c.PipelinePos = 0
	return
}

func (c *Client) Flush() (err error) {
	if c.ActiveConn != nil && c.PipelineActive {
		atomic.StoreUint32(&c.PipelinePos, 0)
		err = c.ActiveConn.Flush()
	}
	return
}

// Receive receives a single reply from the Redis server
func (c *Client) Receive() (reply interface{}, err error) {
	if c.ActiveConn != nil && c.PipelineActive {
		return c.ActiveConn.Receive()
	}
	return
}

func (c *Client) SendAndIncr(commandName string, args redis.Args) (err error) {
	err = c.ActiveConn.Send(commandName, args...)
	if err != nil {
		return err
	}
	// incremement the pipeline
	// flush if required
	err = c.pipeIncr(c.ActiveConn)
	return
}

func (c *Client) pipeIncr(conn redis.Conn) (err error) {
	atomic.AddUint32(&c.PipelinePos, 1)
	if c.PipelinePos >= c.PipelineAutoFlushSize && c.PipelineAutoFlushSize != 0 {
		err = conn.Flush()
		atomic.StoreUint32(&c.PipelinePos, 0)
	}
	return
}

func (c *Client) DoOrSend(cmdName string, args redis.Args, errIn error) (reply interface{}, err error) {
	err = errIn
	if err != nil {
		return
	}
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr(cmdName, args)
	} else {
		reply, err = c.ActiveConn.Do(cmdName, args...)
	}
	return reply, err
}
