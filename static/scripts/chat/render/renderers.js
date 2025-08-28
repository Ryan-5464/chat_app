
const RenderChatElements = (data, overwrite) => RenderElements(APP.ID.CHAT.CONTAINER, ChatElement, data, overwrite);
const RenderMessageElements = (data, overwrite) => RenderElements(APP.ID.MESSAGE.CONTAINER, MessageElement, data, overwrite);
const RenderContactElements = (data, overwrite) => RenderElements(APP.ID.CONTACT.CONTAINER, ContactElement, data, overwrite);
const RenderChatMemberElements = (data, overwrite) => RenderElements(APP.ID.MEMBER.CONTAINER, MemberElement, data, overwrite);
const RenderChatNameElement = (data) => ReplaceElement(APP.ID.CHAT.NAME_INPUT, ChatNameElement, data);
const RenderMessageTextElement = (data) => ReplaceElement(APP.ID.MESSAGE.EDIT_INPUT, MessageTextElement, data)

const DeleteMessageElement = (data) => DeleteElementByDataTag(APP.DATA.MESSAGE.ID, data);
const DeleteChatElement = (data) => DeleteElementByDataTag(APP.DATA.CHAT.ID, data);
const DeleteContactElement = (data) => DeleteElementByDataTag(APP.DATA.CONTACT.ID, data);

function RenderElements(containerId, elemFactory, data, overwrite) {
    const container = document.getElementById(containerId);
    if (!container) {
        throw new Error(`Element with id=${containerId} not found!`);
    };
    if (overwrite == true) {
        container.innerHTML = '';
    };
    Object.values(data).forEach(obj => {
        container.appendChild(elemFactory(obj));
    });
};

function ReplaceElement(elementId, elemFactory, data) {
    const elem = document.getElementById(elementId);
    elem.replaceWith(elemFactory(data));
};

function SetChatToActive(activeChatId) {
    RemoveActiveFromChat();
    const chat = GetElemByDataTag(APP.DATA.CHAT.ID, activeChatId);
    umc = QSelectByClass(chat, APP.CLS.CHAT.UNREAD_MSG_CNT);
    if (!umc) { return; };
    umc.innerHTML = 0;
    umc.classList.add(APP.CLS.HIDDEN);
    chat.classList.add(APP.CLS.ACTIVE);
    AutoScrollToBottom();
};

function SetContactChatToActive(activeContactChatId) {
    RemoveActiveFromChat();
    const chat = GetElemByDataTag(APP.DATA.CONTACT.CHATID, activeContactChatId);
    chat.classList.add(APP.CLS.ACTIVE);
};

function RenderOnlineStatus(status) {
    const onlineStatus = QSelectByClass(document, APP.CLS.STATUS.STATUS);
    onlineStatus.innerHTML = status;
    onlineStatus.classList.value = '';
    changeOnlineStatus(onlineStatus, status);
};

function RenderContactOnlineStatus(status, contact) {
    const onlineStatus = QSelectByClass(contact, APP.CLS.CONTACT.STATUS);
    if (status == APP.CLS.STATUS.STEALTH) {
        onlineStatus.innerHTML = "Offline";
    } else {
        onlineStatus.innerHTML = status;
    };
    onlineStatus.classList.value = '';
    onlineStatus.classList.add(APP.CLS.CONTACT.STATUS);
    changeOnlineStatus(onlineStatus, status);
};

function changeOnlineStatus(elem, status) {
    if (status == APP.CLS.STATUS.ONLINE) { elem.classList.add(APP.CLS.STATUS.ONLINE) };
    if (status == APP.CLS.STATUS.AWAY) { elem.classList.add(APP.CLS.STATUS.AWAY) };
    if (status == APP.CLS.STATUS.BUSY) { elem.classList.add(APP.CLS.STATUS.BUSY) };
    if (status == APP.CLS.STATUS.STEALTH) { elem.classList.add(APP.CLS.STATUS.OFFLINE) };
    if (status == APP.CLS.STATUS.OFFLINE) { elem.classList.add(APP.CLS.STATUS.OFFLINE) };
};

function HideElement(elem) {
    elem.classList.add(APP.CLS.HIDDEN);
};

function ShowElement(elem) {
    elem.classList.remove(APP.CLS.HIDDEN);
};

function RemoveActiveFromChat() {
    const activeChat = QSelectByClass(document, APP.CLS.ACTIVE);
    if (!activeChat) { return; };
    activeChat.classList.remove(APP.CLS.ACTIVE);
};