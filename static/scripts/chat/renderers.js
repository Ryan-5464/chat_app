
const RenderChatElements = (data, overwrite) => RenderElements('chats-container', ChatElement, data, overwrite);
const RenderMessageElements = (data, overwrite) => RenderElements('messages-container', MessageElement, data, overwrite);
const RenderContactElements = (data, overwrite) => RenderElements('contacts-container', ContactElement, data, overwrite);
const RenderChatMemberElements = (data, overwrite) => RenderElements('member-list-container', MemberElement, data, overwrite);
const RenderChatNameElement = (data) => ReplaceElement('chat-name-input', ChatNameElement, data);
const RenderMessageTextElement = (data) => ReplaceElement('edit-message-input', MessageTextElement, data)

const DeleteMessageElement = (data) => DeleteElement(`[data-messageid="${data}"]`);
const DeleteChatElement = (data) => DeleteElement(`[data-chatid="${data}"]`)
const DeleteContactElement = (data) => DeleteElement(`[data-contactid="${data}"]`)


function RenderElements(containerId, elemFactory, data, overwrite) {
    console.log("render elements", containerId, elemFactory, data, overwrite)
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

function DeleteElement(identifier) {
    const elem = document.querySelector(identifier);
    if (elem) {
        elem.remove();
    } else {
        throw new Error(`Failed to find element for identifier = ${identifier}`);
    };
};

function SetChatToActive(activeChatId) {
    const chats = document.querySelectorAll('.active');
    if (chats) {
        Object.values(chats).forEach(chat => {
            chat.classList.remove('active');
        });
    }
    const chat = document.querySelector(`[data-chatid="${activeChatId}"]`);
    chat.classList.add('active')
}

function SetContactChatToActive(activeContactChatId) {
    const chats = document.querySelectorAll('.active');
    if (chats) {
        Object.values(chats).forEach(chat => {
            chat.classList.remove('active');
        });
    }
    const chat = document.querySelector(`[data-contactchatid="${activeContactChatId}"]`);
    chat.classList.add('active')
}
