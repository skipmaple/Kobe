package gitcode

import "testing"

func Test_currentWeather(t *testing.T) {
	type args struct {
		city string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		// TODO: Add test cases.
		{"t1", args{city: "Beijing"}, " 现在的天气是 -3.0 °C, 晴天"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := currentWeather(tt.args.city); gotRes != tt.wantRes {
				t.Errorf("currentWeather() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
