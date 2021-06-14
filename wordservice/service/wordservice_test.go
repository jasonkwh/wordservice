package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	pb "deltatre_grpc/proto/wordservice"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024
var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := Server{}
	s.SetDefaultWords()
	grpcServer := grpc.NewServer()
	pb.RegisterWordServiceServer(
		grpcServer,
		&s,
	)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestAddWords(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewWordServiceClient(conn)
	resp, err := client.AddWords(ctx, &pb.AddWordsRequest{Values:[]string{"test1","test2"}})
	if err != nil {
		t.Fatal(err)
	}
	tempwords, contains1, contains2 := resp.GetWords(), -1, -1
	for index, word := range tempwords {
		if word.Value == "test1" {
			contains1 = index
		}
		if word.Value == "test2" {
			contains2 = index
		}
	}

	//verify results
	if contains1 == -1 || contains2 == -1 {
		t.Errorf("AddWords(): add words testing failed, 1st added record index: %v, 2nd added record index: %v, response: %v", contains1, contains2, resp)
		return
	}
}

func TestSearchWord(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewWordServiceClient(conn)
	resp, err := client.SearchWord(ctx, &pb.SearchWordRequest{Value:"filter"})
	if err != nil {
		t.Fatal(err)
	}
	tempwords, contains, searchcount := resp.GetWords(), -1, -1
	for index, word := range tempwords {
		if word.Value == "filter" {
			contains = index
			searchcount = int(word.SearchCount)
		}
	}

	//verify results
	if contains == -1 || len(tempwords) != 1 || searchcount != 1 {
		t.Errorf("SearchWords(): search words testing failed, record index: %v, returned words length: %v, searchcount: %v, response: %v", contains, len(tempwords), searchcount, resp)
		return
	}
}

func TestUpdateWord(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewWordServiceClient(conn)
	resp, err := client.UpdateWord(ctx, &pb.UpdateWordRequest{OrigValue:"filter",NewValue:"newvalue"})
	if err != nil {
		t.Fatal(err)
	}
	tempwords, contains := resp.GetWords(), -1
	for index, word := range tempwords {
		if word.Value == "newvalue" {
			contains = index
		}
	}

	//verify results
	if contains == -1 {
		t.Errorf("UpdateWord(): update word testing failed, record index: %v, response: %v", contains, resp)
		return
	}
}

func TestTopSearches(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewWordServiceClient(conn)
	testwords := []string{"yes","yes","hello"} //search these words
	for i := 0; i < len(testwords); i++ {
		_, err = client.SearchWord(ctx, &pb.SearchWordRequest{Value:testwords[i]})
		if err != nil {
			t.Fatal(err)
		}
	}
	resp, err := client.TopSearches(ctx, &pb.TopSearchesRequest{})
	if err != nil {
		t.Fatal(err)
	}
	tempwords := resp.GetWords()

	//verify results
	if len(tempwords) != 5 || int(tempwords[0].SearchCount) != 2 || int(tempwords[1].SearchCount) != 1 || tempwords[0].Value != "yes" || tempwords[1].Value != "hello" {
		t.Errorf("TopSearches(): top five searches testing failed, response length: %v, response: %v", len(tempwords), resp)
		return
	}
}