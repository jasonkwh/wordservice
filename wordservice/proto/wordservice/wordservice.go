package wordservice

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

type Server struct {
	UnimplementedWordServiceServer
	words WordStorage
}

func (s *Server) SetDefaultWords() {
	value := []string{"hello","goodbye","simple","list","search","filter","yes","no"}
	s.words.AddWords(&value)
}

func (s *Server) AddWords(ctx context.Context, in *AddWordsRequest) (*WordsResponse, error) {
	s.words.ClearExistInputWords(&in.Values)
	s.words.AddWords(&in.Values)
	return &WordsResponse{Words:s.words.items}, nil
}

func (s *Server) SearchWord(ctx context.Context, in *SearchWordRequest) (*WordsResponse, error) {
	fmt.Println(in.Value)
	return &WordsResponse{Words:s.words.items}, nil
}

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
		if origValueIndex != 0 && in.OrigValue != "" && newValueIndex == 0 && in.NewValue != "" {
			s.words.items[origValueIndex].Value = in.NewValue
			s.words.items[origValueIndex].ModifiedTime = timestamppb.Now()
		}
	}
	return &WordsResponse{Words:s.words.items}, nil
}
