package domain

import (
	"errors"
	"sync/atomic"
)

var ErrKeyNotFound = errors.New("key not found")

type Meta[T comparable] struct {
	Names  [7]T
	Values [7]atomic.Uint32
}

func (m *Meta[T]) GetValue(byKey T) (uint32, error) {
	for ix, n := range m.Names {
		if n == byKey {
			return m.Values[ix].Load(), nil
		}
	}

	return 0,
		ErrKeyNotFound
}

func (m *Meta[T]) Increment(key T) {
	// Scan for exact match
	for i, name := range m.Names {
		if name == key {
			m.Values[i].Add(1)

			return
		}
	}

	// Look for an empty slot
	var zeroT T

	for ix, name := range m.Names {
		if name == zeroT {
			m.Names[ix] = key
			m.Values[ix].Store(1)

			return
		}
	}

	// Matrix is full: Find the absolute weakest link
	lowestIdx := 0
	lowestVal := m.Values[0].Load()

	for ix := 1; ix < 7; ix++ {
		if m.Values[ix].Load() < lowestVal {
			lowestVal = m.Values[ix].Load()
			lowestIdx = ix
		}
	}

	// Space-Saving Eviction
	// We evict the old loser, install the new key, and increment the counter.
	// This ensures a massive spike will instantly overtake the matrix.
	m.Names[lowestIdx] = key
	m.Values[lowestIdx].Store(lowestVal + 1)
}
