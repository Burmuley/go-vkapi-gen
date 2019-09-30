/*
Copyright 2019 Konstantin Vasilev (burmuley@gmail.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"reflect"
	"testing"
)

func Test_convertName(t *testing.T) {
	type args struct {
		jsonName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"TestUnderscores",
			args{"widgets_getPages_response"},
			"WidgetsGetPages",
		},
		{
			"TestNumberFirstByte",
			args{"2fa_required"},
			"TwofaRequired",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertName(tt.args.jsonName); got != tt.want {
				t.Errorf("convertName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkErr(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_convertName1(t *testing.T) {
	type args struct {
		jsonName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertName(tt.args.jsonName); got != tt.want {
				t.Errorf("convertName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getApiNamePrefix(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getApiNamePrefix(tt.args.name); got != tt.want {
				t.Errorf("getApiNamePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getObjectTypeName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getObjectTypeName(tt.args.s); got != tt.want {
				t.Errorf("getObjectTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readLocalSchemaFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLocalSchemaFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readLocalSchemaFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readLocalSchemaFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readSchemaFile(t *testing.T) {
	type args struct {
		fileUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readHTTPSchemaFile(tt.args.fileUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("readHTTPSchemaFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readHTTPSchemaFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
