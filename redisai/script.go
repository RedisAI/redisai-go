package redisai

import (
	"github.com/gomodule/redigo/redis"
)

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
