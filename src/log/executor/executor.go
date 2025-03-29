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

func NewCommandExecutor(ctx *context.Context, executedCommandsKey string) *CommandExecutor {
	return &CommandExecutor{
		ctx:                 ctx,
		executedCommandsKey: executedCommandsKey,
	}
}

func (e *CommandExecutor) Execute(cmd log.Command) {
	if cmd.Id == "" {
		panic(fmt.Errorf("command id should not be empty"))
	}

	isLeader := e.ctx.IsLeader()

	switch cmd.Type {
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

func (e *CommandExecutor) executeSet(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var commandExecutionResult string

	if cmd.Key != e.executedCommandsKey {
		storage.Set(cmd.Key, cmd.NewValue)
		commandExecutionResult = Success
	} else {
		commandExecutionResult = "cannot set value of internal key"
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, commandExecutionResult)
	}
}

func (e *CommandExecutor) executeDelete(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var commandExecutionResult string

	if cmd.Key != e.executedCommandsKey {
		storage.Delete(cmd.Key)
		commandExecutionResult = Success
	} else {
		commandExecutionResult = "cannot delete internal key"
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, commandExecutionResult)
	}
}

func (e *CommandExecutor) executeCompareAndSet(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var commandExecutionResult = Success

	if cmd.Key != e.executedCommandsKey {
		success, err := storage.CompareAndSet(cmd.Key, cmd.OldValue, cmd.NewValue)

		if err != nil {
			commandExecutionResult = err.Error()
		} else if !success {
			commandExecutionResult = "old value does not match"
		}
	} else {
		commandExecutionResult = "cannot executed compare-and-set with internal key"
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, commandExecutionResult)
	}
}

func (e *CommandExecutor) executeAddElement(cmd log.Command) {
	storage := e.ctx.GetKeyValueStorage()
	storage.AddElement(cmd.Key, cmd.SubKey, cmd.NewValue)
	// Don't push anything to log because AddElement is an internal command (not caused by client)
}

// Pushes to log command that will add info about execution of other command
func (e *CommandExecutor) pushCommandForExecutionInfo(commandId string, result string) {
	e.ctx.PushCommand(log.Command{
		Key:      e.executedCommandsKey,
		SubKey:   commandId,
		NewValue: result,
		Type:     log.AddElement,
	})
}
