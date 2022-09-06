package httpx

import (
	"reflect"
	"testing"
)

func TestParseUriQueryToMap(t *testing.T) {
	type args struct {
		query string
	}
	var tests = []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "n1",
			args: args{"name=1&age=2&page=3"},
			want: map[string]interface{}{"name": "1", "age": "2", "page": "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseUriQueryToMap(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseUriQueryToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDaysAgoZeroTime(t *testing.T) {
	type args struct {
		day int
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "n1",
			args: args{day: 2},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDaysAgoZeroTime(tt.args.day); got != tt.want {
				t.Errorf("GetDaysAgoZeroTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeToHuman(t *testing.T) {
	type args struct {
		target int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "n1",
			args: args{target: 1662595200},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeToHuman(tt.args.target); got != tt.want {
				t.Errorf("TimeToHuman() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateDateDir(t *testing.T) {
	type args struct {
		Path string
		prex string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "n1",
			args: args{Path: "/abc", prex: "123_"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateDateDir(tt.args.Path, tt.args.prex); got != tt.want {
				t.Errorf("CreateDateDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
