function addSwitchChatControllerToChatsContainer() {
    const chatContainer = ConfigureChatsController();
    chatContainer.addEventListener('click', (e) => {
        const chatId = e.target.closest('[data-chatid]')?.getAttribute('data-chatid');
        chatContainer.__controller.SwitchChat(chatId)
    }); 
};

function ConfigureChatsController() {
    const chatsContainer = document.getElementById('chats-container')
    chatsContainer.__controller = {
        SwitchChat: (chatId) => SwitchChatHandler(chatId),
    };
    return chatsContainer
};


function addSwitchContactChatControllerToContactsContainer() {
    const contactContainer = ConfigureContactsController();
    contactContainer.addEventListener('click', (e) => {
        const contactChatId = e.target.closest('[data-contactchatid]')?.getAttribute('data-contactchatid');
        contactContainer.__controller.SwitchContactChat(contactChatId)
    }); 
};

function ConfigureContactsController() {
    const contactsContainer = document.getElementById('contacts-container')
    contactsContainer.__controller = {
        SwitchContactChat: (contactChatId) => SwitchContactChatHandler(contactChatId),
    };
    return contactsContainer
};