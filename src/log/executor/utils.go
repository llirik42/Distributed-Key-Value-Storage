package executor

import "distributed-algorithms/src/key-value"

func GetCommandExecutionInfo(
	storage key_value.Storage,
	executedCommandsKey string,
	commandId string,
) (CommandExecutionInfo, bool) {
	value := storage.GetElement(executedCommandsKey, commandId)

	if !value.Exists {
		return CommandExecutionInfo{}, false
	}

	return commandExecutionInfoFromMap(value.Value.(map[string]any)), true
}

func commandExecutionInfoToMap(info CommandExecutionInfo) map[string]any {
	return map[string]any{
		"value":   info.Value,
		"message": info.Message,
		"success": info.Success,
	}
}

func commandExecutionInfoFromMap(info map[string]any) CommandExecutionInfo {
	return CommandExecutionInfo{
		Value:   info["value"],
		Message: info["message"].(string),
		Success: info["success"].(bool),
	}
}
