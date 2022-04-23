package set

type HashSet[K, V comparable] struct {
	Data map[K]V
}

func NewSet[K comparable](keys ...K) *HashSet[K, struct{}] {
	set := &HashSet[K, struct{}]{
		Data: make(map[K]struct{}),
	}

	for _, k := range keys {
		set.Add(k, struct{}{})
	}

	return set
}

func (s *HashSet[K, V]) Add(k K, v V) {
	s.Data[k] = v
}

func (s *HashSet[K, V]) Has(k K) bool {
	_, ok := s.Data[k]
	return ok
}

func (s *HashSet[K, V]) Pop(k K) {
	delete(s.Data, k)
}

func (s *HashSet[K, V]) Len() int {
	return len(s.Data)
}

func (s *HashSet[K, V]) Keys() []K {
	out := make([]K, 0, s.Len())
	for k := range s.Data {
		out = append(out, k)
	}

	return out
}

func (s *HashSet[K, V]) Maps() map[K]V {
	return s.Data
}

func (s *HashSet[K, V]) KeysWithValue(value V) []K {
	out := make([]K, 0, s.Len())
	for k := range s.Data {
		if v, ok := s.Data[k]; ok && v == value {
			out = append(out, k)
		}
	}

	return out
}
