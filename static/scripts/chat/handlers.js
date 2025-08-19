function NewChatHandler(newChatName) {
    console.log("new chat rquest: ", newChatName);
    NewChatRequest(newChatName).then(data => {
        console.log("new chat response handler: ", data)
        HandleNewChatResponse(data);
    }).catch(error => {
        console.error("New chat failed => error: ", error);
    });
};

function AddContactHandler(email) {
    AddContactRequest(email).then(data => {
        console.log("add contact response handler: ", data)
        HandleAddContactResponse(data);
    }).catch(error => {
        console.error("Add contact failed => error: ", error);
    });
};

function DeleteMessageHandler(chatId, messageId, userId) {
    DeleteMessageRequest(chatId, messageId, userId).then(data => {
        console.log("delete messasge response handler: ", data)
        HandleDeleteMessageResponse(data);
    }).catch(error => {
        console.error("Delete message failed => error: ", error);
    });
};

function EditChatNameHandler(newName, chatId) {
    EditChatNameRequest(newName, chatId).then(data => {
        console.log("edit chat name response handler", data)
        HandleEditChatNameResponse(data);
    }).catch(error => {
        console.error("Edit chame name failed => error: ", error);
    });
};

function LeaveChatHandler(chatId, isActive) {
    LeaveChatRequest(chatId).then(data => {
        console.log("leave chat response handler: ", data);
        if (!isActive) { 
            console.log("deleting inactive chat element", chatId)
            DeleteChatElement(chatId);
            return;
        };
        HandleLeaveChatResponse(data);
    }).catch(error => {
        console.error("Leave chat failed => error: ", error);
    });
};

function SwitchChatHandler(chatType, chatId) {
    SwitchChatRequest(chatType, chatId).then(data => {
        console.log("switch chat response handler: ", data)
        HandleSwitchChatResponse(data);
    }).catch(error => {
        console.error("Switch chat failed => error: ", error);
    });
};