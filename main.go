package main

import (
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/transport/grpc"
	"fmt"
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

func main() {
	_, err := raft.NewNode(grpc.ServerFactory{})

	if err != nil {
		fmt.Println(err)
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
