package converters

import (
	"reflect"
	"testing"
)

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
			gotConverted, err := Float32ToByte(tt.args.f)
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
