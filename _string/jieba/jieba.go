package _string

import (
	"github.com/wangbin/jiebago"
	"runtime"
	"strings"
)

type WordServer struct {
	Seg *jiebago.Segmenter
}

func NewWordServer() (*WordServer, error) {
	server := &WordServer{}
	server.Seg = &jiebago.Segmenter{}
	_, file, _, _ := runtime.Caller(0)
	err := server.Seg.LoadDictionary(strings.Replace(file, "jieba.go", "", -1) + "dict.txt")
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (s *WordServer) Search(word string, hmm bool) []string {
	var values []string
	search := s.Seg.CutForSearch(word, hmm)
	for val := range search {
		values = append(values, val)
	}
	return values
}

func (s *WordServer) Cut(word string, hmm bool) []string {
	var values []string
	search := s.Seg.Cut(word, hmm)
	for val := range search {
		values = append(values, val)
	}
	return values
}

func (s *WordServer) CutAll(word string) []string {
	var values []string
	search := s.Seg.CutAll(word)
	for val := range search {
		values = append(values, val)
	}
	return values
}

func (s *WordServer) RemoveWord(word string) {
	s.Seg.DeleteWord(word)
}

func (s *WordServer) AddWord(word string, frequency float64) {
	s.Seg.AddWord(word, frequency)
}
