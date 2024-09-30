package _string

import "github.com/wangbin/jiebago"

type WordServer struct {
	Seg *jiebago.Segmenter
}

func (s *WordServer) NewWordServer() (*WordServer, error) {
	s.Seg = &jiebago.Segmenter{}
	err := s.Seg.LoadDictionary("_string/jieba/dict.txt") //初始化分词的字典，可手动修改该文件
	if err != nil {
		return nil, err
	}
	return s, nil
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
