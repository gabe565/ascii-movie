package movie

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabe565/ascii-movie/movies"
	log "github.com/sirupsen/logrus"
)

func GetInfo(fsys fs.FS, path string) (Info, error) {
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	info := Info{
		Path:    filepath.Clean(path),
		Name:    name,
		Default: path == movies.Default,
	}

	if fsys == nil {
		fsys = os.DirFS(filepath.Dir(path))
		path = filepath.Base(path)
	}
	f, err := fsys.Open(path)
	if err != nil {
		return info, fmt.Errorf("failed to open movie: %w", err)
	}
	defer func(f fs.File) {
		_ = f.Close()
	}(f)

	m, err := NewFromFile(path, f, Padding{}, Padding{})
	if err != nil {
		return info, fmt.Errorf("failed to parse movie: %w", err)
	}
	info.Duration = m.Duration()
	info.NumFrames = len(m.Frames)

	fileInfo, err := f.Stat()
	if err != nil {
		return info, fmt.Errorf("failed to fetch file info: %w", err)
	}
	info.Size = fileInfo.Size()

	return info, nil
}

type Info struct {
	Path      string
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

			info, err := GetInfo(movies.Movies, path)
			if err != nil {
				log.WithError(err).WithField("path", path).Warn("failed to get movie info")
				return nil
			}

			movieInfos = append(movieInfos, info)
			return nil
		},
	); err != nil {
		return nil, err
	}

	return movieInfos, nil
}
