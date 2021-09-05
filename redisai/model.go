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
	MinBatchTimeout() int64
	SetMinBatchTimeout(minBatchSize int64)
}

func modelStoreInterfaceArgs(keyName string, modelInterface ModelInterface) redis.Args {
	return modelStoreFlatArgs(keyName, modelInterface.Backend(), modelInterface.Device(), modelInterface.Tag(), modelInterface.BatchSize(), modelInterface.MinBatchSize(), modelInterface.MinBatchTimeout(), modelInterface.Inputs(), modelInterface.Outputs(), modelInterface.Blob())
}

func modelStoreFlatArgs(keyName, backend, device, tag string, batchsize, minbatchsize, minbatchtimeout int64, inputs, outputs []string, blob []byte) redis.Args {
	args := redis.Args{}.Add(keyName, backend, device)
	if len(tag) > 0 {
		args = args.Add("TAG", tag)
	}
	if batchsize > 0 {
		args = args.Add("BATCHSIZE", batchsize)
		if minbatchsize > 0 {
			args = args.Add("MINBATCHSIZE", minbatchsize)
			if minbatchtimeout > 0 {
				args = args.Add("MINBATCHTIMEOUT", minbatchtimeout)
			}
		}
	}
	if len(inputs) > 0 {
		args = args.Add("INPUTS").Add(len(inputs)).AddFlat(inputs)
	}
	if len(outputs) > 0 {
		args = args.Add("OUTPUTS").Add(len(outputs)).AddFlat(outputs)
	}
	args = args.Add("BLOB")
	args = args.Add(blob)
	return args
}

func modelRunFlatArgs(name string, inputTensorNames, outputTensorNames []string) redis.Args {
	args := redis.Args{name}
	if len(inputTensorNames) > 0 {
		args = args.Add("INPUTS").AddFlat(inputTensorNames)
	}
	if len(outputTensorNames) > 0 {
		args = args.Add("OUTPUTS").AddFlat(outputTensorNames)
	}
	return args
}

func modelExecuteFlatArgs(name string, inputTensorNames, outputTensorNames []string, timeout int64) redis.Args {
	args := redis.Args{name}
	if len(inputTensorNames) > 0 {
		args = args.Add("INPUTS").Add(len(inputTensorNames)).AddFlat(inputTensorNames)
	}
	if len(outputTensorNames) > 0 {
		args = args.Add("OUTPUTS").Add(len(outputTensorNames)).AddFlat(outputTensorNames)
	}
	if timeout > 0 {
		args = args.Add("TIMEOUT").Add(timeout)
	}
	return args
}

func modelGetParseToInterface(reply interface{}, model ModelInterface) (err error) {
	backend, device, tag, blob, batchsize, minbatchsize, inputs, outputs, err := modelGetParseReply(reply)
	if err != nil {
		return err
	}
	model.SetBackend(backend)
	model.SetDevice(device)
	model.SetTag(tag)
	model.SetBlob(blob)
	model.SetBatchSize(batchsize)
	model.SetMinBatchSize(minbatchsize)
	model.SetInputs(inputs)
	model.SetOutputs(outputs)
	return
}

func modelGetParseReply(reply interface{}) (backend, device, tag string, blob []byte, batchsize, minbatchsize int64, inputs, outputs []string, err error) {
	var replySlice []interface{}
	var key string
	inputs = nil
	outputs = nil
	replySlice, err = redis.Values(reply, err)
	if err != nil {
		return
	}
	for pos := 0; pos < len(replySlice); pos += 2 {
		// we need this condition for after parsing err check
		if err != nil {
			break
		}
		key, err = redis.String(replySlice[pos], err)
		if err != nil {
			break
		}
		switch key {
		case "backend":
			backend, err = redis.String(replySlice[pos+1], err)
		case "device":
			device, err = redis.String(replySlice[pos+1], err)
		case "blob":
			blob, err = redis.Bytes(replySlice[pos+1], err)
		case "tag":
			tag, err = redis.String(replySlice[pos+1], err)
		case "batchsize":
			batchsize, err = redis.Int64(replySlice[pos+1], err)
		case "minbatchsize":
			minbatchsize, err = redis.Int64(replySlice[pos+1], err)
		case "inputs":
			// we need to create a temporary slice given redis.Strings creates by default a slice with capacity of the input slice even if it can't be parsed
			// so the solution is to only use the replied slice of redis.Strings in case of success. Otherwise you can have inputs filled with []string(nil)
			var temporaryInputs []string
			temporaryInputs, err = redis.Strings(replySlice[pos+1], err)
			if err == nil {
				inputs = temporaryInputs
			}
		case "outputs":
			// we need to create a temporary slice given redis.Strings creates by default a slice with capacity of the input slice even if it can't be parsed
			// so the solution is to only use the replied slice of redis.Strings in case of success. Otherwise you can have outputs filled with []string(nil)
			var temporaryOutputs []string
			temporaryOutputs, err = redis.Strings(replySlice[pos+1], err)
			if err == nil {
				outputs = temporaryOutputs
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
