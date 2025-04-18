package grpc

import (
	pb "distributed-key-value-storage/generated/proto"
	"distributed-key-value-storage/src/log"
	"distributed-key-value-storage/src/raft/dto"
	"fmt"
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
	mappedEntries := make([]log.Entry, len(msg.Entries))

	for i, el := range msg.Entries {
		mappedEntries[i] = log.Entry{
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
		ConflictTerm:  msg.ConflictTerm,
		ConflictIndex: msg.ConflictIndex,
	}
}

func MapAppendEntriesResponseFromGRPC(msg *pb.AppendEntriesResponse) *dto.AppendEntriesResponse {
	return &dto.AppendEntriesResponse{
		Term:          msg.Term,
		Success:       msg.Success,
		ConflictTerm:  msg.ConflictTerm,
		ConflictIndex: msg.ConflictIndex,
	}
}

func mapLogEntriesToGRPC(entries *[]log.Entry) []*pb.LogEntry {
	result := make([]*pb.LogEntry, len(*entries))

	for i, el := range *entries {
		result[i] = &pb.LogEntry{
			Term:    el.Term,
			Command: mapCommandToGRPC(&el.Command),
		}
	}

	return result
}

func mapCommandToGRPC(command *log.Command) *pb.Command {
	switch command.Type {
	case log.Get:
		return &pb.Command{
			Id: command.Id,
			Type: &pb.Command_Get_{
				Get: &pb.Command_Get{
					Key: command.Key,
				},
			},
		}
	case log.Set:
		return &pb.Command{
			Id: command.Id,
			Type: &pb.Command_Set_{
				Set: &pb.Command_Set{
					Key:   command.Key,
					Value: mapValueToGRPC(&command.NewValue),
				},
			},
		}
	case log.CompareAndSet:
		return &pb.Command{
			Id: command.Id,
			Type: &pb.Command_CompareAndSet_{
				CompareAndSet: &pb.Command_CompareAndSet{
					Key:      command.Key,
					NewValue: mapValueToGRPC(&command.NewValue),
					OldValue: mapValueToGRPC(&command.OldValue),
				},
			},
		}
	case log.Delete:
		return &pb.Command{
			Id: command.Id,
			Type: &pb.Command_Delete_{
				Delete: &pb.Command_Delete{
					Key: command.Key,
				},
			},
		}
	case log.AddElement:
		return &pb.Command{
			Id: command.Id,
			Type: &pb.Command_AddElement_{
				AddElement: &pb.Command_AddElement{
					Key:    command.Key,
					SubKey: command.SubKey,
					Value:  mapValueToGRPC(&command.NewValue),
				},
			},
		}
	default:
		panic(fmt.Errorf("unknown command type: %d\n", command.Type))
	}
}

func mapCommandFromGRPC(command *pb.Command) *log.Command {
	switch command.Type.(type) {
	case *pb.Command_Get_:
		cmd := command.GetGet()

		return &log.Command{
			Id:   command.Id,
			Key:  cmd.Key,
			Type: log.Get,
		}
	case *pb.Command_Set_:
		cmd := command.GetSet()

		return &log.Command{
			Id:       command.Id,
			Key:      cmd.Key,
			NewValue: mapValueFromGRPC(cmd.Value),
			Type:     log.Set,
		}
	case *pb.Command_CompareAndSet_:
		cmd := command.GetCompareAndSet()

		return &log.Command{
			Id:       command.Id,
			Key:      cmd.Key,
			OldValue: mapValueFromGRPC(cmd.OldValue),
			NewValue: mapValueFromGRPC(cmd.NewValue),
			Type:     log.CompareAndSet,
		}
	case *pb.Command_Delete_:
		cmd := command.GetDelete()

		return &log.Command{
			Id:   command.Id,
			Key:  cmd.Key,
			Type: log.Delete,
		}
	case *pb.Command_AddElement_:
		cmd := command.GetAddElement()

		return &log.Command{
			Id:       command.Id,
			Key:      cmd.Key,
			SubKey:   cmd.SubKey,
			NewValue: mapValueFromGRPC(cmd.Value),
			Type:     log.AddElement,
		}
	default: // Unknown type
		panic(fmt.Errorf("unknown of command: %T\n", command.Type))
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
	default:
		panic(fmt.Errorf("unknown type of value: %T\n", v))
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
	default:
		panic(fmt.Errorf("unknown type of value: %T\n", value.Type))
	}
}
