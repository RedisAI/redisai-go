package redisai

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"strconv"
)

type aiclient interface {
	LoadBackend(backend_identifier string, location string) (err error)
	TensorSet(name string, dt string, shape []int, data interface{}) (err error)
	TensorGet(name string, ct string) (data []interface{}, err error)
	ModelSet(name string, backend string, device string, data []byte, inputs []string, outputs []string) (err error)
	ModelGet(name string) (data []interface{}, err error)
	ModelDel(name string) (err error)
	ModelRun(name string, inputs []string, outputs []string) (err error)
	ScriptSet(name string, device string, data string) (err error)
	ScriptGet(name string) (data []interface{}, err error)
	ScriptDel(name string) (err error)
	ScriptRun(name string, fn string, inputs []string, outputs []string) (err error)
}

const (
	// BackendTF represents a TensorFlow backend
	BackendTF = string("TF")
	// BackendTorch represents a Torch backend
	BackendTorch = string("TORCH")
	// BackendONNX represents an ONNX backend
	BackendONNX = string("ORT")

	// DeviceCPU represents a CPU device
	DeviceCPU = string("CPU")
	// DeviceGPU represents a GPU device
	DeviceGPU = string("GPU")

	// TypeFloat represents a float type
	TypeFloat = string("FLOAT")
	// TypeDouble represents a double type
	TypeDouble = string("DOUBLE")
	// TypeInt8 represents a int8 type
	TypeInt8 = string("INT8")
	// TypeInt16 represents a int16 type
	TypeInt16 = string("INT16")
	// TypeInt32 represents a int32 type
	TypeInt32 = string("INT32")
	// TypeInt64 represents a int64 type
	TypeInt64 = string("INT64")
	// TypeUint8 represents a uint8 type
	TypeUint8 = string("UINT8")
	// TypeUint16 represents a uint16 type
	TypeUint16 = string("UINT16")
	// TypeFloat32 is an alias for float
	TypeFloat32 = string("FLOAT")
	// TypeFloat64 is an alias for double
	TypeFloat64 = string("DOUBLE")

	// TensorContentTypeBLOB is an alias for BLOB tensor content
	TensorContentTypeBlob = string("BLOB")

	// TensorContentTypeBLOB is an alias for BLOB tensor content
	TensorContentTypeValues = string("VALUES")

	// TensorContentTypeBLOB is an alias for BLOB tensor content
	TensorContentTypeMeta = string("META")
)

var ErrNil = errors.New("redisai-go: nil returned")

func TensorSetArgs(name string, dt string, dims []int, data interface{}, includeCommandName bool) (redis.Args, error) {
	args := redis.Args{}
	var err error = nil
	if includeCommandName {
		args = args.Add("AI.TENSORSET")
	}
	args = args.Add(name, dt).AddFlat(dims)
	if data != nil {
		var dtype = reflect.TypeOf(data)
		switch dtype {
		case reflect.TypeOf(([]uint8)(nil)):
			fallthrough
		case reflect.TypeOf(([]byte)(nil)):
			args = args.Add("BLOB", data)
		case reflect.TypeOf(""):
			fallthrough
		case reflect.TypeOf(([]int)(nil)):
			fallthrough
		case reflect.TypeOf(([]int8)(nil)):
			fallthrough
		case reflect.TypeOf(([]int16)(nil)):
			fallthrough
		case reflect.TypeOf(([]int32)(nil)):
			fallthrough
		case reflect.TypeOf(([]int64)(nil)):
			fallthrough
		case reflect.TypeOf(([]uint)(nil)):
			fallthrough
		case reflect.TypeOf(([]uint16)(nil)):
			fallthrough
		case reflect.TypeOf(([]float32)(nil)):
			fallthrough
		case reflect.TypeOf(([]float64)(nil)):
			args = args.Add("VALUES").AddFlat(data)
			// unsupported data type
		case reflect.TypeOf(([]uint32)(nil)):
			fallthrough
			// unsupported data type
		case reflect.TypeOf(([]uint64)(nil)):
			fallthrough
			// unsupported data type
		default:
			err = fmt.Errorf("redisai.TensorSetArgs: AI.TENSOR does not support the following type %v", reflect.TypeOf(data))
		}

	}
	return args, err
}



// DataType is a helper that converts a command reply to a DataType.
func replyDataType(reply interface{}, err error) (dt string, outputErr error) {
	if err != nil {
		return "", err
	}
	switch reply := reply.(type) {
	case string:
		switch reply {
		case "FLOAT":
			dt = TypeFloat
		case "DOUBLE":
			dt = TypeDouble
		case "INT8":
			dt = TypeInt8
		case "INT16":
			dt = TypeInt16
		case "INT32":
			dt = TypeInt32
		case "INT64":
			dt = TypeInt64
		case "UINT8":
			dt = TypeUint8
		case "UINT16":
			dt = TypeUint16
		}
		return dt, nil
	case nil:
		return "", ErrNil

	}
	return "", fmt.Errorf("redisai-go: unexpected type for replyDataType, got type %T", reply)
}

func ProcessTensorReplyMeta(resp interface{}, err error) (data []interface{}, outErr error) {
	data, outErr = redis.Values(resp, err)
	if len(data) < 2 {
		err = fmt.Errorf("redisai.TensorGet: AI.TENSORGET returned response with incorrect sizing. expected at least '%d' got '%d'", 2, len(data))
		return data, err
	}
	data[0], outErr = redis.String(data[0], err)
	data[1], outErr = redis.Ints(data[1], err)
	return data, outErr
}

func ProcessTensorReplyBlob(resp []interface{}, err error) ([]interface{}, error) {
	if len(resp) < 3 {
		err = fmt.Errorf("redisai.TensorGet: AI.TENSORGET returned response with incorrect sizing. expected '%d' got '%d'", 3, len(resp))
		return resp, err
	}
	switch resp[0].(string) {
	case TypeFloat:
		resp[2], err = Float32sBytes(resp[2], resp[1].([]int), err)
	case TypeUint8:
		resp[2], err = redis.Bytes(resp[2], err)
	}
	return resp, err
}

func ProcessTensorReplyValues(resp []interface{}, err error) ([]interface{}, error) {
	if len(resp) < 3 {
		err = fmt.Errorf("redisai.TensorGet: AI.TENSORGET returned response with incorrect sizing. expected '%d' got '%d'", 3, len(resp))
		return resp, err
	}
	switch resp[0].(string) {
	case TypeFloat:
		resp[2], err = Float32s(resp[2], err)
	case TypeDouble:
		resp[2], err = redis.Float64s(resp[2], err)
	case TypeInt8:
		resp[2], err = Int8s(resp[2], err)
	case TypeInt16:
		resp[2], err = Int16s(resp[2], err)
	case TypeInt32:
		resp[2], err = redis.Ints(resp[2], err)
	case TypeInt64:
		resp[2], err = redis.Int64s(resp[2], err)
	case TypeUint8:
		resp[2], err = Uint8s(resp[2], err)
	case TypeUint16:
		resp[2], err = Uint16s(resp[2], err)
	}

	//resp[2], err = redis.Values(resp[2], err)

	return resp, err
}

// Float32s is a helper that converts an array command reply to a []float32.
func Float32sBytes(reply interface{}, dimension []int, err error) ([]float32, error) {
	totalResults := 0
	if len(dimension) > 0 {
		totalResults = dimension[0]
	}

	for i := 1; i < len(dimension); i++ {
		totalResults *= dimension[i]
	}
	fmt.Errorf(fmt.Sprintf("Total results %d", totalResults))

	var result = make([]float32, totalResults)
	tr, err := redis.Bytes(reply, err)
	if err != nil {
		return result, err
	}
	buf := bytes.NewReader(tr)
	fmt.Errorf(fmt.Sprintf("Total results %d", totalResults))
	for i := 0; i < totalResults; i++ {
		err := binary.Read(buf, binary.LittleEndian, &result[i])
		if err != nil {
			return result, err
		}
	}
	return result, err
}

// Float32s is a helper that converts an array command reply to a []float32.
func Float32s(reply interface{}, err error) ([]float32, error) {
	var result []float32
	err = sliceHelper(reply, err, "Float32s", func(n int) { result = make([]float32, n) }, func(i int, v interface{}) error {
		p, ok := v.([]byte)
		if !ok {
			return fmt.Errorf("redisai-go: unexpected element type for Float32s, got type %T", v)
		}
		var f, err = strconv.ParseFloat(string(p), 64)
		result[i] = float32(f)
		return err
	})
	return result, err
}

// Uint16s is a helper that converts an array command reply to a []uint16.
func Uint16s(reply interface{}, err error) ([]uint16, error) {
	var result []uint16
	tr, err := redis.Values(reply, err)
	if err != nil {
		return result, err
	}
	for _, num := range tr {
		result = append(result, uint16(num.(int64)))
	}
	return result, err
}

// Int16s is a helper that converts an array command reply to a []int16.
func Int16s(reply interface{}, err error) ([]int16, error) {
	var result []int16
	tr, err := redis.Values(reply, err)
	if err != nil {
		return result, err
	}
	for _, num := range tr {
		result = append(result, int16(num.(int64)))
	}
	return result, err
}

// Uint8s is a helper that converts an array command reply to a []uint8.
func Uint8s(reply interface{}, err error) ([]uint8, error) {
	var result []uint8
	tr, err := redis.Values(reply, err)
	if err != nil {
		return result, err
	}
	for _, num := range tr {
		result = append(result, uint8(num.(int64)))
	}
	return result, err
}

// Int8s is a helper that converts an array command reply to a []int8.
func Int8s(reply interface{}, err error) ([]int8, error) {
	var result []int8
	tr, err := redis.Values(reply, err)
	if err != nil {
		return result, err
	}
	for _, num := range tr {
		result = append(result, int8(num.(int64)))
	}
	return result, err
}

func sliceHelper(reply interface{}, err error, name string, makeSlice func(int), assign func(int, interface{}) error) error {
	if err != nil {
		return err
	}
	switch reply := reply.(type) {
	case []interface{}:
		makeSlice(len(reply))
		for i := range reply {
			if reply[i] == nil {
				continue
			}
			if err := assign(i, reply[i]); err != nil {
				return err
			}
		}
		return nil
	case nil:
		return ErrNil
	}
	return fmt.Errorf("redisai-go: unexpected type for %s, got type %T", name, reply)
}

func float32ToByte(f float32) (converted []byte, err error) {
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, f)
	converted = buf.Bytes()
	return
}
