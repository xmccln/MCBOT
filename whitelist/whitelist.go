package whitelist

import (
	"database/sql"
	"fmt"
	"log"
)

func ReqWhitelist(conn *sql.DB, qq string, gameid string) string {
	// 通过qq查询申请数量
	if gameid == "" {
		return "游戏ID不能为空"
	}

	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM whitelist WHERE qq = ?", qq).Scan(&count)
	if err == sql.ErrNoRows {
		count = 0

	} else if err != nil {
		log.Print(err)
		return fmt.Sprint("意外的错误", err)
	}
	// 添加一条
	if count > 0 {
		return "您申请的白名单正在审批"
	}

	result, err := conn.Exec("INSERT INTO whitelist (qq, gameid, status) VALUES (?, ?, ?)",
		qq, gameid, 0) // status 0 表示待审核状态
	if err != nil {
		log.Print(err)
		return fmt.Sprint("添加白名单申请失败:", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Sprint("添加白名单申请失败:", err)
	}
	return "白名单申请已提交,请等待审核"
}
