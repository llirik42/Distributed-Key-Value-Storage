package main

import (
	"distributed-algorithms/dto"
	"distributed-algorithms/transport/grpc"
	"log"
)

//func electionTimer() {
//	minElectionTimeoutMs := 5000
//	maxElectionTimeoutMs := 7000
//
//	electionTimeoutMs := rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs+1) + minElectionTimeoutMs
//
//	duration := time.Duration(electionTimeoutMs) * time.Millisecond
//
//	log.Println("Starting election timer", duration)
//	timer := time.NewTimer(duration)
//	log.Println("1", duration)
//	<-timer.C
//	log.Println("2", duration)
//	timer.Stop()
//	timer.Reset(duration)
//	log.Println("Timer elapsed")
//}
//
//func leaderLoop() {
//	for {
//		minElectionTimeoutMs := 500
//		maxElectionTimeoutMs := 700
//
//		electionTimeoutMs := rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs+1) + minElectionTimeoutMs
//		duration := time.Duration(electionTimeoutMs) * time.Millisecond
//
//		time.Sleep(duration)
//
//		println("Leader loop elapsed with ", duration)
//	}
//}
//
//func followerLoop() {
//	for {
//		minElectionTimeoutMs := 5000
//		maxElectionTimeoutMs := 7000
//
//		electionTimeoutMs := rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs+1) + minElectionTimeoutMs
//		duration := time.Duration(electionTimeoutMs) * time.Millisecond
//
//		time.Sleep(duration)
//		println("Follower loop elapsed with ", duration)
//	}
//}

func handler1(request dto.RequestVoteRequest) (*dto.RequestVoteResponse, error) {
	log.Printf("Received voting: %v", request)
	return &dto.RequestVoteResponse{}, nil
}

func handler2(request dto.AppendEntriesRequest) (*dto.AppendEntriesResponse, error) {
	log.Printf("Received append: %v", request)
	return &dto.AppendEntriesResponse{}, nil
}

func main() {
	serverFactory := &grpc.ServerFactory{}

	server, err := serverFactory.NewServer("0.0.0.0:8080", handler1, handler2)

	if err != nil {
		panic(err)
	}

	err1 := server.Listen()

	if err1 != nil {
		panic(err1)
	}

	//go followerLoop()
	//go leaderLoop()
	//
	//input, _ := fmt.Scanln()
	//
	//fmt.Println(input)
	//
	//listener, err := net.Listen("tcp", "127.0.0.1:9111")
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//} else {
	//	log.Println("listening on :9111")
	//}
	//
	//go electionTimer()
	//
	//grpcServer := grpc.NewServer()
	//service := &grpc2.RaftServiceServer{
	//	On: nil,
	//}
	//
	//pb.RegisterRaftServiceServer(grpcServer, service)
	//reflection.Register(grpcServer)
	//err = grpcServer.Serve(listener)
	//
	//if err != nil {
	//	log.Fatalf("Error starting server: %v", err)
	//}
}
