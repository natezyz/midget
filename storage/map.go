package storage

import (
	"errors"
	"hash/fnv"
	"strconv"
	"sync"
)

type Map struct {
	Data map[uint64]string
	sync.RWMutex
}

func (s *Map) Init() {
	s.Data = make(map[uint64]string)
}

func (s *Map) Store(url string) string {
	encoded := fnv64(url)
	s.Lock()
	s.Data[encoded] = url
	s.Unlock()
	return strconv.FormatUint(encoded, 10)
}

func (s *Map) Retrieve(encoded string) (string, error) {
	if encoded == "" {
		return "", errors.New("Code was empty")
	}
	key, err := strconv.ParseUint(encoded, 10, 64)
	if err != nil {
		return "", err
	}

	if val, ok := s.Data[key]; ok {
		return val, nil
	}
	return "", errors.New(encoded + " not found")
}

func fnv64(text string) uint64 {
	algo := fnv.New64a()
	algo.Write([]byte(text))
	return algo.Sum64()
}
