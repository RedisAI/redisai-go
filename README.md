[![license](https://img.shields.io/github/license/RediSearch/redisearch-go.svg)](https://github.com/RedisAI/redisai-go)
[![CircleCI](https://circleci.com/gh/RedisAI/redisai-go/tree/master.svg?style=svg)](https://circleci.com/gh/RedisAI/redisai-go/tree/master)
[![GitHub issues](https://img.shields.io/github/release/RedisAI/redisai-go.svg)](https://github.com/RedisAI/redisai-go/releases/latest)
[![Codecov](https://codecov.io/gh/RedisAI/redisai-go/branch/master/graph/badge.svg)](https://codecov.io/gh/RedisAI/redisai-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedisAI/redisai-go)](https://goreportcard.com/report/github.com/RedisAI/redisai-go)
[![GoDoc](https://godoc.org/github.com/RedisAI/redisai-go?status.svg)](https://godoc.org/github.com/RedisAI/redisai-go)

# RedisAI Go Client
[![Forum](https://img.shields.io/badge/Forum-RedisAI-blue)](https://forum.redislabs.com/c/modules/redisai)
[![Gitter](https://badges.gitter.im/RedisLabs/RedisAI.svg)](https://gitter.im/RedisLabs/RedisAI?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

Go client for [RedisAI](http://redisai.io), based on redigo.

# Installing 

```sh
go get github.com/RedisAI/redisai-go/redisai
```

# Supported RedisAI Commands

| Command | Recommended API and godoc  |
| :---          |  ----: |
AI.TENSORSET | [TensorSet](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.TensorSet) and [TensorSetFromTensor](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.TensorSetFromTensor)
AI.TENSORGET | [TensorGet](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.TensorGet) and [TensorGetToTensor](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.TensorGetToTensor)
AI.MODELSET | [ModelSet](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ModelSet) and [ModelSetFromModel](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ModelSetFromModel)
AI.MODELGET | [ModelGet](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ModelGet) and [ModelGetToModel](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ModelGetToModel)
AI.MODELDEL | [ModelDel](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ModelDel)
AI.MODELRUN | [ModelRun](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ModelRun)
AI._MODELSCAN |  
AI.SCRIPTSET | [ScriptSet](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ScriptSet)
AI.SCRIPTGET | [ScriptGet](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ScriptGet)
AI.SCRIPTDEL | [ScriptDel](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ScriptDel)
AI.SCRIPTRUN | [ScriptRun](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.ScriptRun)
AI._SCRIPTSCAN |  
AI.DAGRUN | [DagRun](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.DagRun)
AI.DAGRUN_RO | [DagRunRO](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.DagRunRO)
AI.INFO |  [Info](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.Info)
AI.CONFIG * | [LoadBackend](https://godoc.org/github.com/RedisAI/redisai-go/redisai#Client.LoadBackend)


# Usage Examples
See the [examples](./examples) folder for further feature samples:

## Simple Client 
[(sample code here)](./examples/redisai_simple_client)

```go
package main

import (
	"fmt"
	"github.com/RedisAI/redisai-go/redisai"
	"log"
)

func main() {

	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Set a tensor
	// AI.TENSORSET foo FLOAT 2 2 VALUES 1.1 2.2 3.3 4.4
	_ = client.TensorSet("foo", redisai.TypeFloat, []int64{2, 2}, []float32{1.1, 2.2, 3.3, 4.4})

	// Get a tensor content as a slice of values
	// dt DataType, shape []int, data interface{}, err error
	// AI.TENSORGET foo VALUES
	_, _, fooTensorValues, err := client.TensorGetValues("foo")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fooTensorValues)
	// Output: [1.1 2.2 3.3 4.4]
}
```

## Pipelined Client 
[(sample code here)](./examples/redisai_pipelined_client)
```go
package main

import (
	"fmt"
	"github.com/RedisAI/redisai-go/redisai"
	"log"
)

func main() {

	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Enable pipeline of commands on the client.
	client.Pipeline(3)

	// Set a tensor
	// AI.TENSORSET foo FLOAT 2 2 VALUES 1.1 2.2 3.3 4.4
	err := client.TensorSet("foo1", redisai.TypeFloat, []int64{2, 2}, []float32{1.1, 2.2, 3.3, 4.4})
	if err != nil {
		log.Fatal(err)
	}
	// AI.TENSORSET foo2 FLOAT 1" 1 VALUES 1.1
	err = client.TensorSet("foo2", redisai.TypeFloat, []int64{1, 1}, []float32{1.1})
	if err != nil {
		log.Fatal(err)
	}
	// AI.TENSORGET foo2 META
	_, err = client.TensorGet("foo2", redisai.TensorContentTypeMeta)
	if err != nil {
		log.Fatal(err)
	}
	// Ignore the AI.TENSORSET Reply
	_, err = client.Receive()
	if err != nil {
		log.Fatal(err)
	}
	// Ignore the AI.TENSORSET Reply
	_, err = client.Receive()
	if err != nil {
		log.Fatal(err)
	}
	foo2TensorMeta, err := client.Receive()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(foo2TensorMeta)
	// Output: [FLOAT [1 1]]
}
```
