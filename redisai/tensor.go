package redisai

import (
	"fmt"
	"github.com/RedisAI/redisai-go/redisai/converters"
	"github.com/gomodule/redigo/redis"
	"reflect"
)

// TensorInterface is an interface that represents the skeleton of a tensor ( n-dimensional array of numerical data )
// needed to map it to a RedisAI Model with the proper operations
type TensorInterface interface {

	// Shape returns the size - in each dimension - of the tensor.
	Shape() []int64

	SetShape(shape []int64)

	// NumDims returns the number of dimensions of the tensor.
	NumDims() int64

	// Len returns the number of elements in the tensor.
	Len() int64

	Dtype() reflect.Type

	// Data returns the underlying tensor data
	Data() interface{}
	SetData(interface{})
}

func TensorGetTypeStrFromType(dtype reflect.Type) (typestr string, err error) {
	switch dtype {
	case reflect.TypeOf(([]uint8)(nil)):
		typestr = TypeUint8
	case reflect.TypeOf(([]byte)(nil)):
		typestr = TypeUint8
	case reflect.TypeOf(([]int)(nil)):
		typestr = TypeInt32
	case reflect.TypeOf(([]int8)(nil)):
		typestr = TypeInt8
	case reflect.TypeOf(([]int16)(nil)):
		typestr = TypeInt16
	case reflect.TypeOf(([]int32)(nil)):
		typestr = TypeInt32
	case reflect.TypeOf(([]int64)(nil)):
		typestr = TypeInt64
	case reflect.TypeOf(([]uint)(nil)):
		typestr = TypeUint8
	case reflect.TypeOf(([]uint16)(nil)):
		typestr = TypeUint16
	case reflect.TypeOf(([]float32)(nil)):
		typestr = TypeFloat32
	case reflect.TypeOf(([]float64)(nil)):
		typestr = TypeFloat64
	case reflect.TypeOf(([]uint32)(nil)):
		fallthrough
		// unsupported data type
	case reflect.TypeOf(([]uint64)(nil)):
		fallthrough
		// unsupported data type

	default:
		err = fmt.Errorf("redisai Tensor does not support the following type %v", dtype)
	}
	return
}

func tensorSetFlatArgs(name, dt string, dims []int64, data interface{}) (redis.Args, error) {
	args := redis.Args{}
	var err error = nil
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
			err = fmt.Errorf("redisai.tensorSetFlatArgs: AI.TENSOR does not support the following type %v", reflect.TypeOf(data))
		}
	}
	return args, err
}

func tensorSetInterfaceArgs(keyName string, tensorInterface TensorInterface) (args redis.Args, err error) {
	typestr, err := TensorGetTypeStrFromType(tensorInterface.Dtype())
	if err != nil {
		return
	}
	return tensorSetFlatArgs(keyName, typestr, tensorInterface.Shape(), tensorInterface.Data())
}

func tensorGetParseToInterface(reply interface{}, tensor TensorInterface) (err error) {
	_, shape, data, err := ProcessTensorGetReply(reply, err)
	tensor.SetShape(shape)
	tensor.SetData(data)
	return
}

func ProcessTensorReplyValues(dtype string, reply interface{}) (data interface{}, err error) {
	switch dtype {
	case TypeFloat:
		data, err = converters.Float32s(reply, err)
	case TypeDouble:
		data, err = redis.Float64s(reply, err)
	case TypeInt8:
		data, err = converters.Int8s(reply, err)
	case TypeInt16:
		data, err = converters.Int16s(reply, err)
	case TypeInt32:
		data, err = redis.Ints(reply, err)
	case TypeInt64:
		data, err = redis.Int64s(reply, err)
	case TypeUint8:
		data, err = converters.Uint8s(reply, err)
	case TypeUint16:
		data, err = converters.Uint16s(reply, err)
	}
	return data, err
}

func ProcessTensorGetReply(reply interface{}, errIn error) (dtype string, shape []int64, data interface{}, err error) {
	var replySlice []interface{}
	var key string
	err = errIn
	replySlice, err = redis.Values(reply, err)
	if err != nil {
		return
	}
	for pos := 0; pos < len(replySlice); pos += 2 {
		key, err = redis.String(replySlice[pos], err)
		if err != nil {
			return
		}
		switch key {
		case "dtype":
			dtype, err = redis.String(replySlice[pos+1], err)
			if err != nil {
				return
			}
		case "shape":
			shape, err = redis.Int64s(replySlice[pos+1], err)
			if err != nil {
				return
			}
		case "blob":
			data, err = redis.Bytes(replySlice[pos+1], err)
			if err != nil {
				return
			}
		case "values":
			data, err = ProcessTensorReplyValues(dtype, replySlice[pos+1])
		}
	}
	return
}
