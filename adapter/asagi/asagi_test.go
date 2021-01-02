package asagi

import (
	"testing"
	"time"
)

func TestNYCToUTC(t *testing.T) {
	type args struct {
		timestamp int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "converts NYC Epoch to UTC Epoch",
			args: args{
				timestamp: 0,
			},
			want: time.Unix(0, 0).Add(4 * time.Hour),
		},
		{
			name: "converts 4chan time to UTC",
			args: args{
				timestamp: 1343066185,
			},
			want: time.Unix(1343080585, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NYCToUTC(tt.args.timestamp); got.Unix() != tt.want.Unix() {
				t.Errorf("NYCToUTC() = %v, want %v", got, tt.want)
			}
		})
	}
}
