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

function RemoveContactHandler(contactId, isActive) {
    RemoveContactRequest(contactId).then(data => {
        console.log("remove contact response handler: ", data)
        if (!isActive) { 
            console.log("deleting inactive contact element", contactId)
            DeleteContactElement(contactId);
            return;
        };
        HandleRemoveContactResponse(data);
    }).catch(error => {
        console.error("Remove contact failed => error: ", error);
    });
}

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
        console.error("Edit chat name failed => error: ", error);
    });
};

function EditMessageHandler(messageText, messageId, userId) {
    EditMessageRequest(messageText, messageId, userId).then(data => {
        console.log("edit message response handler", data)
        HandleEditMessageResponse(data);
    }).catch(error => {
        console.error("Edit message failed => error: ", error);
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

function SwitchChatHandler(chatId) {
    SwitchChatRequest(chatId).then(data => {
        console.log("switch chat response handler: ", data)
        HandleSwitchChatResponse(data);
    }).catch(error => {
        console.error("Switch chat failed => error: ", error);
    });
};

function SwitchContactChatHandler(contactChatId) {
    SwitchContactChatRequest(contactChatId).then(data => {
        console.log("switch contact chat response handler: ", data);
        HandleSwitchContactChatResponse(data);
    }).catch(error => {
        console.error("Switch contact chat failed => error: ", error);
    });
};
