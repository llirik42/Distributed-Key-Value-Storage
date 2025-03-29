package executor

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/log"
	"fmt"
)

const (
	Success = "success"
)

type CommandExecutor struct {
	ctx                 *context.Context
	executedCommandsKey string
}

func NewCommandExecutor(ctx *context.Context) *CommandExecutor {
	return &CommandExecutor{
		ctx:                 ctx,
		executedCommandsKey: ctx.GetExecutedCommandsKey(),
	}
}

func (e *CommandExecutor) Execute(cmd log.Command) {
	if cmd.Id == "" {
		panic(fmt.Errorf("command id should not be empty"))
	}

	isLeader := e.ctx.IsLeader()

	switch cmd.Type {
	case log.Get:
		e.executeGet(cmd, isLeader)
	case log.Set:
		e.executeSet(cmd, isLeader)
	case log.CompareAndSet:
		e.executeCompareAndSet(cmd, isLeader)
	case log.Delete:
		e.executeDelete(cmd, isLeader)
	case log.AddElement:
		e.executeAddElement(cmd)
	default:
		panic(fmt.Errorf("unknown type of command %+v", cmd))
	}
}

func (e *CommandExecutor) executeGet(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var info CommandExecutionInfo

	if cmd.Key != e.executedCommandsKey {
		value := storage.Get(cmd.Key)

		if !value.Exists {
			info.Message = "key not found"
			info.Success = false
		} else {
			info.Message = Success
			info.Success = true
			info.Value = value.Value
		}
	} else {
		info.Message = "cannot get value of internal key"
		info.Success = false
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, info)
	}
}

func (e *CommandExecutor) executeSet(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var info CommandExecutionInfo

	if cmd.Key != e.executedCommandsKey {
		storage.Set(cmd.Key, cmd.NewValue)
		info.Message = Success
		info.Success = true
	} else {
		info.Message = "cannot set value of internal key"
		info.Success = false
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, info)
	}
}

func (e *CommandExecutor) executeDelete(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var info CommandExecutionInfo

	if cmd.Key != e.executedCommandsKey {
		storage.Delete(cmd.Key)
		info.Message = Success
		info.Success = true
	} else {
		info.Message = "cannot delete internal key"
		info.Success = false
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, info)
	}
}

func (e *CommandExecutor) executeCompareAndSet(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var info = CommandExecutionInfo{Success: false}

	if cmd.Key != e.executedCommandsKey {
		success, err := storage.CompareAndSet(cmd.Key, cmd.OldValue, cmd.NewValue)

		if err != nil {
			info.Message = err.Error()
		} else if !success {
			info.Message = "old value does not match"
		} else {
			info.Message = Success
			info.Success = true
		}
	} else {
		info.Message = "cannot executed compare-and-set with internal key"
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, info)
	}
}

func (e *CommandExecutor) executeAddElement(cmd log.Command) {
	storage := e.ctx.GetKeyValueStorage()
	storage.AddElement(cmd.Key, cmd.SubKey, cmd.NewValue)
	// Don't push anything to log because AddElement is an internal command (not caused by client)
}

// Pushes to log command that will add info about execution of other command
func (e *CommandExecutor) pushCommandForExecutionInfo(commandId string, info CommandExecutionInfo) {
	e.ctx.PushCommand(log.Command{
		Key:      e.executedCommandsKey,
		SubKey:   commandId,
		NewValue: commandExecutionInfoToMap(info),
		Type:     log.AddElement,
	})
}
