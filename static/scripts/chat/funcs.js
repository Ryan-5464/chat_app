function NewChat(newChatName) {
    NewChatRequest(newChatName).then(data => {
        HandleNewChatResponse(data);
    }).catch(error => {
        console.error("New chat failed => error: ", error);
    });
};

function AddContact(email) {
    AddContactRequest(email).then(data => {
        HandleAddContactResponse(data);
    }).catch(error => {
        console.error("Add contact failed => error: ", error);
    });
};

function DeleteMessage(chatId, messageId, userId) {
    DeleteMessageRequest(chatId, messageId, userId).then(data => {
        HandleDeleteMessageResponse(data);
    }).catch(error => {
        console.error("Delete message failed => error: ", error);
    });
};

function EditChatName(newName, chatId) {
    EditChatNameRequest(newName, chatId).then(data => {
        HandleEditChatNameResponse(data);
    }).catch(error => {
        console.error("Edit chame name failed => error: ", error);
    });
};

function LeaveChat(chatId, isActive) {
    LeaveChatRequest(chatId).then(data => {
        if (!isActive) { return }
        HandleLeaveChatResponse(data);
    }).catch(error => {
        console.error("Leave chat failed => error: ", error);
    })
}

function SwitchChat(chatType, chatId) {
    SwitchChatRequest(chatType, chatId).then(data => {
        HandleSwitchChatResponse(data);
    }).catch(error => {
        console.error("Switch chat failed => error: ", error);
    })
}