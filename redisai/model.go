package redisai

import (
	"github.com/RedisAI/redisai-go/redisai/implementations"
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
	Tag() string
	SetTag(tag string)
}

func modelSetFlatArgs(keyName, backend, device string, inputs, outputs []string, blob []byte) redis.Args {
	modelInterface := implementations.NewModel(backend, device)
	modelInterface.SetInputs(inputs)
	modelInterface.SetOutputs(outputs)
	modelInterface.SetBlob(blob)
	return modelSetInterfaceArgs(keyName, modelInterface)
}

func modelSetInterfaceArgs(keyName string, modelInterface ModelInterface) redis.Args {
	args := redis.Args{keyName}
	if len(modelInterface.Backend()) > 0 {
		args = args.Add(modelInterface.Backend())
	}
	if len(modelInterface.Device()) > 0 {
		args = args.Add(modelInterface.Device())
	}
	if len(modelInterface.Tag()) > 0 {
		args = args.Add("TAG", modelInterface.Tag())
	}
	if len(modelInterface.Inputs()) > 0 {
		args = args.Add("INPUTS").AddFlat(modelInterface.Inputs())
	}
	if len(modelInterface.Outputs()) > 0 {
		args = args.Add("OUTPUTS").AddFlat(modelInterface.Outputs())
	}
	if modelInterface.Blob() != nil {
		args = args.Add("BLOB", modelInterface.Blob())
	}
	return args
}

func modelRunFlatArgs(name string, inputTensorNames, outputTensorNames []string) redis.Args {
	modelInterface := implementations.NewEmptyModel()
	modelInterface.SetInputs(inputTensorNames)
	modelInterface.SetOutputs(outputTensorNames)
	return modelSetInterfaceArgs(name, modelInterface)
}

func modelGetParseToInterface(reply interface{}, model ModelInterface) (err error) {
	var replySlice []interface{}
	var key string
	replySlice, err = redis.Values(reply, err)
	if err != nil {
		return
	}

	var backend string
	var device string
	var blob []byte
	var tag string
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
			model.SetBackend(backend)
		case "device":
			device, err = redis.String(replySlice[pos+1], err)
			if err != nil {
				return
			}
			model.SetDevice(device)
		case "blob":
			blob, err = redis.Bytes(replySlice[pos+1], err)
			if err != nil {
				return
			}
			model.SetBlob(blob)
		case "tag":
			tag, err = redis.String(replySlice[pos+1], err)
			if err != nil {
				return
			}
			model.SetTag(tag)
		}
	}
	return err
}

func modelGetParseReply(reply interface{}) (err error, backend string, device string, blob []byte) {
	modelInterface := implementations.NewEmptyModel()
	err = modelGetParseToInterface(reply, modelInterface)
	return err, modelInterface.Backend(), modelInterface.Device(), modelInterface.Blob()
}

func modelGetFlatArgs(name string) redis.Args {
	args := redis.Args{}.Add(name, "META", "BLOB")
	return args
}

func modelDelFlatArgs(name string) redis.Args {
	args := redis.Args{}.Add(name)
	return args
}
