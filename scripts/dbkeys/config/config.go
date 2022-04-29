package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/tabwriter"

	tmdb "github.com/tendermint/tm-db"
)

// Settings carries configuration values to subcommands.
type Settings struct {
	Spec string // database spec: backend:path
}

// WithDB opens the database associated with s, runs f with it, and then closes
// it. The error reported by f is returned by WithDB.
func (s *Settings) WithDB(f func(tmdb.DB) error) error {
	backend, path, err := ParseDatabaseSpec(s.Spec)
	if err != nil {
		return err
	}

	// tm-db expects the target directory and the path to be separate.
	// Default to the current working directory if one is not specified.
	dir, base := filepath.Split(path)
	if dir == "" {
		dir = "."
	}
	base = strings.TrimSuffix(base, ".db")

	db, err := tmdb.NewDB(base, tmdb.BackendType(backend), dir)
	if err != nil {
		return fmt.Errorf("opening %q: %w", s.Spec, err)
	}
	defer db.Close()
	return f(db)
}

// ParseDatabaseSpec parses spec as a backend:path specification.
func ParseDatabaseSpec(spec string) (backend, path string, _ error) {
	parts := strings.SplitN(spec, ":", 2)
	if len(parts) == 1 {
		return "", "", errors.New("invalid backend:path spec")
	} else if parts[0] == "" || parts[1] == "" {
		return "", "", errors.New("backend and path must be non-empty")
	}
	return parts[0], parts[1], nil
}

type Histogram struct {
	Samples int
	Sum     int64
	Max     int64
	Base    int
	Buckets []int64
}

func (h *Histogram) AddSample(size int) {
	h.Samples++
	z := int64(size)
	h.Sum += z
	if z > h.Max {
		h.Max = z
	}
	if h.Base <= 0 {
		h.Base = 4
	}

	v := size
	for i := 0; ; i++ {
		if len(h.Buckets) <= i {
			h.Buckets = append(h.Buckets, 0)
		}

		w := v / h.Base
		if w == 0 {
			h.Buckets[i]++
			return
		}
		v = w
	}
}

func (h *Histogram) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 1, 8, 2, ' ', 0)
	fmt.Fprintf(tw, "n = %d\n", h.Samples)
	fmt.Fprintf(tw, "max = %d\n", h.Max)
	fmt.Fprintf(tw, "avg = %.1f\n", float64(h.Sum)/float64(h.Samples))
	fmt.Fprintf(tw, "[%d^i]\n", h.Base)
	for i, v := range h.Buckets {
		incr := int(20 * float64(v) / float64(h.Samples))
		fmt.Fprintf(tw, "%-2d\t%-6d\t%s\n", i, v, strings.Repeat("=", incr))
	}
	tw.Flush()
	return io.Copy(w, &buf)
}
