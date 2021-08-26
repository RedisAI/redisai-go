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
	args := redis.Args{name}
	if len(inputTensorNames) > 0 {
		args = args.Add("INPUTS").AddFlat(inputTensorNames)
	}
	if len(outputTensorNames) > 0 {
		args = args.Add("OUTPUTS").AddFlat(outputTensorNames)
	}
	return args
}

func modelGetParseToInterface(reply interface{}, model ModelInterface) (err error) {
	err, backend, device, tag, blob, batchsize, minbatchsize, inputs, outputs := modelGetParseReply(reply)
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

func modelGetParseReply(reply interface{}) (err error, backend string, device string, tag string, blob []byte, batchsize int64, minbatchsize int64, inputs []string, outputs []string) {
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
