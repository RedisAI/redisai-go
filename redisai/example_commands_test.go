package redisai_test

import (
	"fmt"
	"github.com/RedisAI/redisai-go/redisai"
	"github.com/RedisAI/redisai-go/redisai/implementations"
	"io/ioutil"
)

func ExampleClient_TensorSet() {
	// Create a simple client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Set a tensor
	// AI.TENSORSET foo FLOAT 2 2 VALUES 1.1 2.2 3.3 4.4
	err := client.TensorSet("foo", redisai.TypeFloat, []int64{2, 2}, []float32{1.1, 2.2, 3.3, 4.4})

	// print the error (should be <nil>)
	fmt.Println(err)
	// Output: <nil>
}

func ExampleClient_TensorSetFromTensor() {
	// Create a simple client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Build a tensor
	tensor := implementations.NewAiTensor()
	tensor.SetShape([]int64{2, 2})
	tensor.SetData([]float32{1.1, 2.2, 3.3, 4.4})

	// Set a tensor
	// AI.TENSORSET foo FLOAT 2 2 VALUES 1.1 2.2 3.3 4.4
	err := client.TensorSetFromTensor("foo", tensor)

	// print the error (should be <nil>)
	fmt.Println(err)
	// Output: <nil>
}

func ExampleClient_TensorGet() {
	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Set a tensor
	// AI.TENSORSET foo FLOAT 2 2 VALUES 1.1 2.2 3.3 4.4
	_ = client.TensorSet("foo", redisai.TypeFloat, []int64{2, 2}, []float32{1.1, 2.2, 3.3, 4.4})

	// Get a tensor content as a slice of values
	// AI.TENSORGET foo VALUES
	fooTensorValues, err := client.TensorGet("foo", redisai.TensorContentTypeValues)

	fmt.Println(fooTensorValues, err)
	// Output: [FLOAT [2 2] [1.1 2.2 3.3 4.4]] <nil>
}

func ExampleClient_TensorGetToTensor() {
	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Set a tensor
	// AI.TENSORSET foo FLOAT 2 2 VALUES 1.1 2.2 3.3 4.4
	_ = client.TensorSet("foo", redisai.TypeFloat, []int64{2, 2}, []float32{1.1, 2.2, 3.3, 4.4})

	// Get a tensor content as a slice of values
	// AI.TENSORGET foo VALUES
	// Allocate an empty tensor
	fooTensor := implementations.NewAiTensor()
	err := client.TensorGetToTensor("foo", redisai.TensorContentTypeValues, fooTensor)

	// Print the tensor data
	fmt.Println(fooTensor.Data(), err)
	// Output: [1.1 2.2 3.3 4.4] <nil>
}

func ExampleClient_ModelSet() {
	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)
	data, _ := ioutil.ReadFile("./../tests/test_data/creditcardfraud.pb")
	err := client.ModelSet("financialNet", redisai.BackendTF, redisai.DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"})

	// Print the error, which should be <nil> in case of sucessfull modelset
	fmt.Println(err)
	// Output: <nil>
}

func ExampleClient_ModelGet() {
	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)
	data, _ := ioutil.ReadFile("./../tests/test_data/creditcardfraud.pb")
	err := client.ModelSet("financialNet", redisai.BackendTF, redisai.DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"})

	// Print the error, which should be <nil> in case of sucessfull modelset
	fmt.Println(err)

	/////////////////////////////////////////////////////////////
	// The important part of ModelGet example starts here
	reply, err := client.ModelGet("financialNet")
	backend := reply[0]
	device := reply[1]
	// print the error (should be <nil>)
	fmt.Println(err)
	fmt.Println(backend,device)

	// Output:
	// <nil>
	// <nil>
	// TF CPU
}

func ExampleClient_ModelSetFromModel() {
	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Create a model
	model := implementations.NewModel("TF", "CPU")
	model.SetInputs([]string{"transaction", "reference"})
	model.SetOutputs([]string{"output"})
	model.SetBlobFromFile("./../tests/test_data/creditcardfraud.pb")

	err := client.ModelSetFromModel("financialNet", model)

	// Print the error, which should be <nil> in case of successful modelset
	fmt.Println(err)
	// Output: <nil>
}

func ExampleClient_ModelGetToModel() {
	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// Create a model
	model := implementations.NewModel("TF", "CPU")
	model.SetInputs([]string{"transaction", "reference"})
	model.SetOutputs([]string{"output"})

	// Read the model from file
	model.SetBlobFromFile("./../tests/test_data/creditcardfraud.pb")

	// Set the model to RedisAI so that we can afterwards test the modelget
	err := client.ModelSetFromModel("financialNet", model)
	// print the error (should be <nil>)
	fmt.Println(err)

	/////////////////////////////////////////////////////////////
	// The important part of ModelGetToModel example starts here
	// Create an empty load to store the model from RedisAI
	model1 := implementations.NewEmptyModel()
	err = client.ModelGetToModel("financialNet", model1)
	// print the error (should be <nil>)
	fmt.Println(err)

	// print the backend and device info of the model
	fmt.Println(model1.Backend(), model1.Device())

	// Output:
	// <nil>
	// <nil>
	// TF CPU
}

func ExampleClient_ModelRun() {
	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// read the model from file
	data, err := ioutil.ReadFile("./../tests/test_data/graph.pb")

	// set the model to RedisAI
	err = client.ModelSet("example-model", redisai.BackendTF, redisai.DeviceCPU, data, []string{"a", "b"}, []string{"mul"})
	// print the error (should be <nil>)
	fmt.Println(err)

	// set the input tensors
	err = client.TensorSet("a", redisai.TypeFloat32, []int64{1}, []float32{1.1})
	err = client.TensorSet("b", redisai.TypeFloat32, []int64{1}, []float32{4.4})

	// run the model
	err = client.ModelRun("example-model", []string{"a", "b"}, []string{"mul"})
	// print the error (should be <nil>)
	fmt.Println(err)

	// Output:
	// <nil>
	// <nil>
}


func ExampleClient_Info() {
	// Create a client.
	client := redisai.Connect("redis://localhost:6379", nil)

	// read the model from file
	data, err := ioutil.ReadFile("./../tests/test_data/graph.pb")

	// set the model to RedisAI
	err = client.ModelSet("example-info", redisai.BackendTF, redisai.DeviceCPU, data, []string{"a", "b"}, []string{"mul"})
	// print the error (should be <nil>)
	fmt.Println(err)

	// set the input tensors
	err = client.TensorSet("a", redisai.TypeFloat32, []int64{1}, []float32{1.1})
	err = client.TensorSet("b", redisai.TypeFloat32, []int64{1}, []float32{4.4})

	// run the model
	err = client.ModelRun("example-info", []string{"a", "b"}, []string{"mul"})
	// print the error (should be <nil>)
	fmt.Println(err)

	// get the model run info
	info, err := client.Info("example-info")

	// one model runs
	fmt.Println(fmt.Sprintf("Total runs: %s", info["calls"]))

	// Output:
	// <nil>
	// <nil>
	// Total runs: 1
}