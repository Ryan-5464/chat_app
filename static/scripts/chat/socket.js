const socket = new WebSocket(APP.URL.WEBSOCKET);

socket.onopen = function () {
    console.log("WebSocket connection established.");
    GetOnlineStatusHandler();
}

socket.onmessage = function (e) {
    const payload = JSON.parse(e.data);

    if (payload.Chats != null) {
        Object.values(payload.Chats).forEach(chat => { 
            updateUnreadMessageCount(chat);
        });
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

function updateUnreadMessageCount(chat) {
    const chatElem = GetElemByDataTag(APP.DATA.CHAT.ID, chat.Id)
    let umc = QSelectByClass(chatElem, APP.CLS.CHAT.UNREAD_MSG_CNT);
    if (!umc) { return; }
    if (umc.innerHTML == 0) { HideElement(umc); return; };
    if (chatElem.classList.contains(APP.CLS.GEN.ACTIVE)) { return; };
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
