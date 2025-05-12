package main

import (
	"MCWhitelist/Message"
	"MCWhitelist/event"
	"MCWhitelist/whitelist"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var db *sql.DB

func initdb(conn *sql.DB) {
	// 创建表格 whitelist 其中包含 stings类型的 qq 字段 还有gameid字段
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS whitelist (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		qq TEXT NOT NULL,
		gameid TEXT NOT NULL,
		status INTEGER DEFAULT 0
	);
	`
	_, err := conn.Exec(sqlStmt)
	if err != nil {
		log.Fatal("创建表格失败:", err)
	}
	log.Println("数据库初始化成功")

}

func handleConnection(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级连接失败:", err)
		return
	}
	defer conn.Close()

	log.Println("客户端已连接")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息错误:", err)
			break
		}
		msg, err := Message.HandleMessage(message)
		if err != nil {
			log.Println("处理消息错误:", err)
		}

		switch msg.PostType {
		case "message":
			msgData, _ := Message.GetmsgData(msg)

			if event.Isat(msgData, fmt.Sprint(msg.SelfID)) {
				log.Println("被at了")
				if len(msgData) <= 1 {
					log.Println("msgData 长度不足")
					return
				}
				segMap, _ := msgData[1].(map[string]interface{})
				data, _ := segMap["data"].(map[string]interface{})
				part := strings.Fields(data["text"].(string))
				switch part[0] {
				case "help":
					log.Println("help")
				case "申请白名单":
					if len(part) == 2 {

						playload, err := Message.BuildSendPayload("send_group_msg",
							whitelist.ReqWhitelist(db, strconv.FormatInt(msg.UserID, 10), part[1]), msg.GroupID)
						if err != nil {
							log.Println("构建发送数据失败:", err)
							continue
						}
						conn.WriteJSON(playload)
						continue
					} else {
						playload, err := Message.BuildSendPayload("send_group_msg",
							"请提供游戏id", msg.GroupID)
						if err != nil {
							log.Println("构建发送数据失败:", err)
							continue
						}
						conn.WriteJSON(playload)
						continue
					}

				default:
					log.Println("未知指令")
				}
			}
		default:
			log.Printf("未知消息类型: %s", msg.PostType)

		}
	}

	log.Println("客户端已断开连接")
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "D:/poject/Go/t2/MCWhitelist.db")
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()
	initdb(db)

	http.HandleFunc("/ws", handleConnection)

	server := &http.Server{
		Addr: ":8080",
	}

	log.Println("WebSocket服务器启动在 :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
