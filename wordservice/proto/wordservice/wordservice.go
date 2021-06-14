package wordservice

import (
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
	s.words.AddWords(&value)
}

/*
Function which adds new words into storage
 */
func (s *Server) AddWords(ctx context.Context, in *AddWordsRequest) (*WordsResponse, error) {
	s.words.ClearExistInputWords(&in.Values)
	s.words.AddWords(&in.Values)
	return &WordsResponse{Words:s.words.items}, nil
}

/*
Function search occurrence of keyword in word storage
 */
func (s *Server) SearchWord(ctx context.Context, in *SearchWordRequest) (*WordsResponse, error) {
	return &WordsResponse{Words:s.words.GetSearchWords(&in.Value)}, nil
}

/*
Function returns top 5 searches
Order by SearchCount descending, Value ascending
 */
func (s *Server) TopSearches(ctx context.Context, in *TopSearchesRequest) (*WordsResponse, error) {
	return &WordsResponse{Words:s.words.TopSearchWords()}, nil
}

/*
Function to update word
Use sync.WaitGroup to get index of OrigValue & NewValue at the same time
 */
func (s *Server) UpdateWord(ctx context.Context, in *UpdateWordRequest) (*WordsResponse, error) {
	if in.NewValue != in.OrigValue {
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
		if origValueIndex != -1 && in.OrigValue != "" && newValueIndex == -1 && in.NewValue != "" {
			s.words.items[origValueIndex].Value = in.NewValue
			s.words.items[origValueIndex].ModifiedTime = timestamppb.Now()
		}
	}
	return &WordsResponse{Words:s.words.items}, nil
}
