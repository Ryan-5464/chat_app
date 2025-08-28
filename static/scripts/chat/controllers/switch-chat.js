function addSwitchChatControllerToChatsContainer() {
    const chatContainer = ConfigureChatsController();
    chatContainer.addEventListener('click', (e) => {
        const chat = GetClosestTargetByData(e, APP.DATA.CHAT.ID)
        const chatId = GetDataAttribute(chat, APP.DATA.CHAT.ID);
        chatContainer.__controller.SwitchChat(chatId);
    }); 
};

function ConfigureChatsController() {
    const chatsContainer = document.getElementById(APP.CLS.CHAT.CONTAINER);
    chatsContainer.__controller = {
        SwitchChat: (chatId) => SwitchChatHandler(chatId),
    };
    return chatsContainer;
};


function addSwitchContactChatControllerToContactsContainer() {
    const contactContainer = ConfigureContactsController();
    contactContainer.addEventListener('click', (e) => {
        const contact = GetClosestTargetByData(e, APP.DATA.CONTACT.CHATID);
        const contactChatId = GetDataAttribute(contact, APP.DATA.CONTACT.CHATID);
        contactContainer.__controller.SwitchContactChat(contactChatId);
    }); 
};

function ConfigureContactsController() {
    const contactsContainer = document.getElementById(APP.CLS.CONTACT.CONTAINER)
    contactsContainer.__controller = {
        SwitchContactChat: (contactChatId) => SwitchContactChatHandler(contactChatId),
    };
    return contactsContainer
};