package wordservice

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

type Server struct {
	UnimplementedWordServiceServer
	words WordStorage
}

/*
Function sets default words in storage
 */
func (s *Server) SetDefaultWords() {
	value := []string{"hello","goodbye","simple","list","search","filter","yes","no"}
	_ = s.words.AddWords(&value)
}

/*
Function which adds new words into storage
 */
func (s *Server) AddWords(ctx context.Context, in *AddWordsRequest) (*WordsResponse, error) {
	err := s.words.ClearExistInputWords(&in.Values)
	if err != nil {
		return &WordsResponse{}, err
	}
	err = s.words.AddWords(&in.Values)
	if err != nil {
		return &WordsResponse{}, err
	}
	return &WordsResponse{Words:s.words.items}, nil
}

/*
Function search occurrence of keyword in word storage
 */
func (s *Server) SearchWord(ctx context.Context, in *SearchWordRequest) (*WordsResponse, error) {
	w, err := s.words.GetSearchWords(&in.Value)
	return &WordsResponse{Words:w}, err
}

/*
Function returns top 5 searches
Order by SearchCount descending, Value ascending
 */
func (s *Server) TopSearches(ctx context.Context, in *TopSearchesRequest) (*WordsResponse, error) {
	w, err := s.words.TopSearchWords()
	return &WordsResponse{Words:w}, err
}

/*
Function to update word
Use sync.WaitGroup to get index of OrigValue & NewValue at the same time
 */
func (s *Server) UpdateWord(ctx context.Context, in *UpdateWordRequest) (*WordsResponse, error) {
	var err error
	if in.NewValue == in.OrigValue {
		err = errors.New("orig_value and new_value cannot be the same")
	}
	if in.NewValue == "" || in.OrigValue == "" {
		err = errors.New("orig_value or new_value cannot be empty")
	}
	if err == nil {
		origValueIndex, newValueIndex := 0, 0
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			origValueIndex = s.words.GetIndex(&in.OrigValue)
		}()
		go func() {
			defer wg.Done()
			newValueIndex = s.words.GetIndex(&in.NewValue)
		}()
		wg.Wait()
		if origValueIndex == -1 {
			err = errors.New("cannot find orig_value in the storage")
		}
		if newValueIndex != -1 {
			err = errors.New("new_value you are trying to set is already exist")
		}
		if err == nil {
			s.words.items[origValueIndex].Value = in.NewValue
			s.words.items[origValueIndex].ModifiedTime = timestamppb.Now()
		}
	}
	return &WordsResponse{Words:s.words.items}, err
}
