package redisai

import (
	"github.com/RedisAI/redisai-go/redisai/converters"
	"testing"
)

func Test_tensorSetFlatArgs(t *testing.T) {

	f32Bytes, _ := converters.Float32ToByte(1.1)

	type args struct {
		name string
		dt   string
		dims []int
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test:TestTensorSetArgs:[]float32:1", args{"test:TestTensorSetArgs:1", TypeFloat, []int{1}, []float32{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]byte:1", args{"test:TestTensorSetArgs:1", TypeFloat, []int{1}, f32Bytes}, string(TensorContentTypeBlob), false},
		{"test:TestTensorSetArgs:[]int:1", args{"test:TestTensorSetArgs:1", TypeInt32, []int{1}, []int{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int8:1", args{"test:TestTensorSetArgs:1", TypeInt8, []int{1}, []int8{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int16:1", args{"test:TestTensorSetArgs:1", TypeInt16, []int{1}, []int16{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]int64:1", args{"test:TestTensorSetArgs:1", TypeInt64, []int{1}, []int64{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]uint8:1", args{"test:TestTensorSetArgs:1", TypeUint8, []int{1}, []uint8{1}}, string(TensorContentTypeBlob), false},
		{"test:TestTensorSetArgs:[]uint16:1", args{"test:TestTensorSetArgs:1", TypeUint16, []int{1}, []uint16{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]uint32:1", args{"test:TestTensorSetArgs:1", TypeUint8, []int{1}, []uint32{1}}, string(TensorContentTypeBlob), true},
		{"test:TestTensorSetArgs:[]uint64:1", args{"test:TestTensorSetArgs:1", TypeUint16, []int{1}, []uint64{1}}, string(TensorContentTypeValues), true},
		{"test:TestTensorSetArgs:[]float32:1", args{"test:TestTensorSetArgs:1", TypeFloat32, []int{1}, []float32{1}}, string(TensorContentTypeValues), false},
		{"test:TestTensorSetArgs:[]float64:1", args{"test:TestTensorSetArgs:1", TypeFloat64, []int{1}, []float64{1}}, string(TensorContentTypeValues), false},
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
