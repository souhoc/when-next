package discord_test

import (
	"testing"
	"time"

	"github.com/souhoc/when-next/discord"
)

func TestSnowflake_Unix(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		sf   discord.Snowflake
		want int64
	}{
		{
			name: "Given ex",
			sf:   175928847299117063,
			want: 1462015105796,
		},
		{
			name: "My ID",
			sf:   212581406344216578,
			want: 1470753756844,
		},
	}
	for _, tt := range tests {
		time.Now().Unix()
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sf.Unix()
			if got != tt.want {
				t.Errorf("Unix() = %v, want %v", got, tt.want)
			}
		})
	}
}
