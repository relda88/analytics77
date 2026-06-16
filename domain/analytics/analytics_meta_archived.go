package analytics

type MetaArchived[T comparable] struct {
	Names  [7]T
	Values [7]uint32
}

func (m *MetaArchived[T]) GetValue(byKey T) (uint32, error) {
	for ix, n := range m.Names {
		if n == byKey {
			return m.Values[ix], nil
		}
	}

	return 0,
		ErrKeyNotFound
}
