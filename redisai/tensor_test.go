package redisai
//
//import (
//	"github.com/gomodule/redigo/redis"
//	"reflect"
//	"testing"
//)
//
//func TestTensorSetArgs(t *testing.T) {
//	type args struct {
//		name               string
//		dt                 DataType
//		dims               []int
//		data               interface{}
//		includeCommandName bool
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    redis.Args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := TensorSetArgs(tt.args.name, tt.args.dt, tt.args.dims, tt.args.data, tt.args.includeCommandName)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("TensorSetArgs() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("TensorSetArgs() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_processTensorReplyBlob(t *testing.T) {
//	type args struct {
//		resp []interface{}
//		err  error
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []interface{}
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ProcessTensorReplyBlob(tt.args.resp, tt.args.err)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ProcessTensorReplyBlob() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ProcessTensorReplyBlob() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_processTensorReplyMeta(t *testing.T) {
//	type args struct {
//		resp interface{}
//		err  error
//	}
//	tests := []struct {
//		name     string
//		args     args
//		wantData []interface{}
//		wantErr  bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotData, err := ProcessTensorReplyMeta(tt.args.resp, tt.args.err)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ProcessTensorReplyMeta() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(gotData, tt.wantData) {
//				t.Errorf("ProcessTensorReplyMeta() gotData = %v, want %v", gotData, tt.wantData)
//			}
//		})
//	}
//}
//
//func Test_processTensorReplyValues(t *testing.T) {
//	type args struct {
//		resp []interface{}
//		err  error
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []interface{}
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ProcessTensorReplyValues(tt.args.resp, tt.args.err)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ProcessTensorReplyValues() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ProcessTensorReplyValues() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestClient_TensorGet(t *testing.T) {
//	valuesByteSlice := []byte{1, 2, 3, 4, 5}
//
//	valuesFloat32 := []float32{1.1, 2.2, 3.3, 4.4, 5.5}
//	valuesFloat64 := []float64{1.1}
//
//	valuesInt8 := []int8{1}
//	valuesInt16 := []int16{1}
//	valuesInt32 := []int{1}
//	valuesInt64 := []int64{1}
//
//	valuesUint8 := []uint8{1}
//	valuesUint16 := []uint16{1}
//
//	keyByteSlice := "test:TensorGet:[]byte:1"
//	keyFloat32 := "test:TensorGet:TypeFloat32:1"
//	keyFloat64 := "test:TensorGet:TypeFloat64:1"
//
//	keyInt8 := "test:TensorGet:TypeInt8:1"
//	keyInt16 := "test:TensorGet:TypeInt16:1"
//	keyInt32 := "test:TensorGet:TypeInt32:1"
//	keyInt64 := "test:TensorGet:TypeInt64:1"
//
//	keyUint8 := "test:TensorGet:TypeUint8:1"
//	keyUint16 := "test:TensorGet:TypeUint16:1"
//
//	keyUint8Pipelined := "test:TensorGet:TypeUint8:2:Pipelined"
//	shp := []int{1}
//	shpByteSlice := []int{1, 5}
//	simpleClient := Connect("", createPool())
//	simpleClient.TensorSet(keyByteSlice, TypeUint8, shpByteSlice, valuesByteSlice)
//
//	simpleClient.TensorSet(keyFloat32, TypeFloat32, shpByteSlice, valuesFloat32)
//	simpleClient.TensorSet(keyFloat64, TypeFloat64, shp, valuesFloat64)
//
//	simpleClient.TensorSet(keyInt8, TypeInt8, shp, valuesInt8)
//	simpleClient.TensorSet(keyInt16, TypeInt16, shp, valuesInt16)
//	simpleClient.TensorSet(keyInt32, TypeInt32, shp, valuesInt32)
//	simpleClient.TensorSet(keyInt64, TypeInt64, shp, valuesInt64)
//
//	simpleClient.TensorSet(keyUint8, TypeUint8, shp, valuesUint8)
//	simpleClient.TensorSet(keyUint16, TypeUint16, shp, valuesUint16)
//
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		name string
//		ct   TensorContentType
//	}
//	tests := []struct {
//		name         string
//		fields       fields
//		args         args
//		wantDt       DataType
//		wantShape    []int
//		wantData     interface{}
//		compareDt    bool
//		compareShape bool
//		compareData  bool
//		wantErr      bool
//	}{
//		{keyByteSlice, fields{createPool(), false, 0, 0, nil}, args{keyByteSlice, TensorContentTypeBlob}, TypeUint8, shpByteSlice, valuesByteSlice, true, true, true, false},
//
//		{keyFloat32, fields{createPool(), false, 0, 0, nil}, args{keyFloat32, TensorContentTypeValues}, TypeFloat32, shpByteSlice, valuesFloat32, true, true, true, false},
//		{keyFloat64, fields{createPool(), false, 0, 0, nil}, args{keyFloat64, TensorContentTypeValues}, TypeFloat64, shp, valuesFloat64, true, true, true, false},
//
//		{keyInt8, fields{createPool(), false, 0, 0, nil}, args{keyInt8, TensorContentTypeValues}, TypeInt8, shp, valuesInt8, true, true, true, false},
//		{keyInt16, fields{createPool(), false, 0, 0, nil}, args{keyInt16, TensorContentTypeValues}, TypeInt16, shp, valuesInt16, true, true, true, false},
//		{keyInt32, fields{createPool(), false, 0, 0, nil}, args{keyInt32, TensorContentTypeValues}, TypeInt32, shp, valuesInt32, true, true, true, false},
//		{keyInt64, fields{createPool(), false, 0, 0, nil}, args{keyInt64, TensorContentTypeValues}, TypeInt64, shp, valuesInt64, true, true, true, false},
//
//		{keyUint8, fields{createPool(), false, 0, 0, nil}, args{keyUint8, TensorContentTypeValues}, TypeUint8, shp, valuesUint8, true, true, true, false},
//		{keyUint16, fields{createPool(), false, 0, 0, nil}, args{keyUint16, TensorContentTypeValues}, TypeUint16, shp, valuesUint16, true, true, true, false},
//
//		{keyUint8Pipelined, fields{createPool(), true, 1, 0, nil}, args{keyUint8Pipelined, TensorContentTypeValues}, TypeUint8, shp, valuesUint8, true, true, true, false},
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
//			gotResp, err := c.TensorGet(tt.args.name, tt.args.ct)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("TensorGet() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !tt.fields.PipelineActive {
//				if tt.compareDt && !reflect.DeepEqual(gotResp[0], tt.wantDt) {
//					t.Errorf("TensorGet() gotDt = %v, want %v", gotResp[0], tt.wantDt)
//				}
//				if tt.compareShape && !reflect.DeepEqual(gotResp[1], tt.wantShape) {
//					t.Errorf("TensorGet() gotShape = %v, want %v", gotResp[1], tt.wantShape)
//				}
//				if tt.compareData && !reflect.DeepEqual(gotResp[2], tt.wantData) {
//					t.Errorf("TensorGet() gotData = %v, want %v", gotResp[2], tt.wantData)
//				}
//			}
//		})
//	}
//}
//
//func TestClient_TensorGetBlob(t *testing.T) {
//	valuesByte := []byte{1, 2, 3, 4}
//	keyByte := "test:TensorGetBlog:[]byte:1"
//	keyUnexistant := "test:TensorGetMeta:Unexistant"
//
//	shp := []int{1, 4}
//	simpleClient := Connect("", createPool())
//	simpleClient.TensorSet(keyByte, TypeInt8, shp, valuesByte)
//
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
//		name      string
//		fields    fields
//		args      args
//		wantDt    DataType
//		wantShape []int
//		wantData  []byte
//		wantErr   bool
//	}{
//		{keyByte, fields{createPool(), false, 0, 0, nil}, args{keyByte}, TypeInt8, shp, valuesByte, false},
//		{keyUnexistant, fields{createPool(), false, 0, 0, nil}, args{keyUnexistant}, TypeInt8, shp, valuesByte, true},
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
//			gotDt, gotShape, gotData, err := c.TensorGetBlob(tt.args.name)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("TensorGetBlob() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			if tt.wantErr == false {
//				if gotDt != tt.wantDt {
//					t.Errorf("TensorGetBlob() gotDt = %v, want %v", gotDt, tt.wantDt)
//				}
//				if !reflect.DeepEqual(gotShape, tt.wantShape) {
//					t.Errorf("TensorGetBlob() gotShape = %v, want %v", gotShape, tt.wantShape)
//				}
//				if !reflect.DeepEqual(gotData, tt.wantData) {
//					t.Errorf("TensorGetBlob() gotData = %v, want %v", gotData, tt.wantData)
//				}
//			}
//		})
//	}
//}
//
//func TestClient_TensorGetMeta(t *testing.T) {
//
//	keyFloat32 := "test:TensorGetMeta:TypeFloat32:1"
//	keyFloat64 := "test:TensorGetMeta:TypeFloat64:1"
//
//	keyInt8 := "test:TensorGetMeta:TypeInt8:1"
//	keyInt16 := "test:TensorGetMeta:TypeInt16:1"
//	keyInt32 := "test:TensorGetMeta:TypeInt32:1"
//	keyInt64 := "test:TensorGetMeta:TypeInt64:1"
//
//	keyUint8 := "test:TensorGetMeta:TypeUint8:1"
//	keyUint16 := "test:TensorGetMeta:TypeUint16:1"
//
//	keyUnexistant := "test:TensorGetMeta:Unexistant"
//
//	shp := []int{1, 2}
//	simpleClient := Connect("", createPool())
//	simpleClient.TensorSet(keyFloat32, TypeFloat32, shp, nil)
//	simpleClient.TensorSet(keyFloat64, TypeFloat64, shp, nil)
//
//	simpleClient.TensorSet(keyInt8, TypeInt8, shp, nil)
//	simpleClient.TensorSet(keyInt16, TypeInt16, shp, nil)
//	simpleClient.TensorSet(keyInt32, TypeInt32, shp, nil)
//	simpleClient.TensorSet(keyInt64, TypeInt64, shp, nil)
//
//	simpleClient.TensorSet(keyUint8, TypeUint8, shp, nil)
//	simpleClient.TensorSet(keyUint16, TypeUint16, shp, nil)
//
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
//		name      string
//		fields    fields
//		args      args
//		wantDt    DataType
//		wantShape []int
//		wantErr   bool
//	}{
//		{keyFloat32, fields{createPool(), false, 0, 0, nil}, args{keyFloat32}, TypeFloat32, shp, false},
//		{keyFloat64, fields{createPool(), false, 0, 0, nil}, args{keyFloat64}, TypeFloat64, shp, false},
//		{keyInt8, fields{createPool(), false, 0, 0, nil}, args{keyInt8}, TypeInt8, shp, false},
//		{keyInt16, fields{createPool(), false, 0, 0, nil}, args{keyInt16}, TypeInt16, shp, false},
//		{keyInt32, fields{createPool(), false, 0, 0, nil}, args{keyInt32}, TypeInt32, shp, false},
//		{keyInt64, fields{createPool(), false, 0, 0, nil}, args{keyInt64}, TypeInt64, shp, false},
//		{keyUnexistant, fields{createPool(), false, 0, 0, nil}, args{keyUnexistant}, TypeInt64, shp, true},
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
//			gotDt, gotShape, err := c.TensorGetMeta(tt.args.name)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("TensorGetMeta() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if tt.wantErr == false {
//				if gotDt != tt.wantDt {
//					t.Errorf("TensorGetMeta() gotDt = %v, want %v", gotDt, tt.wantDt)
//				}
//				if !reflect.DeepEqual(gotShape, tt.wantShape) {
//					t.Errorf("TensorGetMeta() gotShape = %v, want %v", gotShape, tt.wantShape)
//				}
//			}
//		})
//	}
//}
//
//func TestClient_TensorGetValues(t *testing.T) {
//
//	valuesFloat32 := []float32{1.1, 2.2, 3.3, 4.4, 5.5}
//	valuesFloat64 := []float64{1.1}
//
//	valuesInt8 := []int8{1}
//	valuesInt16 := []int16{1}
//	valuesInt32 := []int{1}
//	valuesInt64 := []int64{1}
//
//	valuesUint8 := []uint8{1}
//	valuesUint16 := []uint16{1}
//	keyFloat32 := "test:TensorGetValues:TypeFloat32:1"
//	keyFloat64 := "test:TensorGetValues:TypeFloat64:1"
//
//	keyInt8 := "test:TensorGetValues:TypeInt8:1"
//	keyInt16 := "test:TensorGetValues:TypeInt16:1"
//	keyInt32 := "test:TensorGetValues:TypeInt32:1"
//	keyInt64 := "test:TensorGetValues:TypeInt64:1"
//
//	keyUint8 := "test:TensorGetValues:TypeUint8:1"
//	keyUint16 := "test:TensorGetValues:TypeUint16:1"
//	keyUnexistant := "test:TensorGetValues:Unexistant"
//
//	shp := []int{1}
//	shp2 := []int{1, 5}
//	simpleClient := Connect("", createPool())
//	simpleClient.TensorSet(keyFloat32, TypeFloat32, shp2, valuesFloat32)
//	simpleClient.TensorSet(keyFloat64, TypeFloat64, shp, valuesFloat64)
//
//	simpleClient.TensorSet(keyInt8, TypeInt8, shp, valuesInt8)
//	simpleClient.TensorSet(keyInt16, TypeInt16, shp, valuesInt16)
//	simpleClient.TensorSet(keyInt32, TypeInt32, shp, valuesInt32)
//	simpleClient.TensorSet(keyInt64, TypeInt64, shp, valuesInt64)
//
//	simpleClient.TensorSet(keyUint8, TypeUint8, shp, valuesUint8)
//	simpleClient.TensorSet(keyUint16, TypeUint16, shp, valuesUint16)
//
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
//		name      string
//		fields    fields
//		args      args
//		wantDt    DataType
//		wantShape []int
//		wantData  interface{}
//		wantErr   bool
//	}{
//		{keyFloat32, fields{createPool(), false, 0, 0, nil}, args{keyFloat32}, TypeFloat32, shp2, valuesFloat32, false},
//		{keyFloat64, fields{createPool(), false, 0, 0, nil}, args{keyFloat64}, TypeFloat64, shp, valuesFloat64, false},
//
//		{keyInt8, fields{createPool(), false, 0, 0, nil}, args{keyInt8}, TypeInt8, shp, valuesInt8, false},
//		{keyInt16, fields{createPool(), false, 0, 0, nil}, args{keyInt16}, TypeInt16, shp, valuesInt16, false},
//		{keyInt32, fields{createPool(), false, 0, 0, nil}, args{keyInt32}, TypeInt32, shp, valuesInt32, false},
//		{keyInt64, fields{createPool(), false, 0, 0, nil}, args{keyInt64}, TypeInt64, shp, valuesInt64, false},
//
//		{keyUint8, fields{createPool(), false, 0, 0, nil}, args{keyUint8}, TypeUint8, shp, valuesUint8, false},
//		{keyUint16, fields{createPool(), false, 0, 0, nil}, args{keyUint16}, TypeUint16, shp, valuesUint16, false},
//
//		{keyUnexistant, fields{createPool(), false, 0, 0, nil}, args{keyUnexistant}, TypeUint16, shp, valuesUint8, true},
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
//			gotDt, gotShape, gotData, err := c.TensorGetValues(tt.args.name)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("TensorGetValues() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if tt.wantErr == false {
//				if gotDt != tt.wantDt {
//					t.Errorf("TensorGetValues() gotDt = %v, want %v", gotDt, tt.wantDt)
//				}
//				if !reflect.DeepEqual(gotShape, tt.wantShape) {
//					t.Errorf("TensorGetValues() gotShape = %v, want %v", gotShape, tt.wantShape)
//				}
//				if !reflect.DeepEqual(gotData, tt.wantData) {
//					t.Errorf("TensorGetValues() gotData = %v, want %v", gotData, tt.wantData)
//				}
//			}
//		})
//	}
//}
//
//func TestClient_TensorSet(t *testing.T) {
//
//	valuesFloat32 := []float32{1.1}
//	valuesFloat64 := []float64{1.1}
//
//	valuesInt8 := []int8{1}
//	valuesInt16 := []int16{1}
//	valuesInt32 := []int{1}
//	valuesInt64 := []int64{1}
//
//	valuesUint8 := []uint8{1}
//	valuesByte := []byte{1}
//
//	valuesUint16 := []uint16{1}
//	valuesUint32 := []uint32{1}
//	valuesUint64 := []uint64{1}
//
//	keyFloat32 := "test:TensorSet:TypeFloat32:1"
//	keyFloat64 := "test:TensorSet:TypeFloat64:1"
//
//	keyInt8 := "test:TensorSet:TypeInt8:1"
//	keyInt16 := "test:TensorSet:TypeInt16:1"
//	keyInt32 := "test:TensorSet:TypeInt32:1"
//	keyInt64 := "test:TensorSet:TypeInt64:1"
//
//	keyByte := "test:TensorSet:Type[]byte:1"
//	keyUint8 := "test:TensorSet:TypeUint8:1"
//	keyUint16 := "test:TensorSet:TypeUint16:1"
//	keyUint32 := "test:TensorSet:TypeUint32:ExpectError:1"
//	keyUint64 := "test:TensorSet:TypeUint64:ExpectError:1"
//
//	keyInt8Meta := "test:TensorSet:TypeInt8:Meta:1"
//	keyInt8MetaPipelined := "test:TensorSet:TypeInt8:Meta:2:Pipelined"
//
//	keyError1 := "test:TestClient_TensorSet:1:FaultyDims"
//
//	shp := []int{1}
//
//	type fields struct {
//		Pool            *redis.Pool
//		PipelineActive  bool
//		PipelineMaxSize uint32
//		PipelinePos     uint32
//		ActiveConn      redis.Conn
//	}
//	type args struct {
//		name string
//		dt   DataType
//		dims []int
//		data interface{}
//	}
//
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{keyFloat32, fields{createPool(), false, 0, 0, nil}, args{keyFloat32, TypeFloat, shp, valuesFloat32}, false},
//		{keyFloat64, fields{createPool(), false, 0, 0, nil}, args{keyFloat64, TypeFloat64, shp, valuesFloat64}, false},
//
//		{keyInt8, fields{createPool(), false, 0, 0, nil}, args{keyInt8, TypeInt8, shp, valuesInt8}, false},
//		{keyInt16, fields{createPool(), false, 0, 0, nil}, args{keyInt16, TypeInt16, shp, valuesInt16}, false},
//		{keyInt32, fields{createPool(), false, 0, 0, nil}, args{keyInt32, TypeInt32, shp, valuesInt32}, false},
//		{keyInt64, fields{createPool(), false, 0, 0, nil}, args{keyInt64, TypeInt64, shp, valuesInt64}, false},
//
//		{keyUint8, fields{createPool(), false, 0, 0, nil}, args{keyUint8, TypeUint8, shp, valuesUint8}, false},
//		{keyUint16, fields{createPool(), false, 0, 0, nil}, args{keyUint16, TypeUint16, shp, valuesUint16}, false},
//		{keyUint32, fields{createPool(), false, 0, 0, nil}, args{keyUint32, TypeUint8, shp, valuesUint32}, true},
//		{keyUint64, fields{createPool(), false, 0, 0, nil}, args{keyUint64, TypeUint8, shp, valuesUint64}, true},
//
//		{keyInt8Meta, fields{createPool(), false, 0, 0, nil}, args{keyInt8Meta, TypeUint8, shp, nil}, false},
//		{keyInt8MetaPipelined, fields{createPool(), true, 1, 0, nil}, args{keyInt8MetaPipelined, TypeUint8, shp, nil}, false},
//
//		{keyByte, fields{createPool(), false, 0, 0, nil}, args{keyByte, TypeUint8, shp, valuesByte}, false},
//
//		{keyError1, fields{createPool(), false, 0, 0, nil}, args{keyError1, TypeFloat, []int{1, 10}, []float32{1}}, true},
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
//			if err := c.TensorSet(tt.args.name, tt.args.dt, tt.args.dims, tt.args.data); (err != nil) != tt.wantErr {
//				t.Errorf("TensorSet() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
