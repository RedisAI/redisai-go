package redisai

import "github.com/gomodule/redigo/redis"

// DagCommandInterface is an interface that represents the skeleton of DAG supported commands
// needed to map it to a RedisAI DAGRUN and DAGURN_RO commands
type DagCommandInterface interface {
	TensorSet(keyName, dt string, dims []int64, data interface{}) DagCommandInterface
	TensorGet(name, format string) DagCommandInterface
	ModelRun(name string, inputTensorNames, outputTensorNames []string) DagCommandInterface
	FlatArgs() (redis.Args, error)
	ParseReply(reply interface{}, err error) ([]interface{}, error)
}

type Dag struct {
	commands []redis.Args
}

func NewDag() *Dag {
	return &Dag{
		commands: make([]redis.Args, 0),
	}
}

func (d *Dag) TensorSet(keyName, dt string, dims []int64, data interface{}) DagCommandInterface {
	args := redis.Args{"AI.TENSORSET"}
	setFlatArgs, err := tensorSetFlatArgs(keyName, dt, dims, data)
	if err == nil {
		args = args.AddFlat(setFlatArgs)
	}
	d.commands = append(d.commands, args)
	return d
}

func (d *Dag) TensorGet(name, format string) DagCommandInterface {
	d.commands = append(d.commands, redis.Args{"AI.TENSORGET", name, format})
	return d
}

func (d *Dag) ModelRun(name string, inputTensorNames, outputTensorNames []string) DagCommandInterface {
	args := redis.Args{"AI.MODELRUN"}
	runFlatArgs := modelRunFlatArgs(name, inputTensorNames, outputTensorNames)
	args = args.AddFlat(runFlatArgs)
	d.commands = append(d.commands, args)
	return d
}

func (d *Dag) FlatArgs() (redis.Args, error) {
	args := redis.Args{}
	for _, command := range d.commands {
		args = args.Add("|>")
		args = args.AddFlat(command)
	}
	return args, nil
}

func (d *Dag) ParseReply(reply interface{}, err error) ([]interface{}, error) {
	return redis.Values(reply, err)
}
