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
