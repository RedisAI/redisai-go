package redisai

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"testing"
)

func TestModelRunArgs(t *testing.T) {
	nameT1 := "test:ModelRunArgs:1:includeCommandName"
	type args struct {
		name               string
		inputs             []string
		outputs            []string
		includeCommandName bool
	}
	tests := []struct {
		name string
		args args
		want redis.Args
	}{
		{nameT1, args{nameT1, []string{}, []string{}, true}, redis.Args{"AI.MODELRUN", nameT1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ModelRunArgs(tt.args.name, tt.args.inputs, tt.args.outputs, tt.args.includeCommandName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModelRunArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTensorSetArgs_TensorContentType(t *testing.T) {

	f32Bytes, _ := float32ToByte(1.1)

	type args struct {
		name               string
		dt                 DataType
		dims               []int
		data               interface{}
		includeCommandName bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test:TestTensorSetArgs:[]float32:1", args{"test:TestTensorSetArgs:1", TypeFloat, []int{1}, []float32{1}, true}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]byte:1", args{"test:TestTensorSetArgs:1", TypeFloat, []int{1}, f32Bytes, true}, string(TensorContentTypeBlob), false},
		{"test:TestTensorSetArgs:[]int:1", args{"test:TestTensorSetArgs:1", TypeInt32, []int{1}, []int{1}, true}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int8:1", args{"test:TestTensorSetArgs:1", TypeInt8, []int{1}, []int8{1}, true}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int16:1", args{"test:TestTensorSetArgs:1", TypeInt16, []int{1}, []int16{1}, true}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int64:1", args{"test:TestTensorSetArgs:1", TypeInt64, []int{1}, []int64{1}, true}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]uint8:1", args{"test:TestTensorSetArgs:1", TypeUint8, []int{1}, []uint8{1}, true}, string(TensorContentTypeBlob), false},
		{"test:TestTensorSetArgs:[]uint16:1", args{"test:TestTensorSetArgs:1", TypeUint16, []int{1}, []uint16{1}, true}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]uint32:1", args{"test:TestTensorSetArgs:1", TypeUint8, []int{1}, []uint32{1}, true}, string(TensorContentTypeBlob), true},
		{"test:TestTensorSetArgs:[]uint64:1", args{"test:TestTensorSetArgs:1", TypeUint16, []int{1}, []uint64{1}, true}, string(TensorContentTypeValues), true},
		{"test:TestTensorSetArgs:[]float32:1", args{"test:TestTensorSetArgs:1", TypeFloat32, []int{1}, []float32{1}, true}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]float64:1", args{"test:TestTensorSetArgs:1", TypeFloat64, []int{1}, []float64{1}, true}, string(TensorContentTypeValues), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TensorSetArgs(tt.args.name, tt.args.dt, tt.args.dims, tt.args.data, tt.args.includeCommandName)

			if (err != nil) != tt.wantErr {
				t.Errorf("TensorSetArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				if got[4].(string) != tt.want {
					t.Errorf("TensorSetArgs() TensorContentType = %v, want %v", got[4], tt.want)
				}
			}
		})
	}
}

func Test_replyDataType(t *testing.T) {

	var r1 interface{} = string("abc")
	var r2 interface{} = int(1)
	var r3 interface{} = string("FLOAT")
	var r4 interface{} = string("DOUBLE")
	var r5 interface{} = string("INT8")
	var r6 interface{} = string("INT16")
	var r7 interface{} = string("INT32")
	var r8 interface{} = string("INT64")
	var r9 interface{} = string("UINT8")
	var r10 interface{} = string("UINT16")
	var r11 interface{} = nil

	var err1 error = fmt.Errorf("")

	type args struct {
		reply interface{}
		err   error
	}
	tests := []struct {
		name    string
		args    args
		wantDt  DataType
		wantErr bool
	}{
		{"test:replyDataType:Error:1", args{r1, err1}, "", true},
		{"test:replyDataType:Error:WrongType:2", args{r2, nil}, "", true},
		{"test:replyDataType:FLOAT:3", args{r3, nil}, TypeFloat, false},
		{"test:replyDataType:DOUBLE:4", args{r4, nil}, TypeDouble, false},
		{"test:replyDataType:INT8:5", args{r5, nil}, TypeInt8, false},
		{"test:replyDataType:INT16:6", args{r6, nil}, TypeInt16, false},
		{"test:replyDataType:INT32:7", args{r7, nil}, TypeInt32, false},
		{"test:replyDataType:INT64:8", args{r8, nil}, TypeInt64, false},
		{"test:replyDataType:UINT8:9", args{r9, nil}, TypeUint8, false},
		{"test:replyDataType:UINT16:10", args{r10, nil}, TypeUint16, false},
		{"test:replyDataType:11:nil", args{r11, nil}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDt, err := replyDataType(tt.args.reply, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("replyDataType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDt != tt.wantDt {
				t.Errorf("replyDataType() gotDt = %v, want %v", gotDt, tt.wantDt)
			}
		})
	}
}

func TestFloat32s(t *testing.T) {

	var r1 interface{} = nil

	type args struct {
		reply interface{}
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    []float32
		wantErr bool
	}{
		{"test:Float32s:1", args{r1, nil}, []float32{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Float32s(tt.args.reply, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Float32s() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false {

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Float32s() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestInt16s(t *testing.T) {

	var r1 interface{} = nil

	type args struct {
		reply interface{}
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    []int16
		wantErr bool
	}{
		{"test:Int16s:1", args{r1, nil}, []int16{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Int16s(tt.args.reply, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Int16s() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false {

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Int16s() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestInt8s(t *testing.T) {

	var r1 interface{} = nil

	type args struct {
		reply interface{}
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    []int8
		wantErr bool
	}{
		{"test:Int8s:1", args{r1, nil}, []int8{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Int8s(tt.args.reply, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Int8s() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false {

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Int8s() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTensorSetArgs(t *testing.T) {
	type args struct {
		name               string
		dt                 DataType
		dims               []int
		data               interface{}
		includeCommandName bool
	}
	tests := []struct {
		name    string
		args    args
		want    redis.Args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TensorSetArgs(tt.args.name, tt.args.dt, tt.args.dims, tt.args.data, tt.args.includeCommandName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TensorSetArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TensorSetArgs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint16s(t *testing.T) {

	var r1 interface{} = nil

	type args struct {
		reply interface{}
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    []uint16
		wantErr bool
	}{
		{"test:Uint16s:1", args{r1, nil}, []uint16{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Uint16s(tt.args.reply, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint16s() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Uint16s() got = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

func TestUint8s(t *testing.T) {

	var r1 interface{} = nil

	type args struct {
		reply interface{}
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    []uint8
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test:Uint8s:1", args{r1, nil}, []uint8{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Uint8s(tt.args.reply, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint8s() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false {

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Uint8s() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_float32ToByte(t *testing.T) {
	type args struct {
		f float32
	}
	tests := []struct {
		name          string
		args          args
		wantConverted []byte
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConverted, err := float32ToByte(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("float32ToByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConverted, tt.wantConverted) {
				t.Errorf("float32ToByte() gotConverted = %v, want %v", gotConverted, tt.wantConverted)
			}
		})
	}
}

func Test_processTensorReplyBlob(t *testing.T) {
	type args struct {
		resp []interface{}
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processTensorReplyBlob(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTensorReplyBlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processTensorReplyBlob() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processTensorReplyMeta(t *testing.T) {
	type args struct {
		resp interface{}
		err  error
	}
	tests := []struct {
		name     string
		args     args
		wantData []interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := processTensorReplyMeta(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTensorReplyMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("processTensorReplyMeta() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func Test_processTensorReplyValues(t *testing.T) {
	type args struct {
		resp []interface{}
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processTensorReplyValues(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTensorReplyValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processTensorReplyValues() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sliceHelper(t *testing.T) {
	type args struct {
		reply     interface{}
		err       error
		name      string
		makeSlice func(int)
		assign    func(int, interface{}) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sliceHelper(tt.args.reply, tt.args.err, tt.args.name, tt.args.makeSlice, tt.args.assign); (err != nil) != tt.wantErr {
				t.Errorf("sliceHelper() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
