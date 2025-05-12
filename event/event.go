package event

func Isat(msgData []interface{}, SelfID string) bool {
	if segment, ok := msgData[0].(map[string]interface{}); ok {
		if segment["type"] == "at" {
			if target, ok := segment["data"].(map[string]interface{}); ok &&
				target["qq"] == SelfID {
				return true
			}
		}
	}
	return false
}
