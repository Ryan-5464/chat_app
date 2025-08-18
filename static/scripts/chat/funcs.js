function NewChatHandler(newChatName) {
    NewChatRequest(newChatName).then(data => {
        HandleNewChatResponse(data);
    }).catch(error => {
        console.error("New chat failed => error: ", error);
    });
};

function AddContactHandler(email) {
    AddContactRequest(email).then(data => {
        HandleAddContactResponse(data);
    }).catch(error => {
        console.error("Add contact failed => error: ", error);
    });
};

function DeleteMessageHandler(chatId, messageId, userId) {
    DeleteMessageRequest(chatId, messageId, userId).then(data => {
        HandleDeleteMessageResponse(data);
    }).catch(error => {
        console.error("Delete message failed => error: ", error);
    });
};

function EditChatNameHandler(newName, chatId) {
    EditChatNameRequest(newName, chatId).then(data => {
        HandleEditChatNameResponse(data);
    }).catch(error => {
        console.error("Edit chame name failed => error: ", error);
    });
};

function LeaveChatHandler(chatId, isActive) {
    LeaveChatRequest(chatId).then(data => {
        if (!isActive) { return }
        HandleLeaveChatResponse(data);
    }).catch(error => {
        console.error("Leave chat failed => error: ", error);
    })
}

function SwitchChatHandler(chatType, chatId) {
    SwitchChatRequest(chatType, chatId).then(data => {
        HandleSwitchChatResponse(data);
    }).catch(error => {
        console.error("Switch chat failed => error: ", error);
    })
}