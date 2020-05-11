package main

import (
	"fmt"
	"github.com/RedisAI/redisai-go/redisai"
	"log"
)

func main() {

	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Enable pipeline of commands on the client, autoFlushing at 3 commands
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
	err, dtype, shape, _ := redisai.ProcessTensorGetReply(client.Receive())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dtype)
	fmt.Println(shape)
	// Output: FLOAT
	//         [1 1]
}
