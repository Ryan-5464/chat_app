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
			%s DATETIME DEFAULT CURRENT_TIMESTAMP,
    		%s TEXT NOT NULL CHECK (%s IN ('%s', '%s'))
		);`,
		ChatTable,
		ChatId,
		Name,
		AdminId,
		CreatedAt,
		LastMsgAt,
		ChatType,
		ChatType,
		Private,
		Group,
	)

	newMemberTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER PRIMARY KEY AUTOINCREMENT,
			%s INTEGER NOT NULL,
			%s INTEGER NOT NULL,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		MemberTable,
		RowId,
		ChatId,
		UserId,
		LastReadMsgId,
	)

	contactsTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER PRIMARY KEY AUTOINCREMENT,
			%s INTEGER NOT NULL,
			%s INTEGER NOT NULL,
			%s DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		ContactsTable,
		RowId,
		Contact1,
		Contact2,
		Established,
	)

	enforcePrivateChatReadOnly := fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS prevent_chat_type_change
		BEFORE UPDATE ON %s
		FOR EACH ROW
		WHEN OLD.%s != NEW.%s
		BEGIN
			SELECT RAISE(ABORT, 'Cannot change chat type after creation');
		END;
	`,
		ChatTable,
		ChatType,
		ChatType,
	)

	enforceTwoMemberLimitForPrivateChats := fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS prevent_extra_private_members
		BEFORE INSERT ON %s
		WHEN (
			(SELECT %s FROM %s WHERE %s = NEW.%s) = '%s' AND
			(SELECT COUNT(*) FROM %s WHERE %s = NEW.%s) >= 2
		)
		BEGIN
			SELECT RAISE(ABORT, 'Private chat cannot have more than 2 members');
		END;
	`,
		MemberTable,
		ChatType,
		ChatTable,
		ChatId,
		ChatId,
		Private,
		MemberTable,
		ChatId,
		ChatId,
	)

	var schema []string
	schema = append(
		schema,
		newUserTable,
		newMessageTable,
		newChatTable,
		newMemberTable,
		contactsTable,
		enforcePrivateChatReadOnly,
		enforceTwoMemberLimitForPrivateChats,
	)

	return schema
}
