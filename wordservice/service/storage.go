package service

import (
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "deltatre_grpc/proto/wordservice"
	"sort"
	"strings"
)

type WordStorage struct {
	items []*pb.Word
}

/*
Function inserts words into storage
 */
func (w *WordStorage) AddWords(values *[]string) error {
	if len(*values) > 0 {
		for _, v := range *values {
			var word pb.Word
			word.Id = int64(len(w.items) + 1)
			word.Value = v
			word.SearchCount = 0
			word.AddedTime = timestamppb.Now()
			w.items = append(w.items, &word)
		}
		return nil
	}
	return errors.New("no new words")
}

/*
Return top five searches
Order by SearchCount descending, Value ascending
 */
func (w *WordStorage) TopSearchWords() ([]*pb.Word, error) {
	if len(w.items) >= 5 {
		sort.SliceStable(w.items, func(i, j int) bool {
			if w.items[i].SearchCount != w.items[j].SearchCount {
				return w.items[i].SearchCount > w.items[j].SearchCount
			}
			return w.items[i].Value < w.items[j].Value
		})
		return w.items[:5], nil
	}
	return nil, errors.New("not enough words in storage")
}

/*
Search for a specified pattern in word storage
 */
func (w *WordStorage) GetSearchWords(value *string) ([]*pb.Word, error) {
	var output []*pb.Word
	if len(w.items) > 0 {
		for _, word := range w.items {
			//Search for a single word
			//if word.Value == *value {
			if strings.Contains(word.Value, *value) {
				word.SearchCount++
				output = append(output, word)
			}
		}
		if len(output) == 0 {
			return output, errors.New("no results")
		}
		return output, nil
	}
	return output, errors.New("storage is empty")
}

/*
Return index/key of a particular word struct in w.items
-1 means word not exist
 */
func (w *WordStorage) GetIndex(value *string) int {
	if len(w.items) > 0 {
		for index, word := range w.items {
			if word.Value == (*value) {
				return index
			}
		}
	}
	return -1
}

/*
Check if w.items has the word or not
 */
func (w *WordStorage) IsContain(value *string) bool {
	if w.GetIndex(value) == -1 {
		return false
	}
	return true
}

/*
Function for AddWords rpc, check existence of input words in storage and clear the identicals
 */
func (w *WordStorage) ClearExistInputWords(values *[]string) error {
	if len(w.items) > 0 {
		if len(*values) > 0 {
			var result []string
			for _, v := range *values {
				if v == "" {
					return errors.New("input word value cannot be empty")
				}
				if w.IsContain(&v) == false {
					result = append(result, v)
				}
			}
			*values = result
			return nil
		}
		return errors.New("no input words")
	}
	return nil
}
