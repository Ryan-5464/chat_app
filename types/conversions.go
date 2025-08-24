package types

import "strconv"

func ConvertInt64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ConvertStringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ToChatId(id string) (ChatId, error) {
	chatId, err := ConvertStringToInt64(id)
	if err != nil {
		return 0, err
	}
	return ChatId(chatId), nil
}

func ToContactId(id string) (ContactId, error) {
	contactId, err := ConvertStringToInt64(id)
	if err != nil {
		return 0, err
	}
	return ContactId(contactId), nil
}

func ToMessageId(id string) (MessageId, error) {
	messageId, err := ConvertStringToInt64(id)
	if err != nil {
		return 0, err
	}
	return MessageId(messageId), nil
}

func ToUserId(id string) (UserId, error) {
	userId, err := ConvertStringToInt64(id)
	if err != nil {
		return 0, err
	}
	return UserId(userId), nil
}
