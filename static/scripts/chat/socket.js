const socket = new WebSocket(APP.URL.WEBSOCKET);

socket.onopen = function () {
    console.log("WebSocket connection established.");
    GetOnlineStatusHandler();
}

socket.onmessage = function (e) {
    const payload = JSON.parse(e.data);

    if (payload.Type == "OnlineStatus") {
        const contact = GetElemByDataTag(document, APP.DATA.CONTACT.ID, payload.UserId)
        RenderContactOnlineStatus(payload.OnlineStatus, contact)
        return;
    }

    if (payload.Type == "RemoveMember") {
        const chatContainer = document.getElementById(APP.ID.CHAT.CONTAINER)
        DeleteElementByDataTag(chatContainer, APP.DATA.CHAT.ID, payload.ChatId)
    }

    if (payload.Type == "AddMember") {
        const activeChat = QSelectByClass(document, APP.CLS.GEN.ACTIVE);
        const activeChatId = GetDataAttribute(activeChat, APP.DATA.CHAT.ID)
        RenderChatElements(payload.Chats, true);
        SetChatToActive(activeChatId)
        return;
    }

    if (payload.Chats != null) {
        if (payload.Messages) {
            if (!payload.Messages[0].IsUserMessage) {
                Object.values(payload.Chats).forEach(chat => { 
                    const chatElem = GetElemByDataTag(document, APP.DATA.CHAT.ID, chat.Id)
                    UpdateUnreadMessageCount(chatElem, chat);
                });
            }
        }
    };

    if (!payload.Messages) return;
    const messageChatId = payload.Messages[0].ChatId;
    
    const activeChat = QSelectByClass(document, APP.CLS.GEN.ACTIVE);
    if (!activeChat) { throw new Error("no active chat found"); };

    const activeChatId = getActiveChatId(activeChat)
    if (!activeChatId) { throw new Error("failed to get active chat id"); };

    if (activeChatId != messageChatId) { return; };
    HandleNewMessageResponse(payload);
    AutoScrollToBottom();
};

function getActiveChatId(activeChat) {
    if (activeChat.classList.contains(APP.CLS.CHAT.TAG)) { 
        return GetDataAttribute(activeChat, APP.DATA.CHAT.ID); 
    };
    if (activeChat.classList.contains(APP.CLS.CONTACT.TAG)) {
        return GetDataAttribute(activeChat, APP.DATA.CONTACT.CHATID);
    };
};

function UpdateUnreadMessageCount(elem, chat) {
    let umc = QSelectByClass(elem, APP.CLS.CHAT.UNREAD_MSG_CNT);
    if (!umc) { return; }
    if (chat.UnreadMessageCount == 0) { HideElement(umc); return; };
    if (elem.classList.contains(APP.CLS.GEN.ACTIVE)) { return; };
    if (umc.innerHTML === chat.UnreadMessageCount) { return; };
    umc.innerHTML = chat.UnreadMessageCount;
    ShowElement(umc);
    PulseElement(umc);
};

const sendContactMessage = (data) => socketSendMessage(data, APP.MSG_TYPE.NEW_CONTACT_MSG);
const sendChatMessage = (data) => socketSendMessage(data, APP.MSG_TYPE.NEW_MSG);

function socketSendMessage(data, msgType) {
    if (data.MsgText) {
        payload = {Type: msgType, Data: data };
        socket.send(JSON.stringify(payload));
    };
};
