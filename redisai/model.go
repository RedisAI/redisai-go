package redisai

import (
	"github.com/gomodule/redigo/redis"
)

// ModelInterface is an interface that represents the skeleton of a model
// needed to map it to a RedisAI Model with the proper operations
type ModelInterface interface {
	Outputs() []string
	SetOutputs(outputs []string)
	Inputs() []string
	SetInputs(inputs []string)
	Blob() []byte
	SetBlob(blob []byte)
	Device() string
	SetDevice(device string)
	Backend() string
	SetBackend(backend string)
}

func modelSetFlatArgs(keyName, backend, device string, inputs, outputs []string, blob []byte) redis.Args {
	args := redis.Args{}.Add(keyName, backend, device)
	if len(inputs) > 0 {
		args = args.Add("INPUTS").AddFlat(inputs)
	}
	if len(outputs) > 0 {
		args = args.Add("OUTPUTS").AddFlat(outputs)
	}
	args = args.Add("BLOB")
	args = args.Add(blob)
	return args
}

func modelSetInterfaceArgs(keyName string, modelInterface ModelInterface) redis.Args {
	return modelSetFlatArgs(keyName, modelInterface.Backend(), modelInterface.Device(), modelInterface.Inputs(), modelInterface.Outputs(), modelInterface.Blob())
}

func modelRunFlatArgs(name string, inputTensorNames, outputTensorNames []string) redis.Args {
	args := redis.Args{}
	args = args.Add(name)
	if len(inputTensorNames) > 0 {
		args = args.Add("INPUTS").AddFlat(inputTensorNames)
	}
	if len(outputTensorNames) > 0 {
		args = args.Add("OUTPUTS").AddFlat(outputTensorNames)
	}
	return args
}

func modelGetParseToInterface(reply interface{}, model ModelInterface) (err error) {
	var backend string
	var device string
	var blob []byte
	err, backend, device, blob = modelGetParseReply(reply)
	if err != nil {
		return err
	}
	model.SetBackend(backend)
	model.SetDevice(device)
	model.SetBlob(blob)
	return
}

func modelGetParseReply(reply interface{}) (err error, backend string, device string, blob []byte) {
	var replySlice []interface{}
	var key string
	replySlice, err = redis.Values(reply, err)
	if err != nil {
		return
	}
	for pos := 0; pos < len(replySlice); pos += 2 {
		key, err = redis.String(replySlice[pos], err)
		if err != nil {
			return
		}
		switch key {
		case "backend":
			backend, err = redis.String(replySlice[pos+1], err)
			if err != nil {
				return
			}
		case "device":
			device, err = redis.String(replySlice[pos+1], err)
			if err != nil {
				return
			}
		case "blob":
			blob, err = redis.Bytes(replySlice[pos+1], err)
			if err != nil {
				return
			}
		}
	}
	return
}

func modelGetFlatArgs(name string) redis.Args {
	args := redis.Args{}.Add(name, "META", "BLOB")
	return args
}

func modelDelFlatArgs(name string) redis.Args {
	args := redis.Args{}.Add(name)
	return args
}
