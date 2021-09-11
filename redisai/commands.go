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
	data[0], data[1], data[2], err = ProcessTensorGetReply(reply, err)
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
	dt, shape, data, err = ProcessTensorGetReply(reply, err)
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
	dt, shape, _, err = ProcessTensorGetReply(reply, err)
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
	dt, shape, dataInterface, err := ProcessTensorGetReply(reply, err)
	data = dataInterface.([]byte)
	return
}

// ModelSet sets a RedisAI model from a blob
func (c *Client) ModelSet(keyName, backend, device string, data []byte, inputs, outputs []string) (err error) {
	args := modelStoreFlatArgs(keyName, backend, device, "", 0, 0, 0, inputs, outputs, data)
	_, err = c.DoOrSend("AI.MODELSTORE", args, nil)
	return
}

// ModelSet sets a RedisAI model from a structure that implements the ModelInterface
func (c *Client) ModelSetFromModel(keyName string, model ModelInterface) (err error) {
	args := modelStoreInterfaceArgs(keyName, model)
	_, err = c.DoOrSend("AI.MODELSTORE", args, nil)
	return
}

// ModelStore sets a RedisAI model from a blob
func (c *Client) ModelStore(keyName, backend, device, tag string, batchsize, minbatchsize, minbatchtimeout int64, inputs, outputs []string, data []byte) (err error) {
	args := modelStoreFlatArgs(keyName, backend, device, tag, batchsize, minbatchsize, minbatchtimeout, inputs, outputs, data)
	_, err = c.DoOrSend("AI.MODELSTORE", args, nil)
	return
}

// ModelStoreFromModel sets a RedisAI model from a structure that implements the ModelInterface
func (c *Client) ModelStoreFromModel(keyName string, model ModelInterface) (err error) {
	args := modelStoreInterfaceArgs(keyName, model)
	_, err = c.DoOrSend("AI.MODELSTORE", args, nil)
	return
}

// ModelGet gets a RedisAI model from the RedisAI server
// The reply will an array, containing at
//    - position 0 the backend used by the model as a String
//    - position 1 the device used to execute the model as a String
//    - position 2 the model's tag as a String
//    - position 3 a blob containing the serialized model (when called with the BLOB argument) as a String
//    - position 4 the maximum size of any batch of incoming requests.
//    - position 5 the minimum size of any batch of incoming requests.
//    - position 6 array reply with one or more names of the model's input nodes (applicable only for TensorFlow models).
//    - position 7 array reply with one or more names of the model's output nodes (applicable only for TensorFlow models).
func (c *Client) ModelGet(keyName string) (data []interface{}, err error) {
	var reply interface{}
	data = make([]interface{}, 8)
	args := modelGetFlatArgs(keyName)
	reply, err = c.DoOrSend("AI.MODELGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7], err = modelGetParseReply(reply)
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
func (c *Client) ModelRun(name string, inputs, outputs []string) (err error) {
	args := modelExecuteFlatArgs(name, inputs, outputs, 0)
	_, err = c.DoOrSend("AI.MODELEXECUTE", args, nil)
	return
}

// ModelExecute runs the model present in the keyName, with the input tensor names, and output tensor names
func (c *Client) ModelExecute(name string, inputs, outputs []string) (err error) {
	args := modelExecuteFlatArgs(name, inputs, outputs, 0)
	_, err = c.DoOrSend("AI.MODELEXECUTE", args, nil)
	return
}

// ModelExecuteWithTimeout runs the model present in the keyName, with the input tensor names, output tensor names and timeout
func (c *Client) ModelExecuteWithTimeout(name string, inputs, outputs []string, timeout int64) (err error) {
	args := modelExecuteFlatArgs(name, inputs, outputs, timeout)
	_, err = c.DoOrSend("AI.MODELEXECUTE", args, nil)
	return
}

// ScriptSet sets a RedisAI script from a blob
func (c *Client) ScriptSet(name, device, scriptSource string) (err error) {
	args := scriptStoreFlatArgs(name, device, "", nil, scriptSource)
	_, err = c.DoOrSend("AI.SCRIPTSET", args, nil)
	return
}

// ScriptSetWithTag sets a RedisAI script from a blob with tag
func (c *Client) ScriptSetWithTag(name, device, scriptSource, tag string) (err error) {
	args := scriptStoreFlatArgs(name, device, tag, nil, scriptSource)
	_, err = c.DoOrSend("AI.SCRIPTSET", args, nil)
	return
}

// ScriptSetFromInteface sets a RedisAI script from a structure that implements the ScriptInterface
func (c *Client) ScriptSetFromInteface(keyName string, script ScriptInterface) (err error) {
	args := scriptStoreInterfaceArgs(keyName, script)
	_, err = c.DoOrSend("AI.SCRIPTSET", args, nil)
	return
}

// ScriptStore store a TorchScript as the value of a key.
func (c *Client) ScriptStore(name, device, scriptSource string, entryPoints []string) (err error) {
	args := scriptStoreFlatArgs(name, device, "", entryPoints, scriptSource)
	_, err = c.DoOrSend("AI.SCRIPTSTORE", args, nil)
	return
}

// ScriptStoreWithTag store a TorchScript as the value of a key with tag.
func (c *Client) ScriptStoreWithTag(name, device, scriptSource string, entryPoints []string, tag string) (err error) {
	args := scriptStoreFlatArgs(name, device, tag, entryPoints, scriptSource)
	_, err = c.DoOrSend("AI.SCRIPTSTORE", args, nil)
	return
}

// ScriptStoreFromInteface store a TorchScript as the value from a structure that implements the ScriptInterface
func (c *Client) ScriptStoreFromInteface(keyName string, script ScriptInterface) (err error) {
	args := scriptStoreInterfaceArgs(keyName, script)
	_, err = c.DoOrSend("AI.SCRIPTSTORE", args, nil)
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
func (c *Client) ScriptRun(name, fn string, inputs, outputs []string) (err error) {
	args := scriptRunFlatArgs(name, fn, inputs, outputs)
	_, err = c.DoOrSend("AI.SCRIPTRUN", args, nil)
	return
}

// ScriptExecute run an already set script
func (c *Client) ScriptExecute(name, fn string, keys, inputs, inputArgs, outputs []string) (err error) {
	args := scriptExecuteFlatArgs(name, fn, keys, inputs, inputArgs, outputs, 0)
	_, err = c.DoOrSend("AI.SCRIPTEXECUTE", args, nil)
	return
}

// ScriptExecuteWithTimeout run an already set script with timeout limitation
func (c *Client) ScriptExecuteWithTimeout(name, fn string, keys, inputs, inputArgs, outputs []string, timeout int64) (err error) {
	args := scriptExecuteFlatArgs(name, fn, keys, inputs, inputArgs, outputs, timeout)
	_, err = c.DoOrSend("AI.SCRIPTEXECUTE", args, nil)
	return
}

func (c *Client) LoadBackend(backend_identifier, location string) (err error) {
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
func (c *Client) DagRun(loadKeys, persistKeys []string, dagCommandInterface DagCommandInterface) ([]interface{}, error) {
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

// DagExecute Direct acyclic graph of operations to run within RedisAI
func (c *Client) DagExecute(loadKeys, persistKeys []string, routing string, timeout int64, dagCommandInterface DagCommandInterface) ([]interface{}, error) {
	commandArgs, err := dagCommandInterface.FlatArgs()
	if err != nil {
		return nil, err
	}
	args := AddDagExecuteArgs(loadKeys, persistKeys, routing, timeout, commandArgs)
	reply, err := c.DoOrSend("AI.DAGEXECUTE", args, nil)
	return dagCommandInterface.ParseReply(reply, err)
}

// DagExecuteRO is the read-only variant of DagExecute
func (c *Client) DagExecuteRO(loadKeys []string, routing string, timeout int64, dagCommandInterface DagCommandInterface) ([]interface{}, error) {
	commandArgs, err := dagCommandInterface.FlatArgs()
	if err != nil {
		return nil, err
	}
	args := AddDagExecuteArgs(loadKeys, nil, routing, timeout, commandArgs)
	reply, err := c.DoOrSend("AI.DAGEXECUTE_RO", args, nil)
	return dagCommandInterface.ParseReply(reply, err)
}

// AddDagRunArgs for AI.DAGRUN and DAGRUN_RO commands.
func AddDagRunArgs(loadKeys, persistKeys []string, commandArgs redis.Args) redis.Args {
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

// AddDagExecuteArgs for AI.DAGEXECUTE and AI.DAGEXECUTE_RO commands.
func AddDagExecuteArgs(loadKeys, persistKeys []string, routing string, timeout int64, commandArgs redis.Args) redis.Args {
	args := redis.Args{}
	if loadKeys != nil {
		args = args.Add("LOAD", len(loadKeys)).AddFlat(loadKeys)
	}
	if persistKeys != nil {
		args = args.Add("PERSIST", len(persistKeys)).AddFlat(persistKeys)
	}
	if len(routing) > 0 {
		args = args.Add("ROUTING", routing)
	}
	if timeout > 0 {
		args = args.Add("TIMEOUT", timeout)
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
