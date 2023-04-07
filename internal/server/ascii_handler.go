package server

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/ahmetb/go-cursor"
	"github.com/gabe565/ascii-movie/generated_frames"
	flag "github.com/spf13/pflag"
	"io"
	"time"
)

func New(flags *flag.FlagSet, serveFlags bool) (handler Handler, err error) {
	handler.ClearExtraLines, err = flags.GetInt(ClearExtraLinesFlag)
	if err != nil {
		return handler, err
	}

	handler.Speed, err = flags.GetFloat64(SpeedFlag)
	if err != nil {
		return handler, err
	}
	if handler.Speed <= 0 {
		return handler, fmt.Errorf("%w: %f", ErrInvalidSpeed, handler.Speed)
	}

	if serveFlags {
		handler.SSHConfig.Enabled, err = flags.GetBool(SSHEnabledFlag)
		if err != nil {
			return handler, err
		}
		handler.SSHConfig.Address, err = flags.GetString(SSHAddressFlag)
		if err != nil {
			return handler, err
		}

		handler.TelnetConfig.Enabled, err = flags.GetBool(TelnetEnabledFlag)
		if err != nil {
			return handler, err
		}
		handler.TelnetConfig.Address, err = flags.GetString(TelnetAddressFlag)
		if err != nil {
			return handler, err
		}
	}

	return handler, nil
}

type Handler struct {
	ClearExtraLines int

	Speed float64

	SSHConfig    ServerConfig
	TelnetConfig ServerConfig
}

func (s *Handler) ServeAscii(w io.Writer) error {
	var buf bytes.Buffer
	buf.Grow(generated_frames.Cap)

	for _, f := range generated_frames.List {
		buf.WriteString(f.Data)

		if _, err := io.Copy(w, &buf); err != nil {
			return err
		}

		time.Sleep(f.CalcDuration(s.Speed))

		buf.Reset()
		buf.WriteString(cursor.MoveUp(f.Height+s.ClearExtraLines) + cursor.ClearScreenDown())
	}
	return nil
}

func (s *Handler) MovieDuration() time.Duration {
	var totalDuration time.Duration
	for _, f := range generated_frames.List {
		totalDuration += f.CalcDuration(s.Speed)
	}
	return totalDuration
}
