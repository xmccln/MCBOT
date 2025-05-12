package Message

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type OnebotGroupMessage struct {
	RawMessage      string      `json:"raw_message"`
	MessageID       int         `json:"message_id"`
	GroupID         int64       `json:"group_id"` // Can be either string or int depending on p.Settings.CompleteFields
	MessageType     string      `json:"message_type"`
	PostType        string      `json:"post_type"`
	SelfID          int64       `json:"self_id"` // Can be either string or int
	Sender          Sender      `json:"sender"`
	SubType         string      `json:"sub_type"`
	Time            int64       `json:"time"`
	Avatar          string      `json:"avatar,omitempty"`
	Echo            string      `json:"echo,omitempty"`
	Message         interface{} `json:"message"` // For array format
	MessageSeq      int         `json:"message_seq"`
	Font            int         `json:"font"`
	UserID          int64       `json:"user_id"`
	RealMessageType string      `json:"real_message_type,omitempty"`  //当前信息的真实类型 group group_private guild guild_private
	RealUserID      string      `json:"real_user_id,omitempty"`       //当前真实uid
	RealGroupID     string      `json:"real_group_id,omitempty"`      //当前真实gid
	IsBindedGroupId bool        `json:"is_binded_group_id,omitempty"` //当前群号是否是binded后的
	IsBindedUserId  bool        `json:"is_binded_user_id,omitempty"`  //当前用户号号是否是binded后的
}

type Sender struct {
	Nickname string `json:"nickname"`
	TinyID   string `json:"tiny_id"`
	UserID   int64  `json:"user_id"`
	Role     string `json:"role,omitempty"`
	Card     string `json:"card,omitempty"`
	Sex      string `json:"sex,omitempty"`
	Age      int32  `json:"age,omitempty"`
	Area     string `json:"area,omitempty"`
	Level    string `json:"level,omitempty"`
	Title    string `json:"title,omitempty"`
}

type OnebotAction struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}

type OneBotGroupMessageParams struct {
	Group_id string `json:"group_id"`
	Message  string `json:"message"`
}

type OneBotPrivateMessageParams struct {
	User_id string `json:"user_id"`
	Message string `json:"message"`
}

func HandleMessage(msg []byte) (*OnebotGroupMessage, error) {
	var message OnebotGroupMessage
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Printf("❌ 解析消息失败: %v", err)
		return nil, err
	}

	// 根据事件类型处理不同消息
	switch message.PostType {
	case "message":
		handleMessageEvent(&message)
	case "notice":
		handleNoticeEvent(&message)
	case "request":
		handleRequestEvent(&message)
	case "meta_event":
		log.Printf("🔄 [元事件] 心跳或生命周期事件")
	default:
		log.Printf("⚠ [未知事件类型] %s", message.PostType)
	}

	return &message, nil
}

// 处理消息事件
func handleMessageEvent(message *OnebotGroupMessage) {
	switch message.MessageType {
	case "private":
		log.Printf("💬 [私聊消息] 发送者: %s(%d)", message.Sender.Nickname, message.UserID)
	case "group":
		log.Printf("👥 [群聊消息] 群号: %d, 发送者: %s(%d), 角色: %s",
			message.GroupID, message.Sender.Nickname, message.UserID, message.Sender.Role)
	default:
		log.Printf("⚠ [未知消息类型] %s", message.MessageType)
		return
	}

	// 处理消息内容
	processMessageContent(message)
}

// 处理通知事件
func handleNoticeEvent(message *OnebotGroupMessage) {
	switch message.SubType {
	case "group_upload":
		log.Printf("📤 [群文件上传] 群号: %d, 用户: %d", message.GroupID, message.UserID)
	case "group_admin":
		log.Printf("👑 [群管理员变动] 群号: %d, 用户: %d", message.GroupID, message.UserID)
	case "group_decrease":
		log.Printf("👋 [群成员减少] 群号: %d, 用户: %d", message.GroupID, message.UserID)
	case "group_increase":
		log.Printf("🎉 [群成员增加] 群号: %d, 用户: %d", message.GroupID, message.UserID)
	case "friend_add":
		log.Printf("🤝 [好友添加] 用户: %d", message.UserID)
	default:
		log.Printf("📢 [通知事件] 类型: %s", message.SubType)
	}
}

// 处理请求事件
func handleRequestEvent(message *OnebotGroupMessage) {
	switch message.SubType {
	case "friend":
		log.Printf("❓ [好友请求] 用户: %d", message.UserID)
	case "group":
		log.Printf("❓ [群请求] 群号: %d, 用户: %d", message.GroupID, message.UserID)
	default:
		log.Printf("❓ [请求事件] 类型: %s", message.SubType)
	}
}

// 处理消息内容
func processMessageContent(message *OnebotGroupMessage) {
	log.Printf("📝 原始消息: %s", message.RawMessage)

	switch content := message.Message.(type) {
	case string:
		// 字符串形式的消息
		log.Printf("📄 字符串消息: %s", content)

	case []interface{}:
		// 数组形式的消息段
		log.Printf("📋 消息段列表:")
		for i, segment := range content {
			segMap, ok := segment.(map[string]interface{})
			if !ok {
				continue
			}

			segType, ok := segMap["type"].(string)
			if !ok {
				continue
			}

			data, ok := segMap["data"].(map[string]interface{})
			if !ok {
				continue
			}

			log.Printf("  📌 消息段[%d]: 类型=%s", i, segType)

			// 根据消息段类型处理
			switch segType {
			case "text":
				if text, ok := data["text"].(string); ok {
					log.Printf("    ✏️ 文本: %s", text)
				}
			case "image":
				log.Printf("    🖼️ 图片:")
				if url, ok := data["url"].(string); ok {
					log.Printf("      🔗 URL: %s", url)
				}
				if file, ok := data["file"].(string); ok {
					log.Printf("      📂 文件: %s", file)
				}
			case "at":
				if qq, ok := data["qq"].(string); ok {
					log.Printf("    👉 @用户: %s", qq)
				}
			case "reply":
				if id, ok := data["id"].(string); ok {
					log.Printf("    💬 回复消息ID: %s", id)
				}
			case "face":
				if id, ok := data["id"].(string); ok {
					log.Printf("    😀 表情ID: %s", id)
				}
			default:
				log.Printf("    ❓ 未处理的类型: %s, 数据: %v", segType, data)
			}
		}

	default:
		log.Printf("⚠️ 未知消息格式: %T", content)
	}

	// 如果有特殊字段，也输出它们
	if message.RealMessageType != "" {
		log.Printf("🏷️ 真实消息类型: %s", message.RealMessageType)
	}
	if message.RealUserID != "" {
		log.Printf("👤 真实用户ID: %s", message.RealUserID)
	}
	if message.RealGroupID != "" {
		log.Printf("👥 真实群组ID: %s", message.RealGroupID)
	}
}

func GetmsgData(msg *OnebotGroupMessage) ([]interface{}, bool) {
	msgData, ok := msg.Message.([]interface{})
	if !ok {
		return nil, false
	}
	if len(msgData) == 0 {
		log.Printf("消息数据为空")
		return nil, false
	}
	return msgData, true
}

func BuildSendPayload(action string, message string, ID int64) (interface{}, error) {
	if action == "send_private_msg" {
		privateParams := OneBotPrivateMessageParams{
			User_id: fmt.Sprint(ID),
			Message: message,
		}
		payload := OnebotAction{
			Action: action,
			Params: privateParams,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			log.Printf("❌ 构建发送数据失败: %v", err)
			return nil, err
		}

		// 反序列化为对象
		var payloadObj map[string]interface{}
		if err := json.Unmarshal(data, &payloadObj); err != nil {
			log.Printf("❌ 反序列化发送数据失败: %v", err)
			return nil, err
		}

		return payloadObj, nil
	} else if action == "send_group_msg" {
		groupParams := OneBotGroupMessageParams{
			Group_id: fmt.Sprint(ID),
			Message:  message,
		}
		payload := OnebotAction{
			Action: action,
			Params: groupParams,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			log.Printf("❌ 构建发送数据失败: %v", err)
			return nil, err
		}

		// 反序列化为对象
		var payloadObj map[string]interface{}
		if err := json.Unmarshal(data, &payloadObj); err != nil {
			log.Printf("❌ 反序列化发送数据失败: %v", err)
			return nil, err
		}

		return payloadObj, nil
	}

	return nil, errors.New("未知的action类型")
}
