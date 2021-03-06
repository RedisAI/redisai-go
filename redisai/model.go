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
	Tag() string
	SetTag(tag string)
	BatchSize() int64
	SetBatchSize(batchSize int64)
	MinBatchSize() int64
	SetMinBatchSize(minBatchSize int64)
}

func modelSetFlatArgs(keyName, backend, device, tag string, inputs, outputs []string, blob []byte) redis.Args {
	args := redis.Args{}.Add(keyName, backend, device)
	if len(tag) > 0 {
		args = args.Add("TAG", tag)
	}
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
	if modelInterface.BatchSize() > 0 {
		args = args.Add("BATCHSIZE", modelInterface.BatchSize())
		if modelInterface.MinBatchSize() > 0 {
			args = args.Add("MINBATCHSIZE", modelInterface.MinBatchSize())
		}
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
	var tag string
	var blob []byte
	err, backend, device, tag, blob = modelGetParseReply(reply)
	if err != nil {
		return err
	}
	model.SetBackend(backend)
	model.SetDevice(device)
	model.SetTag(tag)
	model.SetBlob(blob)
	return
}

func modelGetParseReply(reply interface{}) (err error, backend string, device string, tag string, blob []byte) {
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
		case "tag":
			tag, err = redis.String(replySlice[pos+1], err)
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
