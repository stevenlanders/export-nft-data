package utils

import "math"

func WithPages(start uint64, end uint64, size uint64, handler func(start, end uint64) error) error {
	for pageStart := start; pageStart < end; pageStart += size {
		pageEnd := uint64(math.Min(float64(end), float64(pageStart+size-1)))
		if err := handler(pageStart, pageEnd); err != nil {
			return err
		}
	}
	return nil
}
