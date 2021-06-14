package wordservice

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"sort"
	"strings"
)

type WordStorage struct {
	items []*Word
}

/*
Function inserts words into storage
 */
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

/*
Return top five searches
Order by SearchCount descending, Value ascending
 */
func (w *WordStorage) TopSearchWords() []*Word {
	output := w.items
	sort.SliceStable(output, func(i, j int) bool {
		if output[i].SearchCount != output[j].SearchCount {
			return output[i].SearchCount > output[j].SearchCount
		}
		return output[i].Value < output[j].Value
	})
	return output[:5]
}

/*
Search for a specified pattern in word storage
 */
func (w *WordStorage) GetSearchWords(value *string) []*Word {
	var output []*Word
	for _, word := range w.items {
		if strings.Contains(word.Value, *value) {
			word.SearchCount++
			output = append(output, word)
		}
	}
	return output
}

/*
Return index/key of a particular word struct in w.items
-1 means word not exist
 */
func (w *WordStorage) GetIndex(value *string) int {
	for index, word := range w.items {
		if word.Value == (*value) {
			return index
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
