package main

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertName(tt.args.jsonName); got != tt.want {
				t.Errorf("convertName() = %v, want %v", got, tt.want)
			}
		})
	}
}
