package wordservice

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

type WordStorage struct {
	items []*Word
}

func (w *WordStorage) AddWords(values *[]string) {
	for _, v := range *values {
		var word Word
		word.Id = int64(len(w.items) + 1)
		word.Value = v
		word.SearchCount = 0
		word.AddedTime = timestamppb.Now()
		w.items = append(w.items, &word)
	}
}

func (w *WordStorage) SearchWords(value *string) int {
	for index, word := range w.items {
		if strings.Contains(word.Value, *value) {
			return index
		}
	}
	return 0
}

func (w *WordStorage) GetIndex(value *string) int {
	for index, word := range w.items {
		if word.Value == (*value) {
			return index
		}
	}
	return 0
}

func (w *WordStorage) IsContain(value *string) bool {
	if w.GetIndex(value) == 0 {
		return false
	}
	return true
}

func (w *WordStorage) ClearExistInputWords(values *[]string) {
	if len(w.items) > 0 && len(*values) > 0 {
		var result []string
		for _, v := range *values {
			if v != "" && w.IsContain(&v) == false {
				result = append(result, v)
			}
		}
		*values = result
	}
}