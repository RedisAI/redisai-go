package redisai

import (
	"github.com/gomodule/redigo/redis"
)

// ScriptInterface is an interface that represents the skeleton of a script
type ScriptInterface interface {
	Device() string
	SetDevice(device string)
	Tag() string
	SetTag(tag string)
	Source() string
	SetSource(source string)
	EntryPoints() []string
	SetEntryPoints(entryPoinys []string)
}

func scriptGetParseToInterface(reply interface{}, model ScriptInterface) (err error) {
	var device string
	var tag string
	var source string
	var entryPoints []string
	device, tag, source, entryPoints, err = scriptGetParseReply(reply)
	if err != nil {
		return err
	}
	model.SetDevice(device)
	model.SetTag(tag)
	model.SetSource(source)
	model.SetEntryPoints(entryPoints)
	return
}

func scriptGetParseReply(reply interface{}) (device string, tag string, source string, entryPoints []string, err error) {
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
		case "device":
			device, err = redis.String(replySlice[pos+1], err)
			if err != nil {
				return
			}
		case "tag":
			tag, err = redis.String(replySlice[pos+1], err)
			if err != nil {
				return
			}
		case "source":
			source, err = redis.String(replySlice[pos+1], err)
			if err != nil {
				return
			}
		case "Entry Points":
			entryPoints, err = redis.Strings(replySlice[pos+1], err)
			if err != nil {
				return
			}
		}
	}
	return
}

func scriptGetFlatArgs(name string) redis.Args {
	args := redis.Args{}.Add(name, "META", "SOURCE")
	return args
}

func scriptSetInterfaceArgs(keyName string, scriptInterface ScriptInterface) redis.Args {
	args := redis.Args{keyName}
	if len(scriptInterface.Device()) > 0 {
		args = args.Add(scriptInterface.Device())
	}
	if len(scriptInterface.Tag()) > 0 {
		args = args.Add("TAG", scriptInterface.Tag())
	}
	if len(scriptInterface.Source()) > 0 {
		args = args.Add("SOURCE", scriptInterface.Source())
	}
	return args
}

func scriptRunFlatArgs(name string, fn string, inputs []string, outputs []string) redis.Args {
	args := redis.Args{}.Add(name, fn)
	if len(inputs) > 0 {
		args = args.Add("INPUTS").AddFlat(inputs)
	}
	if len(outputs) > 0 {
		args = args.Add("OUTPUTS").AddFlat(outputs)
	}
	return args
}
