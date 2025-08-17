
const RenderChatElements = (data, overwrite) => RenderElements('chats-container', ChatElement, data, overwrite);
const RenderMessageElements = (data, overwrite) => RenderElements('messages-container', MessageElement, data, overwrite);
const RenderContactElements = (data, overwrite) => RenderElements('contacts-container', ContactElement, data, overwrite);

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

