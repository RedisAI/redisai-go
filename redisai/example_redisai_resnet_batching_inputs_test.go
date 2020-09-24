package redisai_test

import (
	"encoding/json"
	"github.com/RedisAI/redisai-go/redisai"
	"strconv"
)

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
)

// Example of
func Example_resnet50batching() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	var classes map[string]string

	imgPath := "../tests/data/cat.jpg"
	modelPath := "../tests/models/tensorflow/imagenet/resnet50.pb"
	scriptPath := "../tests/models/tensorflow/imagenet/data_processing_script.txt"
	jsonPath := "../tests/data/imagenet_classes.json"

	// Create a client.
	aiClient := redisai.Connect("redis://localhost:6379", nil)

	imgbuf, rect := getRGBImage(imgPath)

	model, _ := ioutil.ReadFile(modelPath)
	script, _ := ioutil.ReadFile(scriptPath)

	aiClient.ModelSet("imagenet_model", redisai.BackendTF, redisai.DeviceCPU, model, []string{"images"}, []string{"output"})
	aiClient.ScriptSet("imagenet_script", redisai.DeviceCPU, string(script))
	aiClient.TensorSet("image", redisai.TypeUint8, []int64{int64(rect.Max.X), int64(rect.Max.Y), 3}, imgbuf.Bytes())
	aiClient.ScriptRun("imagenet_script", "pre_process_3ch", []string{"image"}, []string{"temp1"})
	aiClient.ModelRun("imagenet_model", []string{"temp1"}, []string{"temp2"})
	aiClient.ScriptRun("imagenet_script", "post_process", []string{"temp2"}, []string{"out"})
	_, _, fooTensorValues, _ := aiClient.TensorGetValues("out")

	fmt.Println(fooTensorValues)
	val := fooTensorValues.([]int64)[0]
	byteValue, _ := ioutil.ReadFile(jsonPath)
	json.Unmarshal([]byte(byteValue), &classes)
	fmt.Println(classes[strconv.FormatInt(val, 10)])

	//Output:
	//[281]
	//tabby, tabby catamount
}
