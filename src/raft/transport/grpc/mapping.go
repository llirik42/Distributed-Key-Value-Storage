package grpc

import (
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/src/domain"
	"distributed-algorithms/src/raft/dto"
)

func MapRequestForVoteRequestToGRPC(msg *dto.RequestVoteRequest) *pb.RequestVoteRequest {
	return &pb.RequestVoteRequest{
		Term:         msg.Term,
		CandidateId:  msg.CandidateId,
		LastLogIndex: msg.LastLogIndex,
		LastLogTerm:  msg.LastLogTerm,
	}
}

func MapRequestForVoteRequestFromGRPC(msg *pb.RequestVoteRequest) *dto.RequestVoteRequest {
	return &dto.RequestVoteRequest{
		Term:         msg.Term,
		CandidateId:  msg.CandidateId,
		LastLogIndex: msg.LastLogIndex,
		LastLogTerm:  msg.LastLogTerm,
	}
}

func MapRequestForVoteResponseToGRPC(msg *dto.RequestVoteResponse) *pb.RequestVoteResponse {
	return &pb.RequestVoteResponse{
		Term:        msg.Term,
		VoteGranted: msg.VoteGranted,
	}
}

func MapRequestForVoteResponseFromGRPC(msg *pb.RequestVoteResponse) *dto.RequestVoteResponse {
	return &dto.RequestVoteResponse{
		Term:        msg.Term,
		VoteGranted: msg.VoteGranted,
	}
}

func MapAppendEntriesRequestToGRPC(msg *dto.AppendEntriesRequest) *pb.AppendEntriesRequest {
	return &pb.AppendEntriesRequest{
		Term:         msg.Term,
		LeaderId:     msg.LeaderId,
		PrevLogIndex: msg.PrevLogIndex,
		PrevLogTerm:  msg.PrevLogTerm,
		Entries:      mapLogEntriesToGRPC(&msg.Entries),
		LeaderCommit: msg.LeaderCommit,
	}
}

func MapAppendEntriesRequestFromGRPC(msg *pb.AppendEntriesRequest) *dto.AppendEntriesRequest {
	mappedEntries := make([]dto.LogEntry, len(msg.Entries))

	for i, el := range msg.Entries {
		mappedEntries[i] = dto.LogEntry{
			Index:   el.Index,
			Term:    el.Term,
			Command: *mapCommandFromGRPC(el.Command),
		}
	}

	return &dto.AppendEntriesRequest{
		Term:         msg.Term,
		LeaderId:     msg.LeaderId,
		PrevLogIndex: msg.PrevLogIndex,
		PrevLogTerm:  msg.PrevLogTerm,
		Entries:      mappedEntries,
		LeaderCommit: msg.LeaderCommit,
	}
}

func MapAppendEntriesResponseToGRPC(msg *dto.AppendEntriesResponse) *pb.AppendEntriesResponse {
	return &pb.AppendEntriesResponse{
		Term:          msg.Term,
		Success:       msg.Success,
		ConflictIndex: msg.ConflictIndex,
		ConflictTerm:  msg.ConflictTerm,
	}
}

func MapAppendEntriesResponseFromGRPC(msg *pb.AppendEntriesResponse) *dto.AppendEntriesResponse {
	return &dto.AppendEntriesResponse{
		Term:          msg.Term,
		Success:       msg.Success,
		ConflictIndex: msg.ConflictIndex,
		ConflictTerm:  msg.ConflictTerm,
	}
}

func mapLogEntriesToGRPC(entries *[]dto.LogEntry) []*pb.LogEntry {
	result := make([]*pb.LogEntry, len(*entries))

	for i, el := range *entries {
		result[i] = &pb.LogEntry{
			Index:   el.Index,
			Term:    el.Term,
			Command: mapCommandToGRPC(&el.Command),
		}
	}

	return result
}

func mapCommandToGRPC(command *domain.Command) *pb.Command {
	switch command.Type {
	case domain.Set:
		return &pb.Command{
			Type: &pb.Command_Set_{
				Set: &pb.Command_Set{
					Key:   command.Key,
					Value: mapValueToGRPC(&command.Value),
				},
			},
		}
	case domain.Delete:
		return &pb.Command{
			Type: &pb.Command_Delete_{
				Delete: &pb.Command_Delete{
					Key: command.Key,
				},
			},
		}
	default: // Unknown type
		return nil
	}
}

func mapCommandFromGRPC(command *pb.Command) *domain.Command {
	switch command.Type.(type) {
	case *pb.Command_Delete_:
		cmd := command.GetDelete()

		return &domain.Command{
			Key:   cmd.Key,
			Value: nil,
			Type:  domain.Delete,
		}
	case *pb.Command_Set_:
		cmd := command.GetSet()

		return &domain.Command{
			Key:   cmd.Key,
			Value: mapValueFromGRPC(cmd.Value),
			Type:  domain.Set,
		}
	default: // Unknown type
		return nil
	}
}

func mapValueToGRPC(value *any) *pb.Value {
	v := *value

	switch v.(type) {
	case nil:
		return &pb.Value{
			Type: &pb.Value_Null_{
				Null: &pb.Value_Null{},
			},
		}
	case bool:
		return &pb.Value{
			Type: &pb.Value_Boolean_{
				Boolean: &pb.Value_Boolean{Value: v.(bool)},
			},
		}
	case float64:
		return &pb.Value{
			Type: &pb.Value_Number_{
				Number: &pb.Value_Number{Value: v.(float64)},
			},
		}
	case string:
		return &pb.Value{
			Type: &pb.Value_String_{
				String_: &pb.Value_String{Value: v.(string)},
			},
		}
	case []any:
		pbValue := make([]*pb.Value, len(v.([]any)))

		for i, el := range v.([]any) {
			pbValue[i] = mapValueToGRPC(&el)
		}

		return &pb.Value{
			Type: &pb.Value_Array_{
				Array: &pb.Value_Array{Value: pbValue},
			},
		}
	case map[string]any:
		pbValue := map[string]*pb.Value{}

		for mapKey, mapValue := range v.(map[string]any) {
			pbValue[mapKey] = mapValueToGRPC(&mapValue)
		}

		return &pb.Value{
			Type: &pb.Value_Object_{
				Object: &pb.Value_Object{Value: pbValue},
			},
		}
	default: // Unknown type
		return nil
	}
}

func mapValueFromGRPC(value *pb.Value) any {
	switch value.Type.(type) {
	case *pb.Value_Null_:
		return nil
	case *pb.Value_Boolean_:
		return value.GetBoolean().Value
	case *pb.Value_Number_:
		return value.GetNumber().Value
	case *pb.Value_String_:
		return value.GetString_().Value
	case *pb.Value_Array_:
		a := value.GetArray().Value

		result := make([]any, len(a))

		for i, el := range a {
			result[i] = mapValueFromGRPC(el)
		}

		return result
	case *pb.Value_Object_:
		result := map[string]any{}

		for objKey, objValue := range value.GetObject().Value {
			result[objKey] = mapValueFromGRPC(objValue)
		}

		return result
	default: // Unknown type
		return nil
	}
}
