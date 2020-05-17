package redisai

import (
	"reflect"
	"testing"
)

func Test_modelGetParseReply(t *testing.T) {
	type args struct {
		reply interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantBackend string
		wantDevice  string
		wantTag     string
		wantBlob    []byte
		wantErr     bool
	}{
		{"empty", args{}, "", "", "", nil, true},
		{"negative-wrong-reply", args{[]interface{}{[]interface{}{[]byte("serie 1"), []interface{}{}, []interface{}{[]interface{}{[]byte("AA"), []byte("1")}}}}}, "", "", "", nil, true},
		{"negative-wrong-reply", args{[]interface{}{[]byte("dtype"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, true},
		{"negative-wrong-device", args{[]interface{}{[]byte("device"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, true},
		{"negative-wrong-blob", args{[]interface{}{[]byte("blob"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotBackend, gotDevice, gotTag, gotBlob := modelGetParseReply(tt.args.reply)
			if gotErr != nil && !tt.wantErr {
				t.Errorf("modelGetParseReply() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
			if gotBackend != tt.wantBackend {
				t.Errorf("modelGetParseReply() gotBackend = %v, want %v", gotBackend, tt.wantBackend)
			}
			if gotDevice != tt.wantDevice {
				t.Errorf("modelGetParseReply() gotDevice = %v, want %v", gotDevice, tt.wantDevice)
			}
			if gotTag != tt.wantTag {
				t.Errorf("modelGetParseReply() gotTag = %v, want %v", gotTag, tt.wantTag)
			}
			if !reflect.DeepEqual(gotBlob, tt.wantBlob) {
				t.Errorf("modelGetParseReply() gotBlob = %v, want %v", gotBlob, tt.wantBlob)
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
