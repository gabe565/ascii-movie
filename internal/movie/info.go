package movie

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabe565/ascii-movie/movies"
	log "github.com/sirupsen/logrus"
)

type Info struct {
	Name      string
	Duration  time.Duration
	Default   bool
	NumFrames int
	Size      int64
}

func ListEmbedded() ([]Info, error) {
	var movieInfos []Info

	if err := fs.WalkDir(
		movies.Movies,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			movieLog := log.WithField("path", path)

			f, err := movies.Movies.Open(path)
			if err != nil {
				movieLog.WithError(err).Warn("Failed to open movie")
				return nil
			}

			m, err := NewFromFile(path, f, Padding{}, Padding{})
			if err != nil {
				movieLog.WithError(err).Warn("Failed to parse movie")
			}

			info, err := d.Info()
			if err != nil {
				log.WithError(err).Warn("Failed to fetch file info")
			}

			movieInfos = append(movieInfos, Info{
				Name:      strings.TrimSuffix(path, filepath.Ext(path)),
				Duration:  m.Duration(),
				Default:   path == movies.Default,
				NumFrames: len(m.Frames),
				Size:      info.Size(),
			})
			return nil
		},
	); err != nil {
		return nil, err
	}

	return movieInfos, nil
}
