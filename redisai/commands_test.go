package redisai

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/RedisAI/redisai-go/redisai/implementations"
	"github.com/gomodule/redigo/redis"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var script string = `def bar(tensors: List[Tensor], keys: List[str], args: List[str]):
	a = tensors[0]
	b = tensors[1]
	return a + b

def bar_variadic(tensors: List[Tensor], keys: List[str], args: List[str]):
	a = tensors[0]
	l = tensors[1:]
	return a + l[0]
`

var scriptWithRedisCommands string = `def redis_string_int_to_tensor(redis_value: Any):
    return torch.tensor(int(str(redis_value)))

def int_set_get(tensors: List[Tensor], keys: List[str], args: List[str]):
    key = keys[0]
    value = int(args[0])
    redis.execute("SET", key, str(value))
    res = redis.execute("GET", key)
    return redis_string_int_to_tensor(res)

def func(tensors: List[Tensor], keys: List[str], args: List[str]):
    redis.execute("SET", keys[0], args[0])
    a = torch.stack(tensors).sum()
    b = redis_string_int_to_tensor(redis.execute("GET", keys[0]))
    redis.execute("DEL", keys[0])
    return b + a
`

func TestCommand_TensorSet(t *testing.T) {

	valuesFloat32 := []float32{1.1}
	valuesFloat64 := []float64{1.1}

	valuesInt8 := []int8{1}
	valuesInt16 := []int16{1}
	valuesInt32 := []int64{1}
	valuesInt64 := []int64{1}

	valuesUint8 := []uint8{1}
	valuesByte := []byte{1}

	valuesUint16 := []uint16{1}
	valuesUint32 := []uint32{1}
	valuesUint64 := []uint64{1}

	keyFloat32 := "test:TensorSet:TypeFloat32:1"
	keyFloat64 := "test:TensorSet:TypeFloat64:1"

	keyInt8 := "test:TensorSet:TypeInt8:1"
	keyInt16 := "test:TensorSet:TypeInt16:1"
	keyInt32 := "test:TensorSet:TypeInt32:1"
	keyInt64 := "test:TensorSet:TypeInt64:1"

	keyByte := "test:TensorSet:Type[]byte:1"
	keyUint8 := "test:TensorSet:TypeUint8:1"
	keyUint16 := "test:TensorSet:TypeUint16:1"
	keyUint32 := "test:TensorSet:TypeUint32:ExpectError:1"
	keyUint64 := "test:TensorSet:TypeUint64:ExpectError:1"

	keyInt8Meta := "test:TensorSet:TypeInt8:Meta:1"
	keyInt8MetaPipelined := "test:TensorSet:TypeInt8:Meta:2:Pipelined"

	keyError1 := "test:TestCommand_TensorSet:1:FaultyDims"

	shp := []int64{1}

	type args struct {
		name string
		dt   string
		dims []int64
		data interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{keyFloat32, args{keyFloat32, TypeFloat, shp, valuesFloat32}, false},
		{keyFloat64, args{keyFloat64, TypeFloat64, shp, valuesFloat64}, false},

		{keyInt8, args{keyInt8, TypeInt8, shp, valuesInt8}, false},
		{keyInt16, args{keyInt16, TypeInt16, shp, valuesInt16}, false},
		{keyInt32, args{keyInt32, TypeInt32, shp, valuesInt32}, false},
		{keyInt64, args{keyInt64, TypeInt64, shp, valuesInt64}, false},

		{keyUint8, args{keyUint8, TypeUint8, shp, valuesUint8}, false},
		{keyUint16, args{keyUint16, TypeUint16, shp, valuesUint16}, false},
		{keyUint32, args{keyUint32, TypeUint8, shp, valuesUint32}, true},
		{keyUint64, args{keyUint64, TypeUint8, shp, valuesUint64}, true},

		{keyInt8Meta, args{keyInt8Meta, TypeUint8, shp, nil}, false},
		{keyInt8MetaPipelined, args{keyInt8MetaPipelined, TypeUint8, shp, nil}, false},

		{keyByte, args{keyByte, TypeUint8, shp, valuesByte}, false},

		{keyError1, args{keyError1, TypeFloat, []int64{1, 10}, []float32{1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			if err := c.TensorSet(tt.args.name, tt.args.dt, tt.args.dims, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("TensorSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_FullFromTensor(t *testing.T) {
	tensor := implementations.NewAiTensorWithShape([]int64{1})
	tensor.SetData([]float32{1.0})
	client := createTestClient()
	err := client.TensorSetFromTensor("tensor1", tensor)
	assert.Nil(t, err)
	gotResp, err := client.TensorGet("tensor1", TensorContentTypeValues)
	assert.Nil(t, err)
	if diff := cmp.Diff(TypeFloat32, gotResp[0]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]int64{1}, gotResp[1]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]float32{1.0}, gotResp[2]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	tensor2 := implementations.NewAiTensor()
	err = client.TensorGetToTensor("tensor1", TensorContentTypeValues, tensor2)
	assert.Nil(t, err)
	assert.Equal(t, tensor, tensor2)
	// test for BLOB equality
	tensor2BlobData := implementations.NewAiTensor()
	err = client.TensorGetToTensor("tensor1", TensorContentTypeBlob, tensor2BlobData)
	assert.Nil(t, err)
	_, _, blob, err := client.TensorGetBlob("tensor1")
	assert.Nil(t, err)
	if diff := cmp.Diff(blob, tensor2BlobData.Data()); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	err = client.TensorGetToTensor("tensorDontExist", TensorContentTypeValues, tensor2)
	assert.NotNil(t, err)
}

func TestCommand_TensorGet(t *testing.T) {
	valuesByteSlice := []byte{1, 2, 3, 4, 5}

	valuesFloat32 := []float32{1.1, 2.2, 3.3, 4.4, 5.5}
	valuesFloat64 := []float64{1.1}

	valuesInt8 := []int8{1}
	valuesInt16 := []int16{1}
	valuesInt32 := []int{1}
	valuesInt64 := []int64{1}

	valuesUint8 := []uint8{1}
	valuesUint16 := []uint16{1}

	keyByteSlice := "test:TensorGet:[]byte:1"
	keyFloat32 := "test:TensorGet:TypeFloat32:1"
	keyFloat64 := "test:TensorGet:TypeFloat64:1"

	keyInt8 := "test:TensorGet:TypeInt8:1"
	keyInt16 := "test:TensorGet:TypeInt16:1"
	keyInt32 := "test:TensorGet:TypeInt32:1"
	keyInt64 := "test:TensorGet:TypeInt64:1"

	keyUint8 := "test:TensorGet:TypeUint8:1"
	keyUint16 := "test:TensorGet:TypeUint16:1"

	shp := []int64{1}
	shpByteSlice := []int64{1, 5}
	simpleClient := createTestClient()
	simpleClient.TensorSet(keyByteSlice, TypeUint8, shpByteSlice, valuesByteSlice)

	simpleClient.TensorSet(keyFloat32, TypeFloat32, shpByteSlice, valuesFloat32)
	simpleClient.TensorSet(keyFloat64, TypeFloat64, shp, valuesFloat64)

	simpleClient.TensorSet(keyInt8, TypeInt8, shp, valuesInt8)
	simpleClient.TensorSet(keyInt16, TypeInt16, shp, valuesInt16)
	simpleClient.TensorSet(keyInt32, TypeInt32, shp, valuesInt32)
	simpleClient.TensorSet(keyInt64, TypeInt64, shp, valuesInt64)

	simpleClient.TensorSet(keyUint8, TypeUint8, shp, valuesUint8)
	simpleClient.TensorSet(keyUint16, TypeUint16, shp, valuesUint16)

	type args struct {
		name string
		ct   string
	}
	tests := []struct {
		name         string
		args         args
		wantDt       string
		wantShape    []int64
		wantData     interface{}
		compareDt    bool
		compareShape bool
		compareData  bool
		wantErr      bool
	}{
		{keyByteSlice, args{keyByteSlice, TensorContentTypeBlob}, TypeUint8, shpByteSlice, valuesByteSlice, true, true, true, false},

		{keyFloat32, args{keyFloat32, TensorContentTypeValues}, TypeFloat32, shpByteSlice, valuesFloat32, true, true, true, false},
		{keyFloat64, args{keyFloat64, TensorContentTypeValues}, TypeFloat64, shp, valuesFloat64, true, true, true, false},

		{keyInt8, args{keyInt8, TensorContentTypeValues}, TypeInt8, shp, valuesInt8, true, true, true, false},
		{keyInt16, args{keyInt16, TensorContentTypeValues}, TypeInt16, shp, valuesInt16, true, true, true, false},
		{keyInt32, args{keyInt32, TensorContentTypeValues}, TypeInt32, shp, valuesInt32, true, true, true, false},
		{keyInt64, args{keyInt64, TensorContentTypeValues}, TypeInt64, shp, valuesInt64, true, true, true, false},

		{keyUint8, args{keyUint8, TensorContentTypeValues}, TypeUint8, shp, valuesUint8, true, true, true, false},
		{keyUint16, args{keyUint16, TensorContentTypeValues}, TypeUint16, shp, valuesUint16, true, true, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			gotResp, err := c.TensorGet(tt.args.name, tt.args.ct)
			if (err != nil) != tt.wantErr {
				t.Errorf("TensorGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {

				if diff := cmp.Diff(tt.wantDt, gotResp[0]); diff != "" {
					t.Errorf("TensorGet() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(tt.wantShape, gotResp[1]); diff != "" {
					t.Errorf("TensorGet() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(tt.wantData, gotResp[2]); diff != "" {
					t.Errorf("TensorGet() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestCommand_TensorGetBlob(t *testing.T) {
	valuesByte := []byte{1, 2, 3, 4}
	keyByte := "test:TensorGetBlog:[]byte:1"
	keyUnexistant := "test:TensorGetMeta:Unexistant"

	shp := []int64{1, 4}
	simpleClient := Connect("", createPool())
	simpleClient.TensorSet(keyByte, TypeInt8, shp, valuesByte)

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantDt    string
		wantShape []int64
		wantData  []byte
		wantErr   bool
	}{
		{keyByte, args{keyByte}, TypeInt8, shp, valuesByte, false},
		{keyUnexistant, args{keyUnexistant}, TypeInt8, shp, valuesByte, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			gotDt, gotShape, gotData, err := c.TensorGetBlob(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("TensorGetBlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == false {
				if gotDt != tt.wantDt {
					t.Errorf("TensorGetBlob() gotDt = %v, want %v", gotDt, tt.wantDt)
				}
				if !reflect.DeepEqual(gotShape, tt.wantShape) {
					t.Errorf("TensorGetBlob() gotShape = %v, want %v", gotShape, tt.wantShape)
				}
				if !reflect.DeepEqual(gotData, tt.wantData) {
					t.Errorf("TensorGetBlob() gotData = %v, want %v", gotData, tt.wantData)
				}
			}
		})
	}
}

func TestCommand_TensorGetMeta(t *testing.T) {

	keyFloat32 := "test:TensorGetMeta:TypeFloat32:1"
	keyFloat64 := "test:TensorGetMeta:TypeFloat64:1"

	keyInt8 := "test:TensorGetMeta:TypeInt8:1"
	keyInt16 := "test:TensorGetMeta:TypeInt16:1"
	keyInt32 := "test:TensorGetMeta:TypeInt32:1"
	keyInt64 := "test:TensorGetMeta:TypeInt64:1"

	keyUint8 := "test:TensorGetMeta:TypeUint8:1"
	keyUint16 := "test:TensorGetMeta:TypeUint16:1"

	keyUnexistant := "test:TensorGetMeta:Unexistant"

	shp := []int64{1, 2}
	simpleClient := Connect("", createPool())
	simpleClient.TensorSet(keyFloat32, TypeFloat32, shp, nil)
	simpleClient.TensorSet(keyFloat64, TypeFloat64, shp, nil)

	simpleClient.TensorSet(keyInt8, TypeInt8, shp, nil)
	simpleClient.TensorSet(keyInt16, TypeInt16, shp, nil)
	simpleClient.TensorSet(keyInt32, TypeInt32, shp, nil)
	simpleClient.TensorSet(keyInt64, TypeInt64, shp, nil)

	simpleClient.TensorSet(keyUint8, TypeUint8, shp, nil)
	simpleClient.TensorSet(keyUint16, TypeUint16, shp, nil)

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantDt    string
		wantShape []int64
		wantErr   bool
	}{
		{keyFloat32, args{keyFloat32}, TypeFloat32, shp, false},
		{keyFloat64, args{keyFloat64}, TypeFloat64, shp, false},
		{keyInt8, args{keyInt8}, TypeInt8, shp, false},
		{keyInt16, args{keyInt16}, TypeInt16, shp, false},
		{keyInt32, args{keyInt32}, TypeInt32, shp, false},
		{keyInt64, args{keyInt64}, TypeInt64, shp, false},
		{keyUnexistant, args{keyUnexistant}, TypeInt64, shp, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			gotDt, gotShape, err := c.TensorGetMeta(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("TensorGetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false {
				if gotDt != tt.wantDt {
					t.Errorf("TensorGetMeta() gotDt = %v, want %v", gotDt, tt.wantDt)
				}
				if !reflect.DeepEqual(gotShape, tt.wantShape) {
					t.Errorf("TensorGetMeta() gotShape = %v, want %v", gotShape, tt.wantShape)
				}
			}
		})
	}
}

func TestCommand_TensorGetValues(t *testing.T) {

	valuesFloat32 := []float32{1.1, 2.2, 3.3, 4.4, 5.5}
	valuesFloat64 := []float64{1.1}

	valuesInt8 := []int8{1}
	valuesInt16 := []int16{1}
	valuesInt32 := []int{1}
	valuesInt64 := []int64{1}

	valuesUint8 := []uint8{1}
	valuesUint16 := []uint16{1}
	keyFloat32 := "test:TensorGetValues:TypeFloat32:1"
	keyFloat64 := "test:TensorGetValues:TypeFloat64:1"

	keyInt8 := "test:TensorGetValues:TypeInt8:1"
	keyInt16 := "test:TensorGetValues:TypeInt16:1"
	keyInt32 := "test:TensorGetValues:TypeInt32:1"
	keyInt64 := "test:TensorGetValues:TypeInt64:1"

	keyUint8 := "test:TensorGetValues:TypeUint8:1"
	keyUint16 := "test:TensorGetValues:TypeUint16:1"
	keyUnexistant := "test:TensorGetValues:Unexistant"

	shp := []int64{1}
	shp2 := []int64{1, 5}
	simpleClient := Connect("", createPool())
	simpleClient.TensorSet(keyFloat32, TypeFloat32, shp2, valuesFloat32)
	simpleClient.TensorSet(keyFloat64, TypeFloat64, shp, valuesFloat64)

	simpleClient.TensorSet(keyInt8, TypeInt8, shp, valuesInt8)
	simpleClient.TensorSet(keyInt16, TypeInt16, shp, valuesInt16)
	simpleClient.TensorSet(keyInt32, TypeInt32, shp, valuesInt32)
	simpleClient.TensorSet(keyInt64, TypeInt64, shp, valuesInt64)

	simpleClient.TensorSet(keyUint8, TypeUint8, shp, valuesUint8)
	simpleClient.TensorSet(keyUint16, TypeUint16, shp, valuesUint16)

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantDt    string
		wantShape []int64
		wantData  interface{}
		wantErr   bool
	}{
		{keyFloat32, args{keyFloat32}, TypeFloat32, shp2, valuesFloat32, false},
		{keyFloat64, args{keyFloat64}, TypeFloat64, shp, valuesFloat64, false},

		{keyInt8, args{keyInt8}, TypeInt8, shp, valuesInt8, false},
		{keyInt16, args{keyInt16}, TypeInt16, shp, valuesInt16, false},
		{keyInt32, args{keyInt32}, TypeInt32, shp, valuesInt32, false},
		{keyInt64, args{keyInt64}, TypeInt64, shp, valuesInt64, false},

		{keyUint8, args{keyUint8}, TypeUint8, shp, valuesUint8, false},
		{keyUint16, args{keyUint16}, TypeUint16, shp, valuesUint16, false},

		{keyUnexistant, args{keyUnexistant}, TypeUint16, shp, valuesUint8, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			gotDt, gotShape, gotData, err := c.TensorGetValues(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("TensorGetValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false {
				if gotDt != tt.wantDt {
					t.Errorf("TensorGetValues() gotDt = %v, want %v", gotDt, tt.wantDt)
				}
				if !reflect.DeepEqual(gotShape, tt.wantShape) {
					t.Errorf("TensorGetValues() gotShape = %v, want %v", gotShape, tt.wantShape)
				}
				if !reflect.DeepEqual(gotData, tt.wantData) {
					t.Errorf("TensorGetValues() gotData = %v, want %v", gotData, tt.wantData)
				}
			}
		})
	}
}

func TestCommand_ModelSet(t *testing.T) {

	keyModelSet1 := "test:ModelSet:1"
	keyModelSet1Pipelined := "test:ModelSet:2:Pipelined"
	keyModelSetUnexistant := "test:ModelSet:3:Unexistant"
	dataUnexistant := []byte{}
	data, err := ioutil.ReadFile("./../tests/test_data/creditcardfraud.pb")
	if err != nil {
		t.Errorf("Error preparing for ModelSet(), while reading file. error = %v", err)
		return
	}

	type args struct {
		name    string
		backend string
		device  string
		data    []byte
		inputs  []string
		outputs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{keyModelSet1, args{keyModelSet1, BackendTF, DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"}}, false},
		{keyModelSet1Pipelined, args{keyModelSet1, BackendTF, DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"}}, false},
		{keyModelSetUnexistant, args{keyModelSetUnexistant, BackendTF, DeviceCPU, dataUnexistant, []string{"transaction", "reference"}, []string{"output"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			if err := c.ModelSet(tt.args.name, tt.args.backend, tt.args.device, tt.args.data, tt.args.inputs, tt.args.outputs); (err != nil) != tt.wantErr {
				t.Errorf("ModelSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_ModelStore(t *testing.T) {

	keyModelStore1 := "test:ModelStore:1"
	keyModelStore1Pipelined := "test:ModelStore:2:Pipelined"
	keyModelStoreUnexistant := "test:ModelStore:3:Unexistant"
	dataUnexistant := []byte{}
	data, err := ioutil.ReadFile("./../tests/test_data/creditcardfraud.pb")
	if err != nil {
		t.Errorf("Error preparing for ModelStore(), while reading file. error = %v", err)
		return
	}

	type args struct {
		name            string
		backend         string
		device          string
		tag             string
		batchsize       int64
		minbatchsize    int64
		minbatchtimeout int64
		data            []byte
		inputs          []string
		outputs         []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{keyModelStore1, args{keyModelStore1, BackendTF, DeviceCPU, "", 0, 0, 0, data, []string{"transaction", "reference"}, []string{"output"}}, false},
		{keyModelStore1Pipelined, args{keyModelStore1, BackendTF, DeviceCPU, "", 0, 0, 0, data, []string{"transaction", "reference"}, []string{"output"}}, false},
		{keyModelStore1Pipelined, args{keyModelStore1, BackendTF, DeviceCPU, "", 1, 1, 1, data, []string{"transaction", "reference"}, []string{"output"}}, false},
		{keyModelStoreUnexistant, args{keyModelStoreUnexistant, BackendTF, DeviceCPU, "", 0, 0, 0, dataUnexistant, []string{"transaction", "reference"}, []string{"output"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			if err := c.ModelStore(tt.args.name, tt.args.backend, tt.args.device, tt.args.tag, tt.args.batchsize, tt.args.minbatchsize, tt.args.minbatchtimeout, tt.args.inputs, tt.args.outputs, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ModelStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_ModelGet(t *testing.T) {
	keyModel1 := "test:ModelGetToModel:1"
	keyModelUnexistent1 := "test:ModelGetUnexistent:1"
	data, err := ioutil.ReadFile("./../tests/test_data/creditcardfraud.pb")
	if err != nil {
		t.Errorf("Error preparing for ModelGetToModel(), while issuing ModelSet. error = %v", err)
		return
	}
	simpleClient := createTestClient()
	err = simpleClient.ModelSet(keyModel1, BackendTF, DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"})
	if err != nil {
		t.Errorf("Error preparing for ModelGetToModel(), while issuing ModelSet. error = %v", err)
		return
	}
	type args struct {
		name string
	}
	tests := []struct {
		name             string
		args             args
		wantBackend      string
		wantDevice       string
		wantTag          string
		wantData         []byte
		wantBatchsize    int64
		wantMinbatchsize int64
		wantInputs       []string
		wantOutputs      []string
		wantErr          bool
	}{
		{keyModelUnexistent1, args{keyModelUnexistent1}, BackendTF, DeviceCPU, "", data, 0, 0, nil, nil, true},
		{keyModel1, args{keyModel1}, BackendTF, DeviceCPU, "", data, 0, 0, []string{"transaction", "reference"}, []string{"output"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := createTestClient()
			gotData, err := client.ModelGet(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModelGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotData[0], tt.wantBackend) {
					t.Errorf("ModelGet() gotBackend = %v, want %v. gotBackend Type %v, want Type %v.", gotData[0], tt.wantBackend, reflect.TypeOf(gotData[0]), reflect.TypeOf(tt.wantBackend))
				}
				if !reflect.DeepEqual(gotData[1], tt.wantDevice) {
					t.Errorf("ModelGet() gotDevice = %v, want %v. gotDevice Type %v, want Type %v.", gotData[1], tt.wantDevice, reflect.TypeOf(gotData[1]), reflect.TypeOf(tt.wantDevice))
				}
				if !reflect.DeepEqual(gotData[2], tt.wantTag) {
					t.Errorf("ModelGet() gotTag = %v, want %v. gotTag Type %v, want Type %v.", gotData[2], tt.wantTag, reflect.TypeOf(gotData[2]), reflect.TypeOf(tt.wantTag))
				}
				if !reflect.DeepEqual(gotData[3], tt.wantData) {
					t.Errorf("ModelGet() gotData = %v, want %v. gotData Type %v, want Type %v.", gotData[3], tt.wantData, reflect.TypeOf(gotData[3]), reflect.TypeOf(tt.wantData))
				}
				if !reflect.DeepEqual(gotData[4], tt.wantBatchsize) {
					t.Errorf("ModelGet() gotBatchsize = %v, want %v. gotBatchsize Type %v, want Type %v.", gotData[4], tt.wantBatchsize, reflect.TypeOf(gotData[4]), reflect.TypeOf(tt.wantBatchsize))
				}
				if !reflect.DeepEqual(gotData[5], tt.wantMinbatchsize) {
					t.Errorf("ModelGet() gotMinbatchsize = %v, want %v. gotMinbatchsize Type %v, want Type %v.", gotData[5], tt.wantMinbatchsize, reflect.TypeOf(gotData[5]), reflect.TypeOf(tt.wantMinbatchsize))
				}
				if !reflect.DeepEqual(gotData[6], tt.wantInputs) {
					t.Errorf("ModelGet() gotInputs = %v, want %v. gotInputs Type %v, want Type %v.", gotData[6], tt.wantInputs, reflect.TypeOf(gotData[6]), reflect.TypeOf(tt.wantInputs))
				}
				if !reflect.DeepEqual(gotData[7], tt.wantOutputs) {
					t.Errorf("ModelGet() gotOutputs = %v, want %v. gotOutputs Type %v, want Type %v.", gotData[7], tt.wantOutputs, reflect.TypeOf(gotData[7]), reflect.TypeOf(tt.wantOutputs))
				}
			}

		})
	}
}

func TestCommand_ModelDel(t *testing.T) {
	keyModel1 := "test:ModelDel:1"
	keyModel2 := "test:ModelDel:2:Pipelined"

	data, err := ioutil.ReadFile("./../tests/test_data/creditcardfraud.pb")
	if err != nil {
		t.Errorf("Error preparing for ModelDel(), while issuing ModelSet. error = %v", err)
		return
	}
	simpleClient := createTestClient()
	err = simpleClient.ModelSet(keyModel1, BackendTF, DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"})
	if err != nil {
		t.Errorf("Error preparing for ModelDel(), while issuing ModelSet. error = %v", err)
		return
	}
	err = simpleClient.ModelSet(keyModel2, BackendTF, DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"})
	if err != nil {
		t.Errorf("Error preparing for ModelDel(), while issuing ModelSet. error = %v", err)
		return
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{keyModel1, args{keyModel1}, false},
		{keyModel2, args{keyModel2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			err := c.ModelDel(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModelDel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_ModelExecute(t *testing.T) {
	// test ModelRun, ModelExecute and ModelExecuteWithTimeout
	keyModel1 := "test:ModelExecute:1"
	keyModel2 := "test:ModelExecute:2:Pipelined"
	keyTransaction1 := "test:ModelExecute:transaction:1"
	keyReference1 := "test:ModelExecute:reference:1"
	keyOutput1 := "test:ModelExecute:output:1"

	// preparing for execution
	data, err := ioutil.ReadFile("./../tests/test_data/creditcardfraud.pb")
	if err != nil {
		t.Errorf("Error preparing for ModelExecute(), while issuing ModelSet. error = %v", err)
		return
	}
	simpleClient := Connect("", createPool())
	err = simpleClient.ModelSet(keyModel1, BackendTF, DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"})
	if err != nil {
		t.Errorf("Error preparing for ModelExecute(), while issuing ModelSet. error = %v", err)
		return
	}

	err = simpleClient.ModelSet(keyModel2, BackendTF, DeviceCPU, data, []string{"transaction", "reference"}, []string{"output"})
	if err != nil {
		t.Errorf("Error preparing for ModelExecute(), while issuing ModelSet. error = %v", err)
		return
	}

	errortset := simpleClient.TensorSet(keyTransaction1, TypeFloat, []int64{1, 30}, nil)
	if errortset != nil {
		t.Error(errortset)
	}

	errortsetReference := simpleClient.TensorSet(keyReference1, TypeFloat, []int64{256}, nil)
	if errortsetReference != nil {
		t.Error(errortsetReference)
	}

	// execute
	type args struct {
		name    string
		inputs  []string
		outputs []string
		timeout int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Model1", args{keyModel1, []string{keyTransaction1, keyReference1}, []string{keyOutput1}, 0}, false},
		{"Model2", args{keyModel2, []string{keyTransaction1, keyReference1}, []string{keyOutput1}, 0}, false},
		{"Model2-with-timeout", args{keyModel2, []string{keyTransaction1, keyReference1}, []string{keyOutput1}, 3}, false},
		{"WrongInput", args{keyModel1, []string{keyTransaction1}, []string{keyOutput1}, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := createTestClient()
			if err := client.ModelRun(tt.args.name, tt.args.inputs, tt.args.outputs); (err != nil) != tt.wantErr {
				t.Errorf("ModelRun() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := client.ModelExecute(tt.args.name, tt.args.inputs, tt.args.outputs); (err != nil) != tt.wantErr {
				t.Errorf("ModelExecute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := client.ModelExecuteWithTimeout(tt.args.name, tt.args.inputs, tt.args.outputs, tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("ModelExecuteWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_FullFromModelFlow(t *testing.T) {
	model := implementations.NewModel(BackendTF, DeviceCPU)
	model.SetInputs([]string{"transaction", "reference"})
	model.SetOutputs([]string{"output"})
	client := createTestClient()
	err := model.SetBlobFromFile("./../tests/test_data/creditcardfraud.pb")
	assert.Nil(t, err)
	err = client.ModelSetFromModel("financialNet", model)
	assert.Nil(t, err)
	model1 := implementations.NewEmptyModel()
	err = client.ModelGetToModel("financialNet", model1)
	assert.Equal(t, model.Device(), model1.Device())
	assert.Nil(t, err)
	model1.SetInputs([]string{"transaction", "reference"})
	model1.SetOutputs([]string{"output"})
	model1.SetBatchSize(3)
	model1.SetMinBatchSize(1)
	model1.SetTag("financialTag")
	err = client.ModelStoreFromModel("financialNet1", model1)
	assert.Nil(t, err)
	model2 := implementations.NewEmptyModel()
	err = client.ModelGetToModel("financialNet1", model2)
	assert.Nil(t, err)
	assert.Equal(t, model1.Tag(), model2.Tag())
	assert.Equal(t, model1.BatchSize(), model2.BatchSize())
	assert.Equal(t, model1.MinBatchSize(), model2.MinBatchSize())
	assert.Equal(t, model1.Outputs(), model2.Outputs())
}

func TestCommand_ScriptDel(t *testing.T) {
	keyScript := "test:ScriptDel:1"
	keyScriptPipelined := "test:ScriptDel:2"
	keyScriptUnexistant := "test:ScriptDel:3:Unexistant"
	scriptBin := "def bar(a, b):\n    return a + b\n"
	simpleClient := Connect("", createPool())
	err := simpleClient.ScriptSet(keyScript, DeviceCPU, scriptBin)
	if err != nil {
		t.Errorf("Error preparing for ScriptDel(), while issuing ScriptSet. error = %v", err)
		return
	}
	type fields struct {
		Pool            *redis.Pool
		PipelineActive  bool
		PipelineMaxSize uint32
		PipelinePos     uint32
		ActiveConn      redis.Conn
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{keyScript, fields{createPool(), false, 0, 0, nil}, args{keyScript}, false},
		{keyScriptPipelined, fields{createPool(), true, 1, 0, nil}, args{keyScript}, false},
		{keyScriptUnexistant, fields{createPool(), false, 0, 0, nil}, args{keyScriptUnexistant}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Pool:                  tt.fields.Pool,
				PipelineActive:        tt.fields.PipelineActive,
				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
				PipelinePos:           tt.fields.PipelinePos,
				ActiveConn:            tt.fields.ActiveConn,
			}
			if err := c.ScriptDel(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ScriptDel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_ScriptGet(t *testing.T) {
	keyScript := "test:ScriptGet:1"
	keyScriptPipelined := "test:ScriptGet:2"
	keyScriptEmpty := "test:ScriptGet:3:Empty"
	scriptBin := ""
	simpleClient := Connect("", createPool())
	err := simpleClient.ScriptSet(keyScript, DeviceCPU, scriptBin)
	if err != nil {
		t.Errorf("Error preparing for ScriptGet(), while issuing ScriptSet. error = %v", err)
		return
	}

	keyScript2 := "test:ScriptGet:2"
	keyScriptTag := "keyScriptTag"
	err = simpleClient.ScriptSetWithTag(keyScript2, DeviceCPU, scriptBin, keyScriptTag)
	if err != nil {
		t.Errorf("Error preparing for ScriptGet(), while issuing ScriptSet. error = %v", err)
		return
	}

	type args struct {
		name string
	}
	tests := []struct {
		name           string
		args           args
		wantDeviceType string
		wantData       string
		wantTag        string
		wantEntryPoint []string
		wantErr        bool
	}{
		{keyScript, args{keyScript}, DeviceCPU, "", "", []string{}, false},
		{keyScriptPipelined, args{keyScript}, DeviceCPU, "", "", []string{}, false},
		{keyScriptEmpty, args{keyScriptEmpty}, DeviceCPU, "", "", []string{}, true},
		{keyScriptTag, args{keyScript2}, DeviceCPU, "", keyScriptTag, []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			gotData, err := c.ScriptGet(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScriptGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == false {
				if !reflect.DeepEqual(gotData[0], tt.wantDeviceType) {
					t.Errorf("ScriptGet() gotData = %v, want %v", gotData[0], tt.wantDeviceType)
				}
				if !reflect.DeepEqual(gotData[1], tt.wantTag) {
					t.Errorf("ScriptGet() gotData = %v, want %v", gotData[1], tt.wantTag)
				}
				if !reflect.DeepEqual(gotData[2], tt.wantData) {
					t.Errorf("ScriptGet() gotData = %v, want %v", gotData[2], tt.wantData)
				}
				if !reflect.DeepEqual(gotData[3], tt.wantEntryPoint) {
					t.Errorf("ScriptGet() gotData = %v, want %v", gotData[3], tt.wantEntryPoint)
				}
			}

		})
	}
}

func TestCommand_ScriptGetToInterface(t *testing.T) {
	keyScript := "test:ScriptGet:1"
	simpleClient := Connect("", createPool())
	tag := "bar"

	err := simpleClient.ScriptSetWithTag(keyScript, DeviceCPU, script, tag)
	if err != nil {
		t.Errorf("Error preparing for ScriptGetToInterface(), while issuing ScriptSet. error = %v", err)
		return
	}

	c := createTestClient()
	scriptInterface := implementations.NewEmptyScript()
	err = c.ScriptGetToInterface(keyScript, scriptInterface)
	assert.Nil(t, err)
	assert.Equal(t, scriptInterface.Device(), DeviceCPU)
	assert.Equal(t, scriptInterface.Tag(), tag)
	assert.Equal(t, scriptInterface.EntryPoints(), []string{})
	assert.Equal(t, scriptInterface.Source(), script)
}

func TestCommand_ScriptRun(t *testing.T) {
	keyScript1 := "test:ScriptRun:1"
	keyScript2 := "test:ScriptRun:2:Pipelined"
	keyScript3Empty := "test:ScriptRun:3:Empty"
	scriptBin := "def bar(a, b):\n    return a + b\n"
	tag := "bar"
	simpleClient := Connect("", createPool())

	scriptInterface := implementations.NewEmptyScript()
	scriptInterface.SetDevice(DeviceCPU)
	scriptInterface.SetTag(tag)
	scriptInterface.SetSource(scriptBin)
	err := simpleClient.ScriptSetFromInteface(keyScript1, scriptInterface)
	if err != nil {
		t.Errorf("Error preparing for ScriptRun(), while issuing ScriptSet. error = %v", err)
		return
	}
	err = simpleClient.ScriptSet(keyScript2, DeviceCPU, scriptBin)
	if err != nil {
		t.Errorf("Error preparing for ScriptRun(), while issuing ScriptSet. error = %v", err)
		return
	}

	type fields struct {
		Pool            *redis.Pool
		PipelineActive  bool
		PipelineMaxSize uint32
		PipelinePos     uint32
		ActiveConn      redis.Conn
	}
	type args struct {
		name    string
		fn      string
		inputs  []string
		outputs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{keyScript2, fields{createPool(), true, 1, 0, nil}, args{keyScript2, "", []string{""}, []string{""}}, false},
		{keyScript3Empty, fields{createPool(), false, 0, 0, nil}, args{keyScript3Empty, "", []string{""}, []string{""}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Pool:                  tt.fields.Pool,
				PipelineActive:        tt.fields.PipelineActive,
				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
				PipelinePos:           tt.fields.PipelinePos,
				ActiveConn:            tt.fields.ActiveConn,
			}
			if err := c.ScriptRun(tt.args.name, tt.args.fn, tt.args.inputs, tt.args.outputs); (err != nil) != tt.wantErr {
				t.Errorf("ScriptRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_ScriptSet(t *testing.T) {
	keyScriptError := "test:ScriptSet:Error:1"
	scriptBin := "import abc"

	type args struct {
		name   string
		device string
		data   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{keyScriptError, args{keyScriptError, DeviceCPU, scriptBin}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			if err := c.ScriptSet(tt.args.name, tt.args.device, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ScriptSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_ScriptStore(t *testing.T) {
	keyScriptError := "test:ScriptStore:Error:1"
	scriptBin := "import abc"
	tag := "bar"

	type args struct {
		name        string
		device      string
		data        string
		entryPoints []string
		tag         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no-entrypoints", args{keyScriptError, DeviceCPU, scriptBin, nil, tag}, true},
		{"no-entrypoints", args{keyScriptError, DeviceCPU, scriptBin, []string{}, tag}, true},
		{"non-existing-entrypoint", args{keyScriptError, DeviceCPU, scriptBin, []string{"foo"}, tag}, true},
		{"existing-entrypoint", args{keyScriptError, DeviceCPU, script, []string{"bar"}, ""}, false},
		{"missing-cuda", args{keyScriptError, DeviceGPU, script, []string{"bar"}, ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			if err := c.ScriptStore(tt.args.name, tt.args.device, tt.args.data, tt.args.entryPoints); (err != nil) != tt.wantErr {
				t.Errorf("ScriptSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			if err := c.ScriptStoreWithTag(tt.args.name, tt.args.device, tt.args.data, tt.args.entryPoints, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("ScriptSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_ScriptStoreFromInteface_Flow(t *testing.T) {
	client := createTestClient()
	client.Flush()

	scriptInteface := implementations.NewScript(DeviceCPU)
	scriptInteface.SetSource(script)
	scriptInteface.SetEntryPoints([]string{"bar", "bar_variadic"})
	err := client.ScriptStoreFromInteface("myscript", scriptInteface)
	assert.Nil(t, err)
	scriptInteface.SetTag("{1}")
	err = client.ScriptStoreFromInteface("mynewscript", scriptInteface)
	assert.Nil(t, err)
	scriptInteface1 := implementations.NewEmptyScript()
	err = client.ScriptGetToInterface("mynewscript", scriptInteface1)
	assert.Nil(t, err)
	assert.Equal(t, scriptInteface.Device(), scriptInteface1.Device())
	assert.Equal(t, scriptInteface.Tag(), scriptInteface1.Tag())
	assert.Equal(t, scriptInteface.EntryPoints(), scriptInteface1.EntryPoints())
}

func TestCommand_ScriptExecute(t *testing.T) {
	// execute script using ScriptExecute
	keyScript := "test:myscript"
	keyOutput := "myOutput"
	client := createTestClient()
	client.Flush()

	// preparing for execution
	client.ScriptStore(keyScript, DeviceCPU, script, []string{"bar"})
	client.TensorSet("a", TypeFloat, []int64{2, 2}, []float32{2, 3, 2, 3})
	client.TensorSet("b", TypeFloat, []int64{2, 2}, []float32{2, 3, 2, 3})
	// execute and test the result in the output
	client.ScriptExecute(keyScript, "bar", nil, []string{"a", "b"}, nil, []string{keyOutput})
	gotResp, err := client.TensorGet(keyOutput, TensorContentTypeValues)
	assert.Nil(t, err)
	if diff := cmp.Diff(TypeFloat32, gotResp[0]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]int64{2, 2}, gotResp[1]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]float32{4, 6, 4, 6}, gotResp[2]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
}

func TestCommand_ScriptExecute_RedisCommands(t *testing.T) {
	// test ScriptExecute and ScriptExecuteWithTimeout using a script that comprising redis commands
	client := createTestClient()
	client.Flush()
	keyScript := "test:myscript:rediscommands"
	client.ScriptStore(keyScript, DeviceCPU, scriptWithRedisCommands, []string{"int_set_get", "func"})

	// execute and test the result in the output
	client.ScriptExecute(keyScript, "int_set_get", []string{"x{1}", "{1}"}, nil, []string{"3"}, []string{"y{1}"})
	gotResp, err := client.TensorGet("y{1}", TensorContentTypeValues)
	assert.Nil(t, err)
	if diff := cmp.Diff(TypeInt64, gotResp[0]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]int64{3}, gotResp[2]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}

	// preparing for execution
	client.TensorSet("mytensor1{1}", TypeFloat, []int64{1}, []float32{40})
	client.TensorSet("mytensor2{1}", TypeFloat, []int64{1}, []float32{10})
	client.TensorSet("mytensor3{1}", TypeFloat, []int64{1}, []float32{1})
	// execute and test the result in the output
	client.ScriptExecuteWithTimeout(keyScript, "func", []string{"key{1}"}, []string{"mytensor1{1}", "mytensor2{1}", "mytensor3{1}"}, []string{"3"}, []string{"my_output{1}"}, 1000)
	gotResp, err = client.TensorGet("my_output{1}", TensorContentTypeValues)
	assert.Nil(t, err)
	if diff := cmp.Diff(TypeFloat, gotResp[0]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]float32{54}, gotResp[2]); diff != "" {
		t.Errorf("TestCommand_FullFromTensor() mismatch (-want +got):\n%s", diff)
	}
	reply, _ := client.DoOrSend("GET", redis.Args{"key{1}"}, err)
	assert.Nil(t, reply)
}

func TestCommand_LoadBackend(t *testing.T) {
	keyTest1 := "test:LoadBackend:1:Unexistent"
	keyTest2 := "test:LoadBackend:2:Unexistent:Pipelined"
	type args struct {
		backend_identifier string
		location           string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{keyTest1, args{BackendTF, "unexistant"}, true},
		{keyTest2, args{BackendTF, "unexistant"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			err := c.LoadBackend(tt.args.backend_identifier, tt.args.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBackend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_Info(t *testing.T) {
	c := createTestClient()
	keyModel1 := "test:Info:1"
	data, err := ioutil.ReadFile("./../tests/test_data/graph.pb")
	if err != nil {
		t.Errorf("Error preparing for Info(), while issuing ModelSet. error = %v", err)
		return
	}
	err = c.ModelSet(keyModel1, BackendTF, DeviceCPU, data, []string{"a", "b"}, []string{"mul"})
	if err != nil {
		t.Errorf("Error preparing for Info(), while issuing ModelSet. error = %v", err)
		return
	}

	// first inited info
	info, err := c.Info(keyModel1)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, keyModel1, info["key"])
	assert.Equal(t, DeviceCPU, info["device"])
	assert.Equal(t, BackendTF, info["backend"])
	assert.Equal(t, "0", info["calls"])

	err = c.TensorSet("a", TypeFloat32, []int64{1}, []float32{1.1})
	assert.Nil(t, err)
	err = c.TensorSet("b", TypeFloat32, []int64{1}, []float32{4.4})
	assert.Nil(t, err)
	err = c.ModelExecute(keyModel1, []string{"a", "b"}, []string{"mul"})
	assert.Nil(t, err)
	info, _ = c.Info(keyModel1)
	// one model runs
	assert.Equal(t, "1", info["calls"])

	// reset
	ret, _ := c.ResetStat(keyModel1)
	assert.Equal(t, "OK", ret)
	info, _ = c.Info(keyModel1)
	assert.Equal(t, "0", info["calls"])

	// not exits
	_, err = c.ResetStat("notExits")
	assert.NotNil(t, err)
}

func TestCommand_DagRun(t *testing.T) {
	c := createTestClient()
	keyModel1 := "test:DagRun:mymodel:1"
	data, err := ioutil.ReadFile("./../tests/test_data/graph.pb")
	if err != nil {
		t.Errorf("Error preparing for Info(), while issuing ModelSet. error = %v", err)
		return
	}
	err = c.ModelSet(keyModel1, BackendTF, DeviceCPU, data, []string{"a", "b"}, []string{"mul"})
	assert.Nil(t, err)
	err = c.TensorSet("persisted_tensor_1", TypeFloat32, []int64{1, 2}, []float32{5, 10})
	assert.Nil(t, err)

	type args struct {
		loadKeys            []string
		persistKeys         []string
		dagCommandInterface DagCommandInterface
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"t_wrong_number", args{[]string{"notnumber"}, nil, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, true},
		{"t_load", args{[]string{"persisted_tensor_1"}, []string{"tensor1"}, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, false},
		{"t_load_err", args{[]string{"not_exits_tensor"}, []string{"tensor1"}, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, true},
		{"t1", args{nil, nil, NewDag().TensorSet("a", TypeFloat32, []int64{1}, []float32{1.1})}, false},
		{"t_blob_run", args{nil, nil, NewDag().TensorSet("a", TypeFloat32, []int64{1}, []float32{1.1}).TensorSet("b", TypeFloat32, []int64{1}, []float32{4.4}).ModelRun("test:DagRun:mymodel:1", []string{"a", "b"}, []string{"mul"}).TensorGet("mul", TensorContentTypeBlob)}, false},
		{"t_blob_execute", args{nil, nil, NewDag().TensorSet("a", TypeFloat32, []int64{1}, []float32{1.1}).TensorSet("b", TypeFloat32, []int64{1}, []float32{4.4}).ModelExecute("test:DagRun:mymodel:1", []string{"a", "b"}, []string{"mul"}, 0).TensorGet("mul", TensorContentTypeBlob)}, true},
		{"t_values", args{nil, nil, NewDag().TensorSet("mytensor", TypeFloat32, []int64{1, 2}, []int64{5, 10}).TensorGet("mytensor", TensorContentTypeValues)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			results, err := c.DagRun(tt.args.loadKeys, tt.args.persistKeys, tt.args.dagCommandInterface)
			if (err != nil) != tt.wantErr {
				t.Errorf("DagRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, result := range results {
				ret, ok := result.(string)
				if ok {
					assert.Equal(t, "OK", ret)
					continue
				}
				values, ok := result.([]interface{})
				if ok {
					vs, _ := redis.Strings(values, nil)
					assert.True(t, len(vs) > 0)
					continue
				}
				blobs, ok := result.([]byte)
				if ok {
					assert.True(t, len(blobs) > 0)
					continue
				}
				t.Errorf("DagRun() error unsupported result")
			}
		})
	}
}

func TestCommand_DagRunRO(t *testing.T) {
	c := createTestClient()
	err := c.TensorSet("persisted_tensor", TypeFloat32, []int64{1, 2}, []float32{5, 10})
	assert.Nil(t, err)
	type args struct {
		loadKeys            []string
		dagCommandInterface DagCommandInterface
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"t_1", args{[]string{"persisted_tensor"}, NewDag().TensorGet("persisted_tensor", TensorContentTypeValues)}, false},
		{"t_2", args{nil, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10}).TensorSet("tensor2", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, false},
		{"t_err1", args{[]string{"notnumber"}, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			results, err := c.DagRunRO(tt.args.loadKeys, tt.args.dagCommandInterface)
			if (err != nil) != tt.wantErr {
				t.Errorf("DagRunRO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, result := range results {
				ret, ok := result.(string)
				if ok {
					assert.Equal(t, "OK", ret)
					continue
				}
				values, ok := result.([]interface{})
				if ok {
					vs, _ := redis.Strings(values, nil)
					assert.True(t, len(vs) > 0)
					continue
				}
				blobs, ok := result.([]byte)
				if ok {
					assert.True(t, len(blobs) > 0)
					continue
				}
				t.Errorf("DagRunRO() error unsupported result")
			}
		})
	}
}

func TestCommand_DagExecute(t *testing.T) {
	c := createTestClient()
	keyModel1 := "test:DagRun:mymodel:1"
	data, err := ioutil.ReadFile("./../tests/test_data/graph.pb")
	if err != nil {
		t.Errorf("Error preparing for Info(), while issuing ModelSet. error = %v", err)
		return
	}
	err = c.ModelSet(keyModel1, BackendTF, DeviceCPU, data, []string{"a", "b"}, []string{"mul"})
	assert.Nil(t, err)
	err = c.TensorSet("persisted_tensor_1", TypeFloat32, []int64{1, 2}, []float32{5, 10})
	assert.Nil(t, err)

	type args struct {
		loadKeys            []string
		persistKeys         []string
		routing             string
		timeout             int64
		dagCommandInterface DagCommandInterface
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"t_wrong_number", args{[]string{"notnumber"}, nil, "", 0, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, true},
		{"t_wrong_arguments", args{nil, nil, "", 0, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, true},
		{"t_load", args{[]string{"persisted_tensor_1"}, []string{"tensor1"}, "", 0, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, false},
		{"t_load_err", args{[]string{"not_exits_tensor"}, []string{"tensor1"}, "", 0, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, true},
		{"t1", args{nil, []string{"a"}, "", 0, NewDag().TensorSet("a", TypeFloat32, []int64{1}, []float32{1.1})}, false},
		{"t_blob_run", args{nil, []string{"mul"}, "", 0, NewDag().TensorSet("a", TypeFloat32, []int64{1}, []float32{1.1}).TensorSet("b", TypeFloat32, []int64{1}, []float32{4.4}).ModelRun("test:DagRun:mymodel:1", []string{"a", "b"}, []string{"mul"}).TensorGet("mul", TensorContentTypeBlob)}, true},
		{"t_blob_execute", args{nil, []string{"mul"}, "", 0, NewDag().TensorSet("a", TypeFloat32, []int64{1}, []float32{1.1}).TensorSet("b", TypeFloat32, []int64{1}, []float32{4.4}).ModelExecute("test:DagRun:mymodel:1", []string{"a", "b"}, []string{"mul"}, 0).TensorGet("mul", TensorContentTypeBlob)}, false},
		{"t_values", args{nil, []string{"mytensor"}, "", 0, NewDag().TensorSet("mytensor", TypeFloat32, []int64{1, 2}, []int64{5, 10}).TensorGet("mytensor", TensorContentTypeValues)}, false},
		{"t_values_timeout", args{nil, []string{"mytensor"}, "", 1000, NewDag().TensorSet("mytensor", TypeFloat32, []int64{1, 2}, []int64{5, 10}).TensorGet("mytensor", TensorContentTypeValues)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			results, err := c.DagExecute(tt.args.loadKeys, tt.args.persistKeys, tt.args.routing, tt.args.timeout, tt.args.dagCommandInterface)
			if (err != nil) != tt.wantErr {
				t.Errorf("DagRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, result := range results {
				ret, ok := result.(string)
				if ok {
					assert.Equal(t, "OK", ret)
					continue
				}
				values, ok := result.([]interface{})
				if ok {
					vs, _ := redis.Strings(values, nil)
					assert.True(t, len(vs) > 0)
					continue
				}
				blobs, ok := result.([]byte)
				if ok {
					assert.True(t, len(blobs) > 0)
					continue
				}
				t.Errorf("DagRun() error unsupported result")
			}
		})
	}
}

func TestCommand_DagExecuteRO(t *testing.T) {
	c := createTestClient()
	err := c.TensorSet("persisted_tensor", TypeFloat32, []int64{1, 2}, []float32{5, 10})
	assert.Nil(t, err)
	type args struct {
		loadKeys            []string
		routing             string
		timeout             int64
		dagCommandInterface DagCommandInterface
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"t_1", args{[]string{"persisted_tensor"}, "", 0, NewDag().TensorGet("persisted_tensor", TensorContentTypeValues)}, false},
		{"t_1_timeout", args{[]string{"persisted_tensor"}, "", 1000, NewDag().TensorGet("persisted_tensor", TensorContentTypeValues)}, false},
		{"t_2", args{nil, "{2}", 0, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10}).TensorSet("tensor{2}", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, false},
		{"t_err1", args{[]string{"notnumber"}, "", 0, NewDag().TensorSet("tensor1", TypeFloat32, []int64{1, 2}, []int64{5, 10})}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestClient()
			results, err := c.DagExecuteRO(tt.args.loadKeys, tt.args.routing, tt.args.timeout, tt.args.dagCommandInterface)
			if (err != nil) != tt.wantErr {
				t.Errorf("DagRunRO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, result := range results {
				ret, ok := result.(string)
				if ok {
					assert.Equal(t, "OK", ret)
					continue
				}
				values, ok := result.([]interface{})
				if ok {
					vs, _ := redis.Strings(values, nil)
					assert.True(t, len(vs) > 0)
					continue
				}
				blobs, ok := result.([]byte)
				if ok {
					assert.True(t, len(blobs) > 0)
					continue
				}
				t.Errorf("DagRunRO() error unsupported result")
			}
		})
	}
}

func TestClient_ModelRun(t *testing.T) {
	type fields struct {
		Pool                  *redis.Pool
		PipelineActive        bool
		PipelineAutoFlushSize uint32
		PipelinePos           uint32
		ActiveConn            redis.Conn
	}
	type args struct {
		name              string
		inputTensorNames  []string
		outputTensorNames []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Pool:                  tt.fields.Pool,
				PipelineActive:        tt.fields.PipelineActive,
				PipelineAutoFlushSize: tt.fields.PipelineAutoFlushSize,
				PipelinePos:           tt.fields.PipelinePos,
				ActiveConn:            tt.fields.ActiveConn,
			}
			if err := c.ModelRun(tt.args.name, tt.args.inputTensorNames, tt.args.outputTensorNames); (err != nil) != tt.wantErr {
				t.Errorf("ModelRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_SetBackendsPath(t *testing.T) {
	c := createTestClient()
	ret, err := c.SetBackendsPath("/usr/lib/redis/modules/backends/")
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
}
