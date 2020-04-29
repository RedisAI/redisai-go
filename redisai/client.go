package redisai

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"sync/atomic"
	"time"
)

// Client is a RedisAI client
type Client struct {
	Pool                  *redis.Pool
	PipelineActive        bool
	PipelineAutoFlushSize uint32
	PipelinePos           uint32
	ActiveConn            redis.Conn
}

// Connect intializes a Client
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

// ModelSet sets a RedisAI model from a structure that implements the ModelInterface
func (c *Client) ModelSet(keyName string, model ModelInterface) (err error) {
	args := modelSetInterfaceArgs(keyName,model)
	_, err = c.doOrSend("AI.MODELSET", args)
	return
}

// ModelSet sets a RedisAI model from a blob
func (c *Client) ModelSetFlatCmd(keyName, backend, device string, data []byte, inputs, outputs []string) (err error) {
	args := modelSetFlatArgs(keyName, backend, device, inputs, outputs, data)
	_, err = c.doOrSend("AI.MODELSET", args)
	return
}

// ModelRun runs the model present in the keyName, with the input tensor names, and output tensor names
func (c *Client) ModelRun(name string, inputTensorNames, outputTensorNames []string) (err error) {
	args := modelRunFlatArgs(name, inputTensorNames, outputTensorNames)
	_, err = c.doOrSend("AI.MODELRUN", args)
	return
}

func (c *Client) ModelGet(keyName string, modelIn ModelInterface) (err error) {
	args := modelGetFlatArgs(keyName)
	var reply interface{}
	reply, err = c.doOrSend("AI.MODELGET", args)
	err = modelGetParseToInterface(reply, modelIn)
	return
}

func (c *Client) ModelDel(keyName string) (err error) {
	args := modelDelFlatArgs(keyName)
	_, err = c.doOrSend("AI.MODELDEL", args)
	return
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

func (c *Client) TensorGet(name string, ct string) (data []interface{}, err error) {
	args := redis.Args{}.Add(name, ct)
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr("AI.TENSORGET", args)
	} else {
		resp, err := c.ActiveConn.Do("AI.TENSORGET", args...)
		data, err = ProcessTensorReplyMeta(resp, err)
		switch ct {
		case TensorContentTypeBlob:
			data, err = ProcessTensorReplyBlob(data, err)
		case TensorContentTypeValues:
			data, err = ProcessTensorReplyValues(data, err)
		default:
			err = fmt.Errorf("redisai.TensorGet: Unrecognized TensorContentType. Expected '%s' or '%s', got '%s'", TensorContentTypeBlob, TensorContentTypeValues, ct)
		}
	}
	return
}

// TensorGetValues gets a tensor's values
func (c *Client) TensorGetValues(name string) (dt string, shape []int, data interface{}, err error) {
	resp, err := c.TensorGet(name, TensorContentTypeValues)
	if err != nil {
		return
	}
	if len(resp) != 3 {
		err = fmt.Errorf("redisai.ModelGet: AI.TENSORGET returned response with incorrect sizing. expected '%d' got '%d'", 3, len(resp))
		return dt, shape, data, err
	}
	return resp[0].(string), resp[1].([]int), resp[2], err
}

// TensorGetValues gets a tensor's values
func (c *Client) TensorGetMeta(name string) (dt string, shape []int, err error) {
	resp, err := c.TensorGet(name, TensorContentTypeMeta)
	if err != nil {
		return
	}
	if len(resp) != 2 {
		err = fmt.Errorf("redisai.ModelGet: AI.TENSORGET returned response with incorrect sizing. expected '%d' got '%d'", 2, len(resp))
		return dt, shape, err
	}
	return resp[0].(string), resp[1].([]int), err
}

// TensorGetValues gets a tensor's values
func (c *Client) TensorGetBlob(name string) (dt string, shape []int, data []byte, err error) {
	resp, err := c.TensorGet(name, TensorContentTypeBlob)
	if err != nil {
		return
	}
	if len(resp) != 3 {
		err = fmt.Errorf("redisai.ModelGet: AI.TENSORGET returned response with incorrect sizing. expected '%d' got '%d'", 3, len(resp))
		return dt, shape, data, err
	}
	return resp[0].(string), resp[1].([]int), resp[2].([]byte), err
}

func (c *Client) ScriptGet(name string) (data map[string]string, err error) {
	args := redis.Args{}.Add(name, "META", "SOURCE")
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr("AI.SCRIPTGET", args)
	} else {
		respInitial, err := c.ActiveConn.Do("AI.SCRIPTGET", args...)
		if err != nil {
			return nil, err
		}
		data, err = redis.StringMap(respInitial, err)
	}
	return
}

func (c *Client) ScriptDel(name string) (err error) {
	args := redis.Args{}.Add(name)
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr("AI.SCRIPTDEL", args)
	} else {
		_, err = redis.String(c.ActiveConn.Do("AI.SCRIPTDEL", args...))
	}
	return
}

func (c *Client) LoadBackend(backend_identifier string, location string) (err error) {
	args := redis.Args{}.Add("LOADBACKEND").Add(backend_identifier).Add(location)
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr("AI.CONFIG", args)
	} else {
		_, err = redis.String(c.ActiveConn.Do("AI.CONFIG", args...))
	}
	return
}

// ScriptSet sets a RedisAI script from a blob
func (c *Client) ScriptSet(name string, device string, script_source string) (err error) {
	args := redis.Args{}.Add(name, device, "SOURCE", script_source)
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr("AI.SCRIPTSET", args)
	} else {
		_, err = redis.String(c.ActiveConn.Do("AI.SCRIPTSET", args...))
	}
	return
}

// ScriptSetFromFile sets a RedisAI script from a file
func (c *Client) ScriptSetFromFile(name string, device string, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return c.ScriptSet(name, device, string(data))
}

// ScriptRun runs a RedisAI script
func (c *Client) ScriptRun(name string, fn string, inputs []string, outputs []string) (err error) {
	args := redis.Args{}.Add(name, fn)
	if len(inputs) > 0 {
		args = args.Add("INPUTS").AddFlat(inputs)
	}
	if len(outputs) > 0 {
		args = args.Add("OUTPUTS").AddFlat(outputs)
	}
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr("AI.SCRIPTRUN", args)
	} else {
		_, err = redis.String(c.ActiveConn.Do("AI.SCRIPTRUN", args...))
	}
	return
}

// TensorSet sets a tensor
func (c *Client) TensorSet(name string, dt string, dims []int, data interface{}) (err error) {
	args, err := TensorSetArgs(name, dt, dims, data, false)
	if err != nil {
		return err
	}
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr("AI.TENSORSET", args)
	} else {
		_, err = redis.String(c.ActiveConn.Do("AI.TENSORSET", args...))
	}
	return
}

func (c *Client) doOrSend(cmdName string, args redis.Args) (reply interface{}, err error) {
	c.ActiveConnNX()
	if c.PipelineActive {
		err = c.SendAndIncr(cmdName, args)
	} else {
		reply, err = c.ActiveConn.Do(cmdName, args...)
	}
	return reply, err
}
