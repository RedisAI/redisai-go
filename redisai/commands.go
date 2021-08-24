package redisai

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

// TensorSet sets a tensor
func (c *Client) TensorSet(keyName, dt string, dims []int64, data interface{}) (err error) {
	args, err := tensorSetFlatArgs(keyName, dt, dims, data)
	_, err = c.DoOrSend("AI.TENSORSET", args, err)
	return
}

// TensorSet sets a tensor
func (c *Client) TensorSetFromTensor(keyName string, tensor TensorInterface) (err error) {
	args, err := tensorSetInterfaceArgs(keyName, tensor)
	_, err = c.DoOrSend("AI.TENSORSET", args, err)
	return
}

func (c *Client) TensorGet(name, format string) (data []interface{}, err error) {
	args := redis.Args{}.Add(name, TensorContentTypeMeta, format)
	data = make([]interface{}, 3)
	var reply interface{}
	reply, err = c.DoOrSend("AI.TENSORGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	err, data[0], data[1], data[2] = ProcessTensorGetReply(reply, err)
	return
}

func (c *Client) TensorGetToTensor(name, format string, tensor TensorInterface) (err error) {
	args := redis.Args{}.Add(name, TensorContentTypeMeta, format)
	var reply interface{}
	reply, err = c.DoOrSend("AI.TENSORGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	err = tensorGetParseToInterface(reply, tensor)
	return
}

// TensorGetValues gets a tensor's values
func (c *Client) TensorGetValues(name string) (dt string, shape []int64, data interface{}, err error) {
	args := redis.Args{}.Add(name, TensorContentTypeMeta, TensorContentTypeValues)
	var reply interface{}
	reply, err = c.DoOrSend("AI.TENSORGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	err, dt, shape, data = ProcessTensorGetReply(reply, err)
	return
}

// TensorGetValues gets a tensor's values
func (c *Client) TensorGetMeta(name string) (dt string, shape []int64, err error) {
	args := redis.Args{}.Add(name, TensorContentTypeMeta)
	var reply interface{}
	reply, err = c.DoOrSend("AI.TENSORGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	err, dt, shape, _ = ProcessTensorGetReply(reply, err)
	return
}

// TensorGetValues gets a tensor's values
func (c *Client) TensorGetBlob(name string) (dt string, shape []int64, data []byte, err error) {
	args := redis.Args{}.Add(name, TensorContentTypeMeta, TensorContentTypeBlob)
	var reply interface{}
	reply, err = c.DoOrSend("AI.TENSORGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	err, dt, shape, dataInterface := ProcessTensorGetReply(reply, err)
	data = dataInterface.([]byte)
	return
}

// ModelSet sets a RedisAI model from a blob
func (c *Client) ModelSet(keyName, backend, device string, data []byte, inputs, outputs []string) (err error) {
	args := modelSetFlatArgs(keyName, backend, device, "", inputs, outputs, data)
	_, err = c.DoOrSend("AI.MODELSET", args, nil)
	return
}

// ModelSet sets a RedisAI model from a structure that implements the ModelInterface
func (c *Client) ModelSetFromModel(keyName string, model ModelInterface) (err error) {
	args := modelSetInterfaceArgs(keyName, model)
	_, err = c.DoOrSend("AI.MODELSET", args, nil)
	return
}

// ModelGet gets a RedisAI model from the RedisAI server
// The reply will an array, containing at
//    - position 0 the backend used by the model as a String
//    - position 1 the device used to execute the model as a String
//    - position 2 the model's tag as a String
//    - position 3 a blob containing the serialized model (when called with the BLOB argument) as a String
func (c *Client) ModelGet(keyName string) (data []interface{}, err error) {
	var reply interface{}
	data = make([]interface{}, 4)
	args := modelGetFlatArgs(keyName)
	reply, err = c.DoOrSend("AI.MODELGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	data[0], data[1], data[2], data[3], err = modelGetParseReply(reply)
	return
}

// ModelGetToModel gets a RedisAI model from the RedisAI server as ModelInterface
func (c *Client) ModelGetToModel(keyName string, modelIn ModelInterface) (err error) {
	args := modelGetFlatArgs(keyName)
	var reply interface{}
	reply, err = c.DoOrSend("AI.MODELGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	err = modelGetParseToInterface(reply, modelIn)
	return
}

func (c *Client) ModelDel(keyName string) (err error) {
	args := modelDelFlatArgs(keyName)
	_, err = c.DoOrSend("AI.MODELDEL", args, nil)
	return
}

// ModelRun runs the model present in the keyName, with the input tensor names, and output tensor names
func (c *Client) ModelRun(name string, inputTensorNames, outputTensorNames []string) (err error) {
	args := modelRunFlatArgs(name, inputTensorNames, outputTensorNames)
	_, err = c.DoOrSend("AI.MODELRUN", args, nil)
	return
}

// ScriptSet sets a RedisAI script from a blob
func (c *Client) ScriptSet(name string, device string, script_source string) (err error) {
	args := redis.Args{}.Add(name, device, "SOURCE", script_source)
	_, err = c.DoOrSend("AI.SCRIPTSET", args, nil)
	return
}

// ScriptSetWithTag sets a RedisAI script from a blob with tag
func (c *Client) ScriptSetWithTag(name string, device string, script_source string, tag string) (err error) {
	args := redis.Args{}.Add(name, device, "TAG", tag, "SOURCE", script_source)
	_, err = c.DoOrSend("AI.SCRIPTSET", args, nil)
	return
}

// ScriptSetFromInteface sets a RedisAI script from a structure that implements the ScriptInterface
func (c *Client) ScriptSetFromInteface(keyName string, script ScriptInterface) (err error) {
	args := scriptSetInterfaceArgs(keyName, script)
	_, err = c.DoOrSend("AI.SCRIPTSET", args, nil)
	return
}

// ScriptGet gets a RedisAI script from the RedisAI server
// The reply will an array, containing at
//    - position 0 the script's device as a String
//    - position 1 the scripts's tag as a String
//    - position 2 the script's source code as a String
//    - position 3 an array containing the script entry point functions
func (c *Client) ScriptGet(name string) (data []interface{}, err error) {
	var reply interface{}
	data = make([]interface{}, 4)
	args := scriptGetFlatArgs(name)
	reply, err = c.DoOrSend("AI.SCRIPTGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	data[0], data[1], data[2], data[3], err = scriptGetParseReply(reply)
	return
}

// ScriptGetToInterface gets a RedisAI script from the RedisAI server as ScriptInterface
func (c *Client) ScriptGetToInterface(name string, scriptIn ScriptInterface) (err error) {
	args := scriptGetFlatArgs(name)
	reply, err := c.DoOrSend("AI.SCRIPTGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	err = scriptGetParseToInterface(reply, scriptIn)
	return
}

func (c *Client) ScriptDel(name string) (err error) {
	args := redis.Args{}.Add(name)
	_, err = c.DoOrSend("AI.SCRIPTDEL", args, nil)
	return
}

// ScriptRun runs a RedisAI script
func (c *Client) ScriptRun(name string, fn string, inputs []string, outputs []string) (err error) {
	args := scriptRunFlatArgs(name, fn, inputs, outputs)
	_, err = c.DoOrSend("AI.SCRIPTRUN", args, nil)
	return
}

func (c *Client) LoadBackend(backend_identifier string, location string) (err error) {
	args := redis.Args{}.Add("LOADBACKEND").Add(backend_identifier).Add(location)
	_, err = c.DoOrSend("AI.CONFIG", args, nil)
	return
}

// Returns information about the execution a model or a script.
func (c *Client) Info(key string) (map[string]string, error) {
	reply, err := c.DoOrSend("AI.INFO", redis.Args{key}, nil)
	values, err := redis.Values(reply, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("expects even number of values result")
	}

	m := make(map[string]string, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		k := string(values[i].([]byte))
		switch v := values[i+1].(type) {
		case []byte:
			m[k] = string(values[i+1].([]byte))
		case int64:
			m[k] = strconv.FormatInt(values[i+1].(int64), 10)
		default:
			return nil, fmt.Errorf("unexpected element type for (Ints,String), got type %T", v)
		}
	}
	return m, nil
}

// Resets all statistics associated with the key
func (c *Client) ResetStat(key string) (string, error) {
	return redis.String(c.DoOrSend("AI.INFO", redis.Args{key, "RESETSTAT"}, nil))
}

// Direct acyclic graph of operations to run within RedisAI
func (c *Client) DagRun(loadKeys []string, persistKeys []string, dagCommandInterface DagCommandInterface) ([]interface{}, error) {
	commandArgs, err := dagCommandInterface.FlatArgs()
	if err != nil {
		return nil, err
	}
	args := AddDagRunArgs(loadKeys, persistKeys, commandArgs)
	reply, err := c.DoOrSend("AI.DAGRUN", args, nil)
	return dagCommandInterface.ParseReply(reply, err)
}

// The command is a read-only variant of AI.DAGRUN
func (c *Client) DagRunRO(loadKeys []string, dagCommandInterface DagCommandInterface) ([]interface{}, error) {
	commandArgs, err := dagCommandInterface.FlatArgs()
	if err != nil {
		return nil, err
	}
	args := AddDagRunArgs(loadKeys, nil, commandArgs)
	reply, err := c.DoOrSend("AI.DAGRUN_RO", args, nil)
	return dagCommandInterface.ParseReply(reply, err)
}

// AddDagRunArgs for AI.DAGRUN and DAGRUN_RO commands.
func AddDagRunArgs(loadKeys []string, persistKeys []string, commandArgs redis.Args) redis.Args {
	args := redis.Args{}
	if loadKeys != nil {
		args = args.Add("LOAD", len(loadKeys)).AddFlat(loadKeys)
	}

	if persistKeys != nil {
		args = args.Add("PERSIST", len(persistKeys)).AddFlat(persistKeys)
	}

	if commandArgs != nil {
		args = args.AddFlat(commandArgs)
	}
	return args
}

// Sets the default backends path
func (c *Client) SetBackendsPath(path string) (string, error) {
	return redis.String(c.DoOrSend("AI.CONFIG", redis.Args{"BACKENDSPATH", path}, nil))
}
