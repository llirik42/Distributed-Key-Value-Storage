package executor

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/log"
	"fmt"
)

// Result of executing command
const (
	Success = iota
	Failure
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
		// TODO: handle unexpected error
		return
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

	var executedCommandCode int
	var executedCommandError string

	if cmd.Key != e.executedCommandsKey {
		storage.Set(cmd.Key, cmd.NewValue)
		executedCommandCode = Success
		executedCommandError = ""
	} else {
		executedCommandCode = Failure
		executedCommandError = "cannot set value of internal key"
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, executedCommandCode, executedCommandError)
	}
}

func (e *CommandExecutor) executeDelete(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var executedCommandCode int
	var executedCommandError string

	if cmd.Key != e.executedCommandsKey {
		storage.Delete(cmd.Key)
		executedCommandCode = Success
		executedCommandError = ""
	} else {
		executedCommandCode = Failure
		executedCommandError = "cannot delete internal key"
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, executedCommandCode, executedCommandError)
	}
}

func (e *CommandExecutor) executeCompareAndSet(cmd log.Command, isLeader bool) {
	storage := e.ctx.GetKeyValueStorage()

	var executedCommandCode = Success
	var executedCommandError = ""

	if cmd.Key != e.executedCommandsKey {
		success, err := storage.CompareAndSet(cmd.Key, cmd.OldValue, cmd.NewValue)

		if !success || err != nil {
			executedCommandCode = Failure
		}

		if err != nil {
			executedCommandError = err.Error()
		}
	} else {
		executedCommandCode = Failure
		executedCommandError = "cannot executed compare-and-set with internal key"
	}

	if isLeader {
		e.pushCommandForExecutionInfo(cmd.Id, executedCommandCode, executedCommandError)
	}
}

func (e *CommandExecutor) executeAddElement(cmd log.Command) {
	storage := e.ctx.GetKeyValueStorage()
	storage.AddElement(cmd.Key, cmd.SubKey, cmd.NewValue)
	// Don't push anything to log because AddElement is an internal command (not caused by client)
}

// Pushes to log command that will add info about execution of other command
func (e *CommandExecutor) pushCommandForExecutionInfo(commandId string, code int, err string) {
	e.ctx.PushCommand(log.Command{
		Key:    e.executedCommandsKey,
		SubKey: commandId,
		NewValue: map[string]any{
			"code":  code,
			"error": err,
		},
		Type: log.AddElement,
	})
}
