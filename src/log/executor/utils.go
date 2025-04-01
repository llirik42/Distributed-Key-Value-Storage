package executor

import (
	"distributed-algorithms/src/log"
)

func commandExecutionInfoToMap(info log.CommandExecutionInfo) map[string]any {
	return map[string]any{
		"value":   info.Value,
		"message": info.Message,
		"success": info.Success,
	}
}

func commandExecutionInfoFromMap(info map[string]any) log.CommandExecutionInfo {
	return log.CommandExecutionInfo{
		Value:   info["value"],
		Message: info["message"].(string),
		Success: info["success"].(bool),
	}
}
