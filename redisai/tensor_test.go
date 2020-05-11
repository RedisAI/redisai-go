package redisai

import (
	"github.com/RedisAI/redisai-go/redisai/converters"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func Test_tensorSetFlatArgs(t *testing.T) {

	f32Bytes, _ := converters.Float32ToByte(1.1)

	type args struct {
		name string
		dt   string
		dims []int64
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test:TestTensorSetArgs:[]float32:1", args{"test:TestTensorSetArgs:1", TypeFloat, []int64{1}, []float32{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]byte:1", args{"test:TestTensorSetArgs:1", TypeFloat, []int64{1}, f32Bytes}, string(TensorContentTypeBlob), false},
		{"test:TestTensorSetArgs:[]int:1", args{"test:TestTensorSetArgs:1", TypeInt32, []int64{1}, []int64{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int8:1", args{"test:TestTensorSetArgs:1", TypeInt8, []int64{1}, []int8{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int16:1", args{"test:TestTensorSetArgs:1", TypeInt16, []int64{1}, []int16{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int64:1", args{"test:TestTensorSetArgs:1", TypeInt64, []int64{1}, []int64{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]uint8:1", args{"test:TestTensorSetArgs:1", TypeUint8, []int64{1}, []uint8{1}}, string(TensorContentTypeBlob), false},
		{"test:TestTensorSetArgs:[]uint16:1", args{"test:TestTensorSetArgs:1", TypeUint16, []int64{1}, []uint16{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]uint32:1", args{"test:TestTensorSetArgs:1", TypeUint8, []int64{1}, []uint32{1}}, string(TensorContentTypeBlob), true},
		{"test:TestTensorSetArgs:[]uint64:1", args{"test:TestTensorSetArgs:1", TypeUint16, []int64{1}, []uint64{1}}, string(TensorContentTypeValues), true},
		{"test:TestTensorSetArgs:[]float32:1", args{"test:TestTensorSetArgs:1", TypeFloat32, []int64{1}, []float32{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]float64:1", args{"test:TestTensorSetArgs:1", TypeFloat64, []int64{1}, []float64{1}}, string(TensorContentTypeValues), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tensorSetFlatArgs(tt.args.name, tt.args.dt, tt.args.dims, tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("tensorSetFlatArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				if got[3].(string) != tt.want {
					t.Errorf("tensorSetFlatArgs() TensorContentType = %v, want %v", got[3], tt.want)
				}
			}
		})
	}
}

func TestTensorGetTypeStrFromType(t *testing.T) {
	type args struct {
		dtype reflect.Type
	}
	tests := []struct {
		name        string
		args        args
		wantTypestr string
		wantErr     bool
	}{
		{"uint8", args{reflect.TypeOf(([]uint8)(nil))}, TypeUint8, false},
		{"uint8", args{reflect.TypeOf(([]byte)(nil))}, TypeUint8, false},
		{"uint8", args{reflect.TypeOf(([]uint16)(nil))}, TypeUint16, false},
		{"uint8", args{reflect.TypeOf(([]int)(nil))}, TypeInt32, false},
		{"uint8", args{reflect.TypeOf(([]int8)(nil))}, TypeInt8, false},
		{"uint8", args{reflect.TypeOf(([]int16)(nil))}, TypeInt16, false},
		{"uint8", args{reflect.TypeOf(([]int32)(nil))}, TypeInt32, false},
		{"uint8", args{reflect.TypeOf(([]int64)(nil))}, TypeInt64, false},
		{"uint8", args{reflect.TypeOf(([]uint8)(nil))}, TypeUint8, false},
		{"uint8", args{reflect.TypeOf(([]uint16)(nil))}, TypeUint16, false},
		{"uint8", args{reflect.TypeOf(([]float32)(nil))}, TypeFloat32, false},
		{"uint8", args{reflect.TypeOf(([]float64)(nil))}, TypeFloat64, false},
		{"uint8", args{reflect.TypeOf(([]string)(nil))}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTypestr, err := TensorGetTypeStrFromType(tt.args.dtype)
			if (err != nil) != tt.wantErr {
				t.Errorf("TensorGetTypeStrFromType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTypestr != tt.wantTypestr {
				t.Errorf("TensorGetTypeStrFromType() gotTypestr = %v, want %v", gotTypestr, tt.wantTypestr)
			}
		})
	}
}

func TestProcessTensorGetReply(t *testing.T) {
	type args struct {
		reply interface{}
		errIn error
	}
	tests := []struct {
		name      string
		args      args
		wantDtype string
		wantShape []int64
		wantData  interface{}
		wantErr   bool
	}{
		{"empty", args{}, "", nil, nil, true},
		{"negative-wrong-reply", args{[]interface{}{[]interface{}{[]byte("serie 1"), []interface{}{}, []interface{}{[]interface{}{[]byte("AA"), []byte("1")}}}}, nil}, "", nil, nil, true},
		{"negative-wrong-reply", args{[]interface{}{[]byte("dtype"), []interface{}{[]byte("dtype"), []byte("1")}}, nil}, "", nil, nil, true},
		{"negative-wrong-shape", args{[]interface{}{[]byte("shape"), []byte("string")}, nil}, "", nil, nil, true},
		{"negative-wrong-blob", args{[]interface{}{[]byte("dtype"), []interface{}{[]byte("dtype"), []byte("1")}}, nil}, "", nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotDtype, gotShape, gotData := ProcessTensorGetReply(tt.args.reply, tt.args.errIn)
			if gotErr != nil && !tt.wantErr {
				t.Errorf("ProcessTensorGetReply() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantDtype, gotDtype); diff != "" {
				t.Errorf("ProcessTensorGetReply() gotDtype mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantShape, gotShape); diff != "" {
				t.Errorf("ProcessTensorGetReply() gotShape mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantData, gotData); diff != "" {
				t.Errorf("ProcessTensorGetReply() gotData mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
