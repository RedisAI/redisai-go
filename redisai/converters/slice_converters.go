package converters

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

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
		return errNil
	}
	return fmt.Errorf("redisai-go: unexpected type for %s, got type %T", name, reply)
}

func Float32ToByte(f float32) (converted []byte, err error) {
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, f)
	converted = buf.Bytes()
	return
}

var errNil = errors.New("redisai-go: nil returned")

// Float32s is a helper that converts an array command reply to a []float32.
func Float32sBytes(reply interface{}, dimension []int, err error) ([]float32, error) {
	totalResults := 0
	if len(dimension) > 0 {
		totalResults = dimension[0]
	}

	for i := 1; i < len(dimension); i++ {
		totalResults *= dimension[i]
	}

	var result = make([]float32, totalResults)
	tr, err := redis.Bytes(reply, err)
	if err != nil {
		return result, err
	}
	buf := bytes.NewReader(tr)
	for i := 0; i < totalResults; i++ {
		err := binary.Read(buf, binary.LittleEndian, &result[i])
		if err != nil {
			return result, err
		}
	}
	return result, err
}
