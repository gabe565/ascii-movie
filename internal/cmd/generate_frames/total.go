package main

import (
	"bufio"
	"io"
)

func countNewlines(r io.ReadSeeker) (int, error) {
	scanner := bufio.NewScanner(r)
	var total int
	for scanner.Scan() {
		total += 1
	}
	if err := scanner.Err(); err != nil {
		return total, err
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return total, err
	}

	return total, nil
}
