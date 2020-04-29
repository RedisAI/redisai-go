package redisai
//
//import (
//	"reflect"
//	"testing"
//)
//import "github.com/gomodule/redigo/redis"
//
//func TestClient_ScriptDel(t *testing.T) {
//	keyScript := "test:ScriptDel:1"
//	keyScriptPipelined := "test:ScriptDel:2"
//	keyScriptUnexistant := "test:ScriptDel:3:Unexistant"
//	scriptBin := "def bar(a, b):\n    return a + b\n"
//	simpleClient := Connect("", createPool())
//	err := simpleClient.ScriptSet(keyScript, DeviceCPU, scriptBin)
//	if err != nil {
//		t.Errorf("Error preparing for ScriptDel(), while issuing ScriptSet. error = %v", err)
//		return
//	}
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		name string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{keyScript, fields{createPool(), false, 0, 0, nil}, args{keyScript}, false},
//		{keyScriptPipelined, fields{createPool(), true, 1, 0, nil}, args{keyScript}, false},
//		{keyScriptUnexistant, fields{createPool(), false, 0, 0, nil}, args{keyScriptUnexistant}, true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			if err := c.ScriptDel(tt.args.name); (err != nil) != tt.wantErr {
//				t.Errorf("ScriptDel() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestClient_ScriptGet(t *testing.T) {
//	keyScript := "test:ScriptGet:1"
//	keyScriptPipelined := "test:ScriptGet:2"
//	keyScriptEmpty := "test:ScriptGet:3:Empty"
//	scriptBin := ""
//	simpleClient := Connect("", createPool())
//	err := simpleClient.ScriptSet(keyScript, DeviceCPU, scriptBin)
//	if err != nil {
//		t.Errorf("Error preparing for ScriptGet(), while issuing ScriptSet. error = %v", err)
//		return
//	}
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		name string
//	}
//	tests := []struct {
//		name           string
//		fields         fields
//		args           args
//		wantDeviceType DeviceType
//		wantData       string
//		wantErr        bool
//	}{
//		{keyScript, fields{createPool(), false, 0, 0, nil}, args{keyScript}, DeviceCPU, "", false},
//		{keyScriptPipelined, fields{createPool(), true, 1, 0, nil}, args{keyScript}, DeviceCPU, "", false},
//		{keyScriptEmpty, fields{createPool(), false, 0, 0, nil}, args{keyScriptEmpty}, DeviceCPU, "", true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			gotData, err := c.ScriptGet(tt.args.name)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ScriptGet() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if tt.wantErr == false && !tt.fields.PipelineActive {
//				if !reflect.DeepEqual(gotData["device"], tt.wantDeviceType) {
//					t.Errorf("ScriptGet() gotData = %v, want %v", gotData["device"], tt.wantDeviceType)
//				}
//				if !reflect.DeepEqual(gotData["source"], tt.wantData) {
//					t.Errorf("ScriptGet() gotData = %v, want %v", gotData["source"], tt.wantData)
//				}
//			}
//
//		})
//	}
//}
//
//func TestClient_ScriptRun(t *testing.T) {
//	keyScript1 := "test:ScriptRun:1"
//	keyScript2 := "test:ScriptRun:2:Pipelined"
//	keyScript3Empty := "test:ScriptRun:3:Empty"
//	scriptBin := "def bar(a, b):\n    return a + b\n"
//	simpleClient := Connect("", createPool())
//	err := simpleClient.ScriptSet(keyScript1, DeviceCPU, scriptBin)
//	if err != nil {
//		t.Errorf("Error preparing for ScriptRun(), while issuing ScriptSet. error = %v", err)
//		return
//	}
//	err = simpleClient.ScriptSet(keyScript2, DeviceCPU, scriptBin)
//	if err != nil {
//		t.Errorf("Error preparing for ScriptRun(), while issuing ScriptSet. error = %v", err)
//		return
//	}
//
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		name    string
//		fn      string
//		inputs  []string
//		outputs []string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{keyScript2, fields{createPool(), true, 1, 0, nil}, args{keyScript2, "", []string{""}, []string{""}}, false},
//		{keyScript3Empty, fields{createPool(), false, 0, 0, nil}, args{keyScript3Empty, "", []string{""}, []string{""}}, true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			if err := c.ScriptRun(tt.args.name, tt.args.fn, tt.args.inputs, tt.args.outputs); (err != nil) != tt.wantErr {
//				t.Errorf("ScriptRun() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestClient_ScriptSet(t *testing.T) {
//	keyScriptError := "test:ScriptSet:Error:1"
//	scriptBin := "import abc"
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		name   string
//		device DeviceType
//		data   string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{keyScriptError, fields{createPool(), false, 0, 0, nil}, args{keyScriptError, DeviceCPU, scriptBin}, true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			if err := c.ScriptSet(tt.args.name, tt.args.device, tt.args.data); (err != nil) != tt.wantErr {
//				t.Errorf("ScriptSet() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestClient_ScriptSetFromFile(t *testing.T) {
//	keyScript1 := "test:ScriptSetFromFile:1:DontExist"
//	keyScript2 := "test:ScriptSetFromFile:2"
//	keyScript3Pipelined := "test:ScriptSetFromFile:3"
//
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		name   string
//		device DeviceType
//		path   string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{keyScript1, fields{createPool(), false, 0, 0, nil}, args{keyScript1, DeviceCPU, "./../tests/testdata/dontexist"}, true},
//		{keyScript2, fields{createPool(), false, 0, 0, nil}, args{keyScript2, DeviceCPU, "./../tests/testdata/script.txt"}, false},
//		{keyScript3Pipelined, fields{createPool(), true, 1, 0, nil}, args{keyScript3Pipelined, DeviceCPU, "./../tests/testdata/script.txt"}, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				Pool:                  tt.fields.Pool,
//				PipelineActive:        tt.fields.PipelineActive,
//				PipelineAutoFlushSize: tt.fields.PipelineMaxSize,
//				PipelinePos:           tt.fields.PipelinePos,
//				ActiveConn:            tt.fields.ActiveConn,
//			}
//			if err := c.ScriptSetFromFile(tt.args.name, tt.args.device, tt.args.path); (err != nil) != tt.wantErr {
//				t.Errorf("ScriptSetFromFile() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
