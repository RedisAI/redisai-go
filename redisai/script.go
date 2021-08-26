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

func scriptGetParseToInterface(reply interface{}, script ScriptInterface) (err error) {
	device, tag, source, entryPoints, err := scriptGetParseReply(reply)
	if err != nil {
		return err
	}
	script.SetDevice(device)
	script.SetTag(tag)
	script.SetSource(source)
	script.SetEntryPoints(entryPoints)
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
		// we need this condition for after parsing err check
		if err != nil {
			break
		}
		key, err = redis.String(replySlice[pos], err)
		if err != nil {
			return
		}
		switch key {
		case "device":
			device, err = redis.String(replySlice[pos+1], err)
		case "tag":
			tag, err = redis.String(replySlice[pos+1], err)
		case "source":
			source, err = redis.String(replySlice[pos+1], err)
		case "Entry Points":
			// we need to create a temporary slice given redis.Strings creates by default a slice with capacity of the input slice even if it can't be parsed
			// so the solution is to only use the replied slice of redis.Strings in case of success. Otherwise you can have entryPoints filled with []string(nil)
			var temporaryEntryPoints []string
			temporaryEntryPoints, err = redis.Strings(replySlice[pos+1], err)
			if err == nil {
				entryPoints = temporaryEntryPoints
			}
		}
	}
	return
}

func scriptGetFlatArgs(name string) redis.Args {
	args := redis.Args{}.Add(name, "META", "SOURCE")
	return args
}

func scriptStoreInterfaceArgs(keyName string, scriptInterface ScriptInterface) redis.Args {
	return ScriptStoreFlatArgs(keyName, scriptInterface.Device(), scriptInterface.Tag(), scriptInterface.EntryPoints(), scriptInterface.Source())
}

func ScriptStoreFlatArgs(keyName, device, tag string, entryPoints []string, source string) redis.Args {
	args := redis.Args{keyName}
	args = args.Add(device)
	if len(tag) > 0 {
		args = args.Add("TAG", tag)
	}
	if len(entryPoints) > 0 {
		args = args.Add("ENTRY_POINTS").Add(len(entryPoints)).AddFlat(entryPoints)
	}
	args = args.Add("SOURCE", source)
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
