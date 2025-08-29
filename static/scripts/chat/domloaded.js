document.addEventListener("DOMContentLoaded", function() {
    SetupListeners();
    QSelectAllByClass(document, APP.CLS.MESSAGE.TAG).forEach(msg => { formatMessageDate(msg) });
    AutoScrollToBottom();
});

function formatMessageDate(msgElem) {
    const dateElem = createMessageDate(msgElem);
    const header = CreateElement({classes:[APP.CLS.MESSAGE.HEADER]});
    header.appendChild(dateElem);
    msgElem.insertBefore(header, msgElem.querySelector(APP.CLS.MESSAGE.TEXT));
};

function createMessageDate(msgElem) {
    const createdAt = new Date(msgElem.dataset[APP.DATA.GEN.CREATED]);
    const lastEditAt = new Date(msgElem.dataset[APP.DATA.MESSAGE.LAST_EDIT]);
    if (createdAt < lastEditAt) {
        return CreateElement({
            classes:[APP.CLS.MESSAGE.LAST_EDIT], 
            innerHTML:`Edited: ${FormatDate(msgElem.dataset[APP.DATA.MESSAGE.LAST_EDIT])}`,
        });
    } else {
        return CreateElement({
            classes:[APP.CLS.MESSAGE.CREATED], 
            innerHTML:`Sent: ${FormatDate(msgElem.dataset[APP.DATA.GEN.CREATED])}`,
        });
    };
};

function AutoScrollToBottom() {
    const container = document.getElementById(APP.ID.MESSAGE.CONTAINER);
    requestAnimationFrame(() => { container.scrollTop = container.scrollHeight; });
};

function SetupListeners() {
    addChatModalListenerToChatContainer();
    addMessageModalListenerToMessageContainer();
    addContactModalListenerToContactContainer();
    addSwitchChatControllerToChatsContainer();
    addSwitchContactChatControllerToContactsContainer();
    addNewChatEventListenerToNewChatInput();
    addAddContactEventListenerToAddContactInput();
    ConfigureMemberListModal();
    addNewMsgListenerToMsgInput();
    addNewMsgListenerToSendMsgButton();
};

function addNewMsgListenerToMsgInput() {
    const input = document.getElementById(APP.ID.MESSAGE.INPUT);
    input.addEventListener('keydown', function (event) {
        if (event.key === "Enter") {
            const msgText = input.value.trim();
            const chat = QSelectByClass(document, APP.CLS.GEN.ACTIVE);
            let chatId = GetDataAttribute(chat, APP.DATA.CHAT.ID);
            if (chatId == null) {
                chatId = GetDataAttribute(chat, APP.DATA.CONTACT.CHATID);
                sendContactMessage({MsgText: msgText, ChatId: chatId});
            } else {
                sendChatMessage({MsgText: msgText, ChatId: chatId});
            };
            input.value = '';
        };
    });
};

function addNewMsgListenerToSendMsgButton() {
    const button = document.getElementById(APP.ID.MESSAGE.SEND_BTN);
    button.addEventListener('click', function () {
        const input = document.getElementById(APP.ID.MESSAGE.INPUT);
        const msgText = input.value.trim();
        const chat = QSelectByClass(document, APP.CLS.GEN.ACTIVE);
        let chatId = GetDataAttribute(chat, APP.DATA.CHAT.ID);
        if (chatId == null) {
            chatId = GetDataAttribute(chat, APP.DATA.CONTACT.CHATID);
            sendContactMessage({MsgText: msgText, ChatId: chatId});
        } else {
            sendChatMessage({MsgText: msgText, ChatId: chatId});
        };
        input.value = '';
    });
};
