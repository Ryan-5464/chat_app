package schema

import (
	"fmt"
)

func Get() []string {
	newUserTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER PRIMARY KEY AUTOINCREMENT,
			%s TEXT,
			%s TEXT NOT NULL UNIQUE,
			%s TEXT NOT NULL,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		UserTable,
		UserId,
		Name,
		Email,
		PwdHash,
		CreatedAt,
	)

	newMessageTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER PRIMARY KEY AUTOINCREMENT,
			%s INTEGER NOT NULL,
			%s INTEGER NOT NULL,
			%s TEXT NOT NULL,
			%s INTEGER,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		MessageTable,
		MessageId,
		UserId,
		ChatId,
		MsgText,
		ReplyId,
		CreatedAt,
		LastEditAt,
	)

	newChatTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER PRIMARY KEY AUTOINCREMENT,
			%s TEXT NOT NULL UNIQUE,
			%s INTEGER NOT NULL,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		ChatTable,
		ChatId,
		Name,
		AdminId,
		CreatedAt,
		LastMsgAt,
	)

	newMemberTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER NOT NULL,
			%s INTEGER NOT NULL,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (%s, %s)
		);`,
		MemberTable,
		ChatId,
		UserId,
		LastReadMsgId,
		ChatId,
		UserId,
	)

	friendsTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER NOT NULL,
			%s INTEGER NOT NULL,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (%s, %s)
		);`,
		FriendsTable,
		UserAId,
		UserBId,
		FriendSince,
		UserAId,
		UserBId,
	)

	var schema []string
	schema = append(
		schema,
		newUserTable,
		newMessageTable,
		newChatTable,
		newMemberTable,
		friendsTable,
	)

	return schema
}
