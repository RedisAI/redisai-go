package redisai

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_modelGetParseReply(t *testing.T) {
	type args struct {
		reply interface{}
	}
	tests := []struct {
		name                string
		args                args
		wantBackend         string
		wantDevice          string
		wantTag             string
		wantBlob            []byte
		wantBatchsize       int64
		wantMinbatchsize    int64
		wantMinbatchtimeout int64
		wantInputs          []string
		wantOutputs         []string
		wantErr             bool
	}{
		{"empty", args{}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"negative-wrong-reply", args{[]interface{}{[]interface{}{[]byte("serie 1"), []interface{}{}, []interface{}{[]interface{}{[]byte("AA"), []byte("1")}}}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"negative-wrong-reply", args{[]interface{}{[]byte("dtype"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"positive-backend", args{[]interface{}{[]byte("backend"), []byte(BackendTF)}}, BackendTF, "", "", nil, 0, 0, 0, nil, nil, false},
		{"negative-wrong-device", args{[]interface{}{[]byte("device"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"positive-device", args{[]interface{}{[]byte("device"), []byte(DeviceGPU)}}, "", DeviceGPU, "", nil, 0, 0, 0, nil, nil, false},
		{"negative-wrong-batchsize", args{[]interface{}{[]byte("batchsize"), []interface{}{[]byte("1")}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"positive-batchsize", args{[]interface{}{[]byte("batchsize"), int64(1)}}, "", "", "", nil, 1, 0, 0, nil, nil, false},
		{"negative-wrong-minbatchsize", args{[]interface{}{[]byte("minbatchsize"), []interface{}{[]byte("1")}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"positive-minbatchsize", args{[]interface{}{[]byte("minbatchsize"), int64(1)}}, "", "", "", nil, 0, 1, 0, nil, nil, false},
		{"negative-wrong-minbatchtimeout", args{[]interface{}{[]byte("minbatchtimeout"), []interface{}{[]byte("1")}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"positive-minbatchtimeout", args{[]interface{}{[]byte("minbatchtimeout"), int64(1)}}, "", "", "", nil, 0, 0, 1, nil, nil, false},
		{"negative-wrong-inputs", args{[]interface{}{[]byte("inputs"), []interface{}{[]interface{}{[]byte("bar"), []byte("foo")}}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"positive-inputs", args{[]interface{}{[]byte("inputs"), []interface{}{[]byte("bar"), []byte("foo")}}}, "", "", "", nil, 0, 0, 0, []string{"bar", "foo"}, nil, false},
		{"negative-wrong-output", args{[]interface{}{[]byte("output"), []interface{}{[]interface{}{[]byte("output")}}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"positive-output", args{[]interface{}{[]byte("outputs"), []interface{}{[]byte("output")}}}, "", "", "", nil, 0, 0, 0, nil, []string{"output"}, false},
		{"negative-wrong-blob", args{[]interface{}{[]byte("blob"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, 0, 0, 0, nil, nil, true},
		{"positive-blob", args{[]interface{}{[]byte("blob"), []byte("blob")}}, "", "", "", []byte("blob"), 0, 0, 0, nil, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBackend, gotDevice, gotTag, gotBlob, gotBatchsize, gotMinbatchsize, gotInputs, gotOutputs, gotMinbatchtimeout, gotErr := modelGetParseReply(tt.args.reply)
			if gotErr != nil && !tt.wantErr {
				t.Errorf("modelGetParseReply() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
			if gotBackend != tt.wantBackend {
				t.Errorf("modelGetParseReply() gotBackend = %v, want %v. gotErr = %v", gotBackend, tt.wantBackend, gotErr)
			}
			if gotDevice != tt.wantDevice {
				t.Errorf("modelGetParseReply() gotDevice = %v, want %v. gotErr = %v", gotDevice, tt.wantDevice, gotErr)
			}
			if gotTag != tt.wantTag {
				t.Errorf("modelGetParseReply() gotTag = %v, want %v. gotErr = %v", gotTag, tt.wantTag, gotErr)
			}
			if gotBatchsize != tt.wantBatchsize {
				t.Errorf("modelGetParseReply() gotBatchsize = %v, want %v. gotErr = %v", gotBatchsize, tt.wantBatchsize, gotErr)
			}
			if gotMinbatchsize != tt.wantMinbatchsize {
				t.Errorf("modelGetParseReply() gotMinbatchsize = %v, want %v. gotErr = %v", gotMinbatchsize, tt.wantMinbatchsize, gotErr)
			}
			if gotMinbatchtimeout != tt.wantMinbatchtimeout {
				t.Errorf("modelGetParseReply() gotMinbatchsize = %v, want %v. gotErr = %v", gotMinbatchsize, tt.wantMinbatchsize, gotErr)
			}
			assert.EqualValues(t, gotInputs, tt.wantInputs, "modelGetParseReply() gotInputs = %v, want %v. gotErr = %v", gotInputs, tt.wantInputs, gotErr)
			assert.EqualValues(t, gotOutputs, tt.wantOutputs, "modelGetParseReply() gotOutputs = %v, want %v. gotErr = %v", gotOutputs, tt.wantOutputs, gotErr)
			if !reflect.DeepEqual(gotBlob, tt.wantBlob) {
				t.Errorf("modelGetParseReply() gotBlob = %v, want %v. gotErr = %v", gotBlob, tt.wantBlob, gotErr)
			}
		})
	}
}

func Test_modelGetParseToInterface(t *testing.T) {
	type args struct {
		reply interface{}
		model ModelInterface
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"negative-wrong-reply", args{[]interface{}{[]interface{}{[]byte("serie 1"), []interface{}{}, []interface{}{[]interface{}{[]byte("AA"), []byte("1")}}}}, nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := modelGetParseToInterface(tt.args.reply, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("modelGetParseToInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
