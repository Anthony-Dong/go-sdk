package collections

import "sync"

type Set interface {
	ToSlice() []string
	Put(key ...string)
	Delete(key ...string)
	Contains(key string) bool
	Size() int
}

type set struct {
	sync.RWMutex
	body map[string]struct{}
}

func NewSet(slice []string) Set {
	body := SliceToMap(slice)
	return &set{
		body: body,
	}
}

func NewSetInit() Set {
	return &set{
		body: make(map[string]struct{}, 8),
	}
}

func NewSetInitSize(size int) Set {
	return &set{
		body: make(map[string]struct{}, size),
	}
}

func NewSetV2(slice ...string) Set {
	body := SliceToMap(slice)
	return &set{
		body: body,
	}
}

func (s *set) Contains(key string) bool {
	s.RLock()
	defer s.RUnlock()
	_, isExist := s.body[key]
	return isExist
}

func (s *set) ToSlice() []string {
	result := make([]string, 0, len(s.body))
	s.RLock()
	defer s.RUnlock()
	for k := range s.body {
		result = append(result, k)
	}
	return result
}

func (s *set) Put(key ...string) {
	s.Lock()
	defer s.Unlock()
	for _, elem := range key {
		s.body[elem] = struct{}{}
	}
}

func (s *set) Delete(key ...string) {
	s.Lock()
	defer s.Unlock()
	for _, elem := range key {
		delete(s.body, elem)
	}
}

func (s *set) Size() int {
	return len(s.body)
}

func SliceToMap(slice []string) map[string]struct{} {
	if len(slice) == 0 {
		return map[string]struct{}{}
	}
	result := make(map[string]struct{}, len(slice))
	for _, elem := range slice {
		result[elem] = struct{}{}
	}
	return result
}
