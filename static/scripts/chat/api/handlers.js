
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

function DeleteMessageHandler(messageId, userId) {
    DeleteMessageRequest(messageId, userId).then(data => {
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

function ChangeOnlineStatusHandler(status) {
    ChangeOnlineStatusRequest(status).then(data => {
        console.log("change online status response handler: ", data);
        HandleChangeOnlineStatusResponse(data);
    }).catch(error => {
        console.error("change online status failed => error: ", error);
    });
}

function GetOnlineStatusHandler() {
    GetOnlineStatusRequest().then(data => {
        console.log("get online status response handler: ", data);
        HandleGetOnlineStatusResponse(data);
    }).catch(error => {
        console.error("get online status failed => error: ", error);
    });
}

function SwitchContactChatHandler(contactChatId) {
    SwitchContactChatRequest(contactChatId).then(data => {
        console.log("switch contact chat response handler: ", data);
        HandleSwitchContactChatResponse(data);
    }).catch(error => {
        console.error("Switch contact chat failed => error: ", error);
    });
};

function DisplayMemberListHandler(chatId) {
    GetMemberListRequest(chatId).then(data => {
        console.log("display member list: ", data);
        HandleGetMemberListResponse(data);
    }).catch(error => {
        console.error("Display member list failed => error: ", error);
    });
};

function AddMemberToChatHandler(email, chatId) {
    AddMemberToChatRequest(email, chatId).then(data => {
        console.log("add member to chat: ", data)
        HandleAddMemberResponse(data);
    }).catch(error => {
        console.error("Add member to chat failed => error: ", error);
    });
};

function RemoveMemberHandler(chatId, userId) {
    RemoveMemberRequest(chatId, userId).then(data => {
        console.log("remove member from chat: ", data)
        const members = QSelectAllByClass(document, APP.CLS.MEMBER.TAG)
        Object.values(members).forEach(member => {
            if (GetDataAttribute(member, APP.DATA.USER.ID) === userId) {
                member.remove()
            }
        })
    }).catch(error => {
        console.error("Remove member from chat failed => error: ", error);
    });
}