package progressbar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProgressBar_Generate(t *testing.T) {
	type fields struct {
		Phases []string
	}
	type args struct {
		n     time.Duration
		total time.Duration
		width int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"0%", fields{Phases: DefaultPhases}, args{
			n: 0, total: 1000,
			width: 50,
		}, "                                                "},
		{"100%", fields{Phases: DefaultPhases}, args{
			n: 1000, total: 1000,
			width: 50,
		}, "████████████████████████████████████████████████"},
		{"82.3%", fields{Phases: DefaultPhases}, args{
			n: 823, total: 1000,
			width: 50,
		}, "███████████████████████████████████████▌        "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProgressBar{
				Phases: tt.fields.Phases,
			}
			assert.Equal(t, tt.want, p.Generate(tt.args.n, tt.args.total, tt.args.width))
		})
	}
}
