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

function HandleNewChatResponse(data) {
    const callbacks = {
        Chats: (data) => RenderChatElements(data, false), 
        Messages: (data) => RenderMessageElements(data, true)
    };
    return HandleResponse(data, callbacks);
};

function HandleSwitchChatResponse(data) {
    const callbacks = {Messages: (data) => RenderMessageElements(data, true)};
    return HandleResponse(data, callbacks);
};

function HandleResponse(data, callbacks) {
    Object.entries(data).forEach(([key, value]) => {
        console.log("key, value = ", key, value);
        console.log("callbacks: ", callbacks);
        if (callbacks[key]) {
            callbacks[key](value);
        } else {
            throw new Error(`No callback found for key: ${key}`);
        };
    });
};