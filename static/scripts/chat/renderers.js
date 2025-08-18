
const RenderChatElements = (data, overwrite) => RenderElements('chats-container', ChatElement, data, overwrite);
const RenderMessageElements = (data, overwrite) => RenderElements('messages-container', MessageElement, data, overwrite);
const RenderContactElements = (data, overwrite) => RenderElements('contacts-container', ContactElement, data, overwrite);
const RenderChatNameElement = (data) => ReplaceElement('chat-name-input', ChatNameElement, data);

const DeleteMessageElement = (data) => DeleteElement(`[data-messageid="${data.MessageId}"]`);
const DeleteChatElement = (data) => DeleteElement(`[data-chatid="${data.ChatId}"]`)

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

function DeleteElement(identifier) {
    const elem = document.querySelector(identifier);
    if (elem) {
        elem.remove();
    } else {
        throw new Error(`Failed to find element for identifier = ${identifier}="${value}" `);
    };
};

