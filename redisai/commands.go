package redisai

import (
	"github.com/gomodule/redigo/redis"
)

// TensorSet sets a tensor
func (c *Client) TensorSet(keyName, dt string, dims []int, data interface{}) (err error) {
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
func (c *Client) TensorGetValues(name string) (dt string, shape []int, data interface{}, err error) {
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
func (c *Client) TensorGetMeta(name string) (dt string, shape []int, err error) {
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
func (c *Client) TensorGetBlob(name string) (dt string, shape []int, data []byte, err error) {
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
	args := modelSetFlatArgs(keyName, backend, device, inputs, outputs, data)
	_, err = c.DoOrSend("AI.MODELSET", args, nil)
	return
}

// ModelSet sets a RedisAI model from a structure that implements the ModelInterface
func (c *Client) ModelSetFromModel(keyName string, model ModelInterface) (err error) {
	args := modelSetInterfaceArgs(keyName, model)
	_, err = c.DoOrSend("AI.MODELSET", args, nil)
	return
}

func (c *Client) ModelGet(keyName string) (data []interface{}, err error) {
	var reply interface{}
	data = make([]interface{}, 3)
	args := modelGetFlatArgs(keyName)
	reply, err = c.DoOrSend("AI.MODELGET", args, nil)
	if err != nil || reply == nil {
		return
	}
	err, data[0], data[1], data[2] = modelGetParseReply(reply)
	return
}

func (c *Client) ModelGeToModel(keyName string, modelIn ModelInterface) (err error) {
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

func (c *Client) ScriptGet(name string) (data map[string]string, err error) {
	args := redis.Args{}.Add(name, "META", "SOURCE")
	respInitial, err := c.DoOrSend("AI.SCRIPTGET", args, nil)
	if err != nil || respInitial == nil {
		return
	}
	data, err = redis.StringMap(respInitial, err)
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
