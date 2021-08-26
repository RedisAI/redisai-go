package redisai

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_scriptGetParseReply(t *testing.T) {
	type args struct {
		reply interface{}
	}

	tests := []struct {
		name            string
		args            args
		wantDevice      string
		wantTag         string
		wantSource      string
		wantEntryPoints []string
		wantErr         bool
	}{
		{"empty", args{}, "", "", "", nil, true},
		{"negative-wrong-reply", args{[]interface{}{[]interface{}{[]byte("serie 1"), []interface{}{}, []interface{}{[]interface{}{[]byte("AA"), []byte("1")}}}}}, "", "", "", nil, true},
		{"negative-wrong-reply", args{[]interface{}{[]byte("dtype"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, true},
		{"negative-wrong-device", args{[]interface{}{[]byte("device"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, true},
		{"positive-device", args{[]interface{}{[]byte("device"), []byte(DeviceGPU)}}, DeviceGPU, "", "", nil, false},
		{"negative-wrong-tag", args{[]interface{}{[]byte("tag"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, true},
		{"positive-tag", args{[]interface{}{[]byte("tag"), []byte("bar")}}, "", "bar", "", nil, false},
		{"negative-wrong-source", args{[]interface{}{[]byte("source"), []interface{}{[]byte("dtype"), []byte("1")}}}, "", "", "", nil, true},
		{"positive-source", args{[]interface{}{[]byte("source"), []byte("return")}}, "", "", "return", nil, false},
		{"negative-wrong-entry-points", args{[]interface{}{[]byte("Entry Points"), []interface{}{[]interface{}{[]byte("bar"), []byte("foo")}}}}, "", "", "", nil, true},
		{"positive-entry-points", args{[]interface{}{[]byte("Entry Points"), []interface{}{[]byte("bar"), []byte("foo")}}}, "", "", "", []string{"bar", "foo"}, false},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			gotDevice, gotTag, gotSource, gotEntryPoints, gotErr := scriptGetParseReply(tt.args.reply)

			if gotErr != nil && !tt.wantErr {

				t.Errorf("modelGetParseReply() gotErr = %v, want %v", gotErr, tt.wantErr)

			}
			if gotDevice != tt.wantDevice {

				t.Errorf("modelGetParseReply() gotDevice = %v, want %v. gotErr = %v", gotDevice, tt.wantDevice, gotErr)

			}
			if gotTag != tt.wantTag {

				t.Errorf("modelGetParseReply() gotTag = %v, want %v. gotErr = %v", gotTag, tt.wantTag, gotErr)

			}
			if gotSource != tt.wantSource {

				t.Errorf("modelGetParseReply() gotSource = %v, want %v. gotErr = %v", gotSource, tt.wantSource, gotErr)

			}
			assert.EqualValues(t, gotEntryPoints, tt.wantEntryPoints, "modelGetParseReply() gotEntryPoints = %v, want %v. gotErr = %v", gotEntryPoints, tt.wantEntryPoints, gotErr)
		})
	}
}
