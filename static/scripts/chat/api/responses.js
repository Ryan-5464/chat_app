function HandleLeaveChatResponse(data) {
    const callbacks = { 
        NewActiveChatId: (data) => SetChatToActive(data), 
        Chats: (data) => RenderChatElements(data, true), 
        Messages: (data) => RenderMessageElements(data, true)
    };
    return HandleResponse(data, callbacks);
};

function HandleAddContactResponse(data) {
    const callbacks = {
        Contacts: (data) => RenderContactElements(data, false),
        Messages: (data) => RenderMessageElements(data, true)
    };
    return  HandleResponse(data, callbacks);
};

function HandleRemoveContactResponse(data) {
    const callbacks = {
        NewActiveChatId: (data) => SetChatToActive(data),
        Contacts: (data) => RenderContactElements(data, true), 
        Messages: (data) => RenderMessageElements(data, true),
    };
    return HandleResponse(data, callbacks);
};

function HandleDeleteMessageResponse(data) {
    const callbacks = { 
        MessageId: (data) => DeleteMessageElement(data)
    };
    return HandleResponse(data, callbacks);
};

function HandleEditChatNameResponse(data) {
    const callbacks = { 
        Name: (data) => RenderChatNameElement(data)
    };
    return HandleResponse(data, callbacks);
};

function HandleGetOnlineStatusResponse(data) {
    const callbacks = {
        OnlineStatus: (data) => RenderOnlineStatus(data),
    };
    return HandleResponse(data, callbacks);
};

function HandleEditMessageResponse(data) {
    const callbacks = {
        MsgText: (data) => RenderMessageTextElement(data),
    };
    return HandleResponse(data, callbacks);
};

function HandleNewChatResponse(data) {
    const callbacks = {
        Chats: (data) => RenderChatElements(data, false),
        ActiveChatId: (data) => SetChatToActive(data),
        Messages: (data) => RenderMessageElements(data, true)
    };
    return HandleResponse(data, callbacks);
};

function HandleSwitchChatResponse(data) {
    const callbacks = {
        ActiveChatId: (data) => SetChatToActive(data),
        Messages: (data) => RenderMessageElements(data, true)
    };
    return HandleResponse(data, callbacks);
};

function HandleSwitchContactChatResponse(data) {
    const callbacks = {
        ActiveContactChatId: (data) => SetContactChatToActive(data),
        Messages: (data) => RenderMessageElements(data, true),
    }
    return HandleResponse(data, callbacks)
}

function HandleNewMessageResponse(data) {
    const callbacks = {
        Messages: (data) => RenderMessageElements(data, false),
    };
    return HandleResponse(data, callbacks);
};

function HandleGetMemberListResponse(data) {
    const callbacks = {
        Members: (data) => RenderChatMemberElements(data, true),
    };
    return HandleResponse(data, callbacks);
};

function HandleAddMemberResponse(data) {
    const callbacks = {
        Members: (data) => RenderChatMemberElements(data, false),
    }
    return HandleResponse(data, callbacks);
}

function HandleChangeOnlineStatusResponse(data) {
    const callbacks = {
        Status: (data) => RenderOnlineStatus(data)
    };
    return HandleResponse(data, callbacks);
};

function HandleResponse(data, callbacks) {
    Object.entries(data).forEach(([key, value]) => {
        console.log("key, value = ", key, value);
        console.log("callbacks: ", callbacks);
        if (callbacks[key]) {
            callbacks[key](value);
        } else {
            console.log(`Warning => No callback found for key: ${key}`);
        };
    });
};
